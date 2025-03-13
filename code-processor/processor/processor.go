package processor

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func ExecuteCode(cli *client.Client, code, lang string) (string, string, error) {
	var (
		filename        string
		containerConfig *container.Config
	)

	// Создание конфигурации по выбранному языку
	switch lang {
	case "go":
		filename = "main.go"
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "golang:1.24-alpine",
			Cmd:        []string{"go", "run", "/app/main.go"},
		}
	case "python":
		filename = "main.py"
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "python:alpine",
			Cmd:        []string{"python", "/app/main.py"},
		}
	case "c":
		filename = "main.c"
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "gcc:latest",
			Cmd:        []string{"sh", "-c", "gcc /app/main.c -o main && ./main"},
		}
	case "cpp":
		filename = "main.cpp"
		containerConfig = &container.Config{
			WorkingDir: "/app",
			Image:      "gcc:latest",
			Cmd:        []string{"sh", "-c", "g++ /app/main.cpp -o main && ./main"},
		}
	default:
		return "", "", fmt.Errorf("unsupported language: %s", lang)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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
	containerResp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %v", err)
	}

	// Копирование кода в контейнер
	err = cli.CopyToContainer(ctx, containerResp.ID, "/app", bytes.NewReader(tarBuffer.Bytes()), container.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
	if err != nil {
		return "", "", fmt.Errorf("failed to copy code to container: %v", err)
	}

	// Запуск контейнера
	if err := cli.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	// Ожидание завершения контейнера
	statusCh, errCh := cli.ContainerWait(ctx, containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			cli.ContainerRemove(context.Background(), containerResp.ID, container.RemoveOptions{Force: true})
			return "", "", fmt.Errorf("container runtime error: %v", err)
		}
	case status := <-statusCh:
		log.Printf("Container exited with status: %v", status)
	}

	// Получение логов
	logs, err := cli.ContainerLogs(ctx, containerResp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logs.Close()

	// Удаление контейнера
	if err := cli.ContainerRemove(context.Background(), containerResp.ID, container.RemoveOptions{Force: true}); err != nil {
		log.Printf("Warning: Failed to remove container: %v", err)
	}

	// Копирование логов
	var stdoutBuffer, stderrBuffer bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdoutBuffer, &stderrBuffer, logs); err != nil {
		return "", "", fmt.Errorf("failed to copy logs: %v", err)
	}

	return stdoutBuffer.String(), stderrBuffer.String(), nil
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

func LoadImages(cli *client.Client) error {
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
