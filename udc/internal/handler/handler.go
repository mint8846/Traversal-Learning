package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mint8846/Traversal-Learning/udc/internal/model"
	"github.com/mint8846/Traversal-Learning/udc/internal/service"
)

var udcID = ""
var internalErrorCode = "999"
var udcKey []byte = nil

func SetId(id string) {
	udcID = id
}

func Connect(c echo.Context) error {
	return c.JSON(http.StatusOK, model.ConnectResponse{ID: udcID})
}

func Model(c echo.Context) error {
	dirPath := filepath.Join(service.Default.Cfg.NFSPath, service.Default.Cfg.HostName)

	newKey, fileName, err := service.Default.File.EncryptFile(service.Default.Cfg.ModelPath, dirPath, udcKey)
	if err != nil {
		log.Printf("EncryptFile failed: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "EncryptFile failed"})
	}
	log.Printf("Model: fileName(%s)", fileName)

	// Set key if it was nil
	if udcKey == nil {
		udcKey = newKey
	}

	udcKey, err := service.Default.OTP.EncryptKey(udcKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "EncryptKey failed"})
	}
	log.Printf("Model: UDC key(%s)", udcKey)

	nfsURL, err := service.Default.NFS.GenerateNFSUrl()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "Create Model Path failed"})
	}
	log.Printf("Model: nfsPath(%s)", nfsURL)

	return c.JSON(http.StatusOK, model.SetupResponse{
		Key:      udcKey,
		Path:     nfsURL,
		FileName: fileName,
	})
}

func Result(c echo.Context) error {
	req := new(model.ResultDataRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	resultDir := filepath.Join(service.Default.Cfg.NFSPath, service.Default.Cfg.HostName)
	log.Printf("Result: result dir path(%s)", resultDir)

	_, err := service.Default.File.DecryptFile(resultDir, req.FileName, service.Default.Cfg.ResultDir, udcKey)
	if err != nil {
		return fmt.Errorf("result: decrpyt fail %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond) // wait response time ..
		log.Printf("Result: Process completed, sending SIGTERM to self...")

		if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
			log.Printf("Failed to send SIGTERM: %v", err)
			os.Exit(0) // if SIGTERM fail.
		}
	}()

	return c.JSON(http.StatusOK, nil)
}
