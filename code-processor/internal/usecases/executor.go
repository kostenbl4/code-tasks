package usecases

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	"github.com/docker/docker/pkg/stdcopy"
)

type CodeExecutor interface {
	Execute(context.Context, string, string) (string, string, error)
}

func NewCodeExecutor(client *client.Client) CodeExecutor {

	if err := loadImages(client); err != nil {
		log.Fatal(err)
	}

	return codeExecutor{
		client: client,
	}
}
func loadImages(cli *client.Client) error {
	ctx := context.Background()

	if err := pullImage("golang:1.24-alpine", cli, &ctx); err != nil {
		return err
	}

	if err := pullImage("gcc:latest", cli, &ctx); err != nil {
		return err
	}

	if err := pullImage("python:alpine", cli, &ctx); err != nil {
		return err
	}
	log.Println("all images loaded")
	return nil
}

func pullImage(imageName string, cli *client.Client, ctx *context.Context) error {
	res, err := cli.ImagePull(*ctx, imageName, image.PullOptions{})
	// проверяем на ошибку
	if err != nil {
		return err
	}
	// в конце программы закрываем данные, которые получили. Не сам образ, а именно сообщение о его удачном получении
	defer res.Close()

	// из-за того, что полученные данные храняться в io.ReadCloser, их можно вывести в консоль таким образом
	io.Copy(os.Stdout, res)
	return nil
}

type codeExecutor struct {
	client *client.Client
}

func (ce codeExecutor) Execute(ctx context.Context, code, lang string) (string, string, error) {

	filename, err := getFilenameByLang(lang)
	if err != nil {
		return "", "", err
	}

	containerConfig, err := getConfigByLang(lang)
	if err != nil {
		return "", "", err
	}

	hostConfig := setLimits(container.HostConfig{})

	// Создание тар-архива для копирования кода в контейнер
	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)
	tarHeader := &tar.Header{
		Name: filename,
		Size: int64(len(code)),
	}
	if err := tarWriter.WriteHeader(tarHeader); err != nil {
		return "", "", err
	}
	if _, err := tarWriter.Write([]byte(code)); err != nil {
		return "", "", err
	}
	if err := tarWriter.Close(); err != nil {
		return "", "", err
	}

	// Создание и запуск контейнера
	containerResp, err := ce.client.ContainerCreate(ctx, containerConfig, &hostConfig, nil, nil, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %v", err)
	}

	// Копирование кода в контейнер
	err = ce.client.CopyToContainer(ctx, containerResp.ID, "/app", bytes.NewReader(tarBuffer.Bytes()), container.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
	if err != nil {
		return "", "", fmt.Errorf("failed to copy code to container: %v", err)
	}

	// Запуск контейнера
	if err := ce.client.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Ожидание завершения контейнера
	statusCh, errCh := ce.client.ContainerWait(ctx, containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			ce.client.ContainerRemove(context.Background(), containerResp.ID, container.RemoveOptions{Force: true})
			return "", "", fmt.Errorf("container runtime error: %v", err)
		}
	case status := <-statusCh:
		log.Printf("Container exited with status: %v", status)
	}

	// Получение логов
	logs, err := ce.client.ContainerLogs(ctx, containerResp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logs.Close()

	// Удаление контейнера
	if err := ce.client.ContainerRemove(context.Background(), containerResp.ID, container.RemoveOptions{Force: true}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}

	// Копирование логов
	var stdoutBuffer, stderrBuffer bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuffer, &stderrBuffer, logs); err != nil {
		return "", "", fmt.Errorf("failed to copy logs: %v", err)
	}

	return stdoutBuffer.String(), stderrBuffer.String(), nil
}

func getConfigByLang(lang string) (*container.Config, error) {

	var containerConfig *container.Config

	// Создание конфигурации по выбранному языку
	switch lang {
	case "go":
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "golang:1.24-alpine",
			Cmd:        []string{"go", "run", "/app/main.go"},
		}
	case "python":
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "python:alpine",
			Cmd:        []string{"python", "/app/main.py"},
		}
	case "c":
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "gcc:latest",
			Cmd:        []string{"sh", "-c", "gcc /app/main.c -o main && ./main"},
		}
	case "cpp":
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "gcc:latest",
			Cmd:        []string{"sh", "-c", "g++ /app/main.cpp -o main && ./main"},
		}
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
	return containerConfig, nil
}

func getFilenameByLang(lang string) (string, error) {
	switch lang {
	case "go":
		return "main.go", nil
	case "python":
		return "main.py", nil
	case "c":
		return "main.c", nil
	case "cpp":
		return "main.cpp", nil
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
}

func setLimits(hostConfig container.HostConfig) container.HostConfig {
	hostConfig.Resources = container.Resources{
		CPUPeriod: 1000000,
		CPUQuota:  1000000,
		Memory:    128 * 1024 * 1024,
	}
	return hostConfig
}
