package service

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mint8846/Traversal-Learning/udc/internal/config"
)

type NFSService struct {
	cfg *config.Config
}

func NewNFSService(cfg *config.Config) *NFSService {
	return &NFSService{
		cfg: cfg,
	}
}

func (n *NFSService) Connect(url string) (string, error) {
	commands := [][]string{
		{"mkdir", "-p", n.cfg.NFSPath},
		{"mount", "-t", "nfs", url, n.cfg.NFSPath},
	}

	for _, cmd := range commands {
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			return "", fmt.Errorf("failed to execute command %s: %s", cmd, err)
		}
	}
	return n.cfg.NFSPath, nil
}

func (n *NFSService) GenerateNFSUrl() (string, error) {
	host := n.cfg.ServerHost

	if host == "" {
		serverIP, err := n.getServerIPAddress()
		if err != nil {
			log.Printf("GenerateNFSUrl: get public IP fail %v", err)
			return "", err
		}

		host = serverIP
	}

	return fmt.Sprintf("%s:/%s/", host, n.cfg.HostName), nil
}

func (n *NFSService) SetupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal, cleaning NFS mount path...")
		n.cleanup()

		log.Println("Cleanup completed, exiting...")
		os.Exit(0)
	}()
}

func (n *NFSService) getServerIPAddress() (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func (n *NFSService) cleanup() {
	if err := os.RemoveAll(filepath.Join(n.cfg.NFSPath, n.cfg.HostName)); err != nil {
		log.Printf("Warning: Failed to clean mount path files: %v", err)
	} else {
		log.Printf("Mount path files cleaned: %s", n.cfg.NFSPath)
	}
}
