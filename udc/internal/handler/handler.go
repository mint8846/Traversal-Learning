package handler

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mint8846/Traversal-Learning/udc/internal/filter"
	"github.com/mint8846/Traversal-Learning/udc/internal/model"
	"github.com/mint8846/Traversal-Learning/udc/internal/service"
	"github.com/mint8846/Traversal-Learning/udc/internal/utils"
)

var (
	internalErrorCode = "999"
)

func Model(c echo.Context) error {
	seedKey, err := filter.GetSessionID(c)
	if err != nil {
		return err
	}

	encryptKey, err := service.Default.OTP.Generate(seedKey, service.Default.Cfg.OTPPeriod, time.Now())
	if err != nil {
		log.Printf("Model: OTP Generate fail %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "key generate fail"})
	}

	sessionID := utils.HashB64([]byte(seedKey))
	dirPath := service.Default.NFS.GetPath(sessionID)

	fileName, err := service.Default.File.EncryptFile(service.Default.Cfg.ModelPath, dirPath, encryptKey)
	if err != nil {
		log.Printf("EncryptFile failed: %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "EncryptFile failed"})
	}
	log.Printf("Model: path(%s/%s)", dirPath, fileName)

	nfsURL, err := service.Default.NFS.GenerateNFSUrl(service.Default.Cfg.ServerHost, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "Create Model Path failed"})
	}
	log.Printf("Model: nfsPath(%s)", nfsURL)

	return c.JSON(http.StatusOK, model.SetupResponse{
		ID:       service.Default.Cfg.ID,
		Path:     nfsURL,
		FileName: fileName,
	})
}

func Result(c echo.Context) error {
	seedKey, err := filter.GetSessionID(c)
	if err != nil {
		return err
	}

	decryptKey, err := service.Default.OTP.Generate(seedKey, service.Default.Cfg.OTPPeriod, time.Now())
	if err != nil {
		log.Printf("Result: OTP Generate fail %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "key generate fail"})
	}

	req := new(model.ResultDataRequest)
	if err = c.Bind(req); err != nil {
		return err
	}
	sessionKey := utils.HashB64([]byte(seedKey))
	resultDir := service.Default.NFS.GetPath(sessionKey)
	log.Printf("Result: result dir path(%s)", resultDir)

	filePath, err := service.Default.File.DecryptFile(resultDir, req.FileName, service.Default.Cfg.ResultDir, decryptKey)
	if err != nil {
		log.Printf("Result: decrypt fail %v", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: internalErrorCode, Message: "decrypt fail"})
	}

	go processModelAndCleanup(sessionKey, filePath, resultDir)

	return c.JSON(http.StatusOK, nil)
}

func processModelAndCleanup(sessionKey, resultPath, resultDir string) {
	if err := os.RemoveAll(resultDir); err != nil {
		log.Printf("Result: clean(%s) failed(%v)", resultDir, err)
	}

	if err := service.Default.Runner.ExecuteModel(service.Default.Cfg.ModelPath, resultPath); err != nil {
		log.Printf("Result: (%s)model execute fial(%v)", sessionKey, err)
	}
}
