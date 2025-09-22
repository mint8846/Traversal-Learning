package service

import "github.com/mint8846/Traversal-Learning/udc/internal/config"

type Container struct {
	File   *FileService
	NFS    *NFSService
	Cfg    *config.Config
	OTP    *OTPService
	Runner *RunnerService
}

var Default *Container

func Initialize(cfg *config.Config) {
	Default = &Container{
		Cfg:    cfg,
		File:   &FileService{cfg: cfg},
		NFS:    &NFSService{cfg: cfg},
		OTP:    &OTPService{},
		Runner: &RunnerService{cfg: cfg},
	}

	Default.NFS.SetupSignalHandler()
}
