package main

import (
	"log"

	"github.com/mint8846/Traversal-Learning/odc/internal/client"
	"github.com/mint8846/Traversal-Learning/odc/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	udcClient := client.NewUDCClient(cfg)

	modelPath, err := udcClient.GetModel()
	if err != nil {
		log.Fatal("Failed GetModel", err)
		return
	}

	if err = udcClient.ExecuteModel(modelPath); err != nil {
		log.Fatal("Failed ExecuteModel", err)
		return
	}

	if err = udcClient.EncryptResult(); err != nil {
		log.Fatal("Failed EncryptResult", err)
		return
	}
}
