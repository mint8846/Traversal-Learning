package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/mint8846/Traversal-Learning/odc/internal/config"
)

// FileService example, use AES-256 as an example of implementing AES-256 encryption.
// Encryption strength can be enhanced based on implementation.
type FileService struct {
	cfg *config.Config
}

func NewFileService(cfg *config.Config) *FileService {
	return &FileService{cfg: cfg}
}

const (
	BufferSize = 1024 * 1024 // 1MB
)

func (f *FileService) EncryptFile(srcFile, dstDir string, key []byte) (string, error) {
	inputFile, err := os.Open(srcFile)
	if err != nil {
		log.Printf("EncryptFile: read file fail %v", err)
		return "", err
	}
	defer inputFile.Close()

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		log.Fatal("EncryptFile: Failed to create directory:", err)
	}

	outputFile, err := os.Create(filepath.Join(dstDir, uuid.NewString()))
	if err != nil {
		log.Printf("EncryptFile: write file create fail %v", err)
		return "", err
	}
	defer outputFile.Close()

	crypto, err := NewCryptoServiceWithKey(key)
	if err != nil {
		log.Printf("EncryptFile: set crypto fail %v", err)
		return "", err
	}

	log.Printf("encrypt file start(%s)", srcFile)
	if err = f.processFile(inputFile, outputFile, crypto); err != nil {
		return "", err
	}
	log.Printf("encrypt file complete(%s)", srcFile)

	return filepath.Base(outputFile.Name()), nil
}

func (f *FileService) DecryptFile(srcDirPath, fileName, dstDirPath string, key []byte) (string, error) {
	inputFile, err := os.Open(filepath.Join(srcDirPath, fileName))
	if err != nil {
		log.Printf("DecryptFile: read file fail %v", err)
		return "", err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(filepath.Join(dstDirPath, fileName))
	if err != nil {
		log.Printf("DecryptFile: write file create fail %v", err)
		return "", err
	}
	defer outputFile.Close()

	log.Printf("decrypt file start(%s)", fileName)
	crypto, err := NewCryptoServiceWithKey(key)
	if err != nil {
		log.Printf("EncryptFile: set crypto fail %v", err)
		return "", err
	}

	if err = f.processFile(inputFile, outputFile, crypto); err != nil {
		return "", err
	}
	log.Printf("decrypt file complete(%s)", fileName)

	return outputFile.Name(), nil
}

func (f *FileService) Write(path, data string) error {
	inputFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("write: file open fail %v", err)
	}
	defer inputFile.Close()

	if _, err := inputFile.Write([]byte(data)); err != nil {
		return fmt.Errorf("write: fail %v", err)
	}

	return nil
}

func (f *FileService) processFile(inputFile, outputFile *os.File, crypto *CryptoService) error {
	stream, err := crypto.CreateCipher()
	if err != nil {
		return err
	}

	inputBuffer := make([]byte, BufferSize)
	outputBuffer := make([]byte, BufferSize)

	for {
		// file read
		n, err := inputFile.Read(inputBuffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read file fail: %v", err)
		}

		if n == 0 {
			break // EOF
		}

		// Encrypt or Decrypt
		crypto.ProcessData(stream, inputBuffer, outputBuffer, n)

		// write file
		if _, err := outputFile.Write(outputBuffer[:n]); err != nil {
			return fmt.Errorf("write file fail: %v", err)
		}
	}

	return nil
}
