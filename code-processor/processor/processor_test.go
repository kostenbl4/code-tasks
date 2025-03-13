package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteCode(t *testing.T) {
	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "code")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}`
	// Сохраняем код в файл
	mainFile := filepath.Join(tmpDir, "main.go")
	t.Logf("Saving code to %s", mainFile)
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	// считаем код обратно
	data, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatal(err)
	}
	if data == nil {
		t.Fatal("code is empty")
	}
	t.Logf("Read code from %s", mainFile)
	// выводим его
	t.Log(string(data))

}
