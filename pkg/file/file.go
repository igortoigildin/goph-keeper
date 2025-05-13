package file

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/igortoigildin/goph-keeper/internal/client/grpc/models"
)

type File struct {
	FilePath   string
	buffer     *bytes.Buffer
	OutputFile *os.File
}

func NewFile() *File {
	return &File{
		buffer: &bytes.Buffer{},
	}
}

func (f *File) SetFile(fileName, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	f.FilePath = filepath.Join(path, fileName)
	file, err := os.Create(f.FilePath)
	if err != nil {
		return err
	}
	f.OutputFile = file
	return nil
}

func (f *File) Write(chunk []byte) error {
	if f.OutputFile == nil {
		return nil
	}
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) Close() error {
	return f.OutputFile.Close()
}

func (f *File) Remove() error {
	err := os.Remove(f.FilePath)
	if err != nil {
		return err
	}

	return nil
}

func SaveFileToDisk(file models.File, dir string) error {
	path := filepath.Join(dir, file.Filename)

	err := os.WriteFile(path, file.Data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении файла '%s': %w", path, err)
	}
	return nil
}
