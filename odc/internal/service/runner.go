package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mint8846/Traversal-Learning/odc/internal/config"
)

type RunnerService struct {
	cfg *config.Config
}

func NewRunnerService(cfg *config.Config) *RunnerService {
	return &RunnerService{cfg: cfg}
}

func (r *RunnerService) ExecuteModel(modelPath string) error {
	// 1. Load docker image from tar
	imageName, err := r.loadContainerImage(modelPath)
	if err != nil {
		return err
	}
	log.Printf("ExecuteModel: imageName: %s", imageName)

	if err = r.runContainer(); err != nil {
		return err
	}
	return nil
}

func (r *RunnerService) CheckResultData(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if os.IsNotExist(err) {
		return fmt.Errorf("CheckResultData: result file not exist(%s)", path)
	}
	log.Printf("CheckResultData: file size(%d)", info.Size())
	return nil
}

func (r *RunnerService) loadContainerImage(modelPath string) (string, error) {
	//log.Printf("loadContainerImage: Loading docker image from: %s", r.cfg.ModelPath)
	cmd := exec.Command("docker", "load", "-i", modelPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("loadContainerImage: docker load failed: %v", err)
	}

	// ex) "Loaded image: nginx:latest"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Loaded image:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				imageName := strings.TrimSpace(parts[1]) + ":" + parts[2]
				return imageName, nil
			}
		}
	}

	return "", fmt.Errorf("could not find loaded image name")
}

func (r *RunnerService) runContainer() error {
	if err := exec.Command("chmod", "+x", r.cfg.ModelScript).Run(); err != nil {
		return fmt.Errorf("runContainer: chmod failed: %v", err)
	}
	cmd := exec.Command(r.cfg.ModelScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runContainer: failed: %v, output: %s", err, string(output))
	}

	log.Printf("runContainer success(%s)", output)
	return nil
}
