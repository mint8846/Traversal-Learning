package service

import (
	"fmt"
	"os/exec"

	"github.com/mint8846/Traversal-Learning/odc/internal/config"
)

type NFSService struct {
	cfg *config.Config
}

func NewNFSService(cfg *config.Config) *NFSService {
	return &NFSService{
		cfg: cfg,
	}
}

func (n *NFSService) Connect(url string) error {
	commands := [][]string{
		{"mkdir", "-p", n.cfg.NFSPath},
		{"mount", "-t", "nfs", url, n.cfg.NFSPath},
	}

	for _, cmd := range commands {
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			return fmt.Errorf("failed to execute command %s: %s", cmd, err)
		}
	}

	return nil
}

func (n *NFSService) GetPath() string {
	return n.cfg.NFSPath
}
