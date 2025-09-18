package service

import "github.com/mint8846/Traversal-Learning/udc/internal/config"

type Container struct {
	File *FileService
	OTP  *OTPService
	NFS  *NFSService
	Cfg  *config.Config
}

var Default *Container

func Initialize(cfg *config.Config) {
	Default = &Container{
		Cfg:  cfg,
		File: &FileService{cfg: cfg},
		OTP:  &OTPService{cfg: cfg},
		NFS:  &NFSService{cfg: cfg},
	}

	Default.NFS.SetupSignalHandler()
}
