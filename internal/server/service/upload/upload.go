package service

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type UploadService struct{}

func New() *UploadService {
	return &UploadService{}
}

func (f *UploadService) Upload(stream desc.UploadV1_UploadFileServer) error {
	file := fl.NewFile()
	var fileSize uint32
	fileSize = 0
	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), "client_files")
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return err
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		log.Printf("received a chunk with size: %d\n", fileSize)

		if err := file.Write(chunk); err != nil {
			log.Println(err)
			return err
		}
	}

	fmt.Println("result:", file.FilePath, fileSize)
	fileName := filepath.Base(file.FilePath)
	
	return stream.SendAndClose(&desc.UploadFileResponse{FileName: fileName, Size: fileSize})
}
