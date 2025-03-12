package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func ExecuteCode(cli *client.Client, code string, lang string) (string, string, error) {

	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "code")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	var filename string
	var containerConfig *container.Config
	// В зависимости от языка выбираем конфигурацию
	switch lang {
	case "go":
		filename = "main.go"
		containerConfig = &container.Config{
			Image: "golang:1.24-alpine",
			Cmd:   []string{"go", "run", "/app/main.go"},
		}
	case "python":
		filename = "main.py"
		containerConfig = &container.Config{
			Image: "python:alpine",
			Cmd:   []string{"python", "/app/main.py"},
		}
	case "c":
		filename = "main.c"
		containerConfig = &container.Config{
			Image: "gcc:latest",
			Cmd:   []string{"sh", "-c", "gcc /app/main.c -o main && ./main"},
		}
	case "cpp":
		filename = "main.cpp"
		containerConfig = &container.Config{
			Image: "gcc:latest",
			Cmd:   []string{"sh", "-c", "g++ /app/main.cpp -o main && ./main"},
		}
	default:
		return "", "", fmt.Errorf("unknown language: %s", lang)
	}
	// Сохраняем код в файл
	mainGo := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(mainGo, []byte(code), 0644); err != nil {
		return "", "", err
	}

	// Конфигурация контейнера
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	hostConfig := &container.HostConfig{
		Binds:       []string{tmpDir + ":/app"},
		NetworkMode: "none",
	}

	// Создаем контейнер
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", "", fmt.Errorf("container create error: %v", err)
	}

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("container start error: %v", err)
	}

	// Ожидаем завершения
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			err = cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{
				Force: true,
			})
			if err != nil {
				return "", "", fmt.Errorf("container remove error: %v", err)
			}
			return "", "", fmt.Errorf("container runtime error: %v", err)
		}
	case status := <-statusCh:
		log.Printf("Exit status: %d", status.StatusCode)
	}

	// Получаем логи
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("logs error: %v", err)
	}
	defer out.Close()

	// Удаляем контейнер после получения логов
	err = cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return "", "", fmt.Errorf("container remove error: %v", err)
	}

	// Копируем вывод
	var outputBuf, errorBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outputBuf, &errorBuf, out)
	if err != nil {
		return "", "", fmt.Errorf("copy error: %v", err)
	}

	return outputBuf.String(), errorBuf.String(), nil
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
