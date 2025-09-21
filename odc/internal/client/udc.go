package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mint8846/Traversal-Learning/odc/internal/config"
	"github.com/mint8846/Traversal-Learning/odc/internal/model"
	"github.com/mint8846/Traversal-Learning/odc/internal/service"
	"github.com/mint8846/Traversal-Learning/odc/internal/utils"
)

type UDCClient struct {
	cfg    *config.Config
	http   *HTTPClient
	otp    *service.OTPService
	nfs    *service.NFSService
	file   *service.FileService
	runner *service.RunnerService
}

func NewUDCClient(cfg *config.Config) *UDCClient {
	httpClient := NewHTTPClient(cfg.UDC.Url, cfg.UDC.HTTPClientTimeout)
	httpClient.AddDefaultHeader("X-Container-ID", cfg.ID)
	httpClient.AddDefaultHeader("Content-Type", "application/json")

	return &UDCClient{
		cfg:    cfg,
		http:   httpClient,
		otp:    &service.OTPService{},
		nfs:    service.NewNFSService(cfg),
		file:   service.NewFileService(cfg),
		runner: service.NewRunnerService(cfg),
	}
}

func (u *UDCClient) Connect() error {
	resp, err := u.http.Get("/api/connect", WithHeader("X-SEED-ID", u.cfg.OTP.Seed))

	if err != nil {
		return fmt.Errorf("connect: request Fail: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("connect: response Read Fail: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connect: response Error(%d) %s", resp.StatusCode, string(body))
	}

	var response model.ConnectResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Connect: body data(%s)", body)
		return fmt.Errorf("JSON parsing error: %w", err)
	}

	if u.cfg.UDC.ID != "" && u.cfg.UDC.ID != response.ID {
		log.Fatalf("Connect: Access denied ID %s is not allowed (expected: %s)", response.ID, u.cfg.UDC.ID)
	}

	u.http.AddDefaultHeader("X-Session-Key", u.cfg.SessionKey)

	log.Printf("Connect success (%s)", utils.HashB64([]byte(u.cfg.OTP.Seed)))
	return nil
}

func (u *UDCClient) GetModel() (string, error) {
	requestTime := time.Now()

	resp, err := u.http.Post("/api/model", nil)
	if err != nil {
		return "", fmt.Errorf("GetModel: Request Fail: %v", err)
	}
	defer resp.Body.Close()

	body, err := u.http.GetBody(resp)
	if err != nil {
		return "", fmt.Errorf("GetModel: response error %v", err)
	}

	var response model.SetupResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("GetModel: body data(%s)", body)
		return "", fmt.Errorf("JSON parsing error: %v", err)
	}

	if err = u.nfs.Connect(response.Path); err != nil {
		return "", fmt.Errorf("GetModel: NFS Connect fail %v", err)
	}

	key, err := u.getKey(requestTime)
	if err != nil {
		return "", fmt.Errorf("GetModel: get decrypt key fail %v", err)
	}

	log.Printf("GetModel: model path(%s/%s)", u.nfs.GetPath(""), response.FileName)
	modelPath, err := u.file.DecryptFile(u.nfs.GetPath(""), response.FileName, "/tmp", key)

	if err != nil {
		return "", fmt.Errorf("GetModel: model decrpyt fail %v", err)
	}
	return modelPath, nil
}

func (u *UDCClient) ExecuteModel(modelPath string) error {
	log.Printf("ExecuteModel: start(%s)", modelPath)
	if err := u.runner.ExecuteModel(modelPath); err != nil {
		return fmt.Errorf("GetModel: model execute fail %v", err)
	}

	if err := u.runner.CheckResultData(u.cfg.ModelOutputPath); err != nil {
		return fmt.Errorf("GetModel: model result file error %v", err)
	}
	return nil
}

func (u *UDCClient) EncryptResult() error {
	key, err := u.getKey(time.Now())
	if err != nil {
		return fmt.Errorf("EncryptResult: get encrypt key fail %v", err)
	}

	fileName, err := u.file.EncryptFile(u.cfg.ModelOutputPath, u.nfs.GetPath(""), key)
	if err != nil {
		return fmt.Errorf("EncryptResult: EncryptFile error %v", err)
	}
	log.Printf("EncryptResult: fileName(%s)", fileName)

	// From the start of encryption to transmission, it must not exceed the config's OTP.Period time (encryption key changes)
	// To resolve this, passing the encryption start time is necessary
	resp, err := u.http.Post("/api/result", model.ResultDataRequest{FileName: fileName})
	if err != nil {
		return fmt.Errorf("EncryptResult: send result info fail %v", err)
	}
	defer resp.Body.Close()

	if _, err = u.http.GetBody(resp); err != nil {
		log.Printf("EncryptResult: response error %v", err)
		return err
	}

	if err = u.file.Write(u.cfg.ResultFile, "success"); err != nil {
		log.Printf("EncryptResult: write error %v", err)
		return err
	}

	return nil
}

func (u *UDCClient) getKey(time time.Time) ([]byte, error) {
	key, err := u.otp.Generate(u.cfg.OTP.Seed, u.cfg.OTP.Period, time)
	if err != nil {
		return nil, err
	}

	return key, nil
}
