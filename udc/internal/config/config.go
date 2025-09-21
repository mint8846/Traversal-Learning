package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ID          string
	ResultDir   string
	ServerHost  string
	ODC         ODC
	OTPPeriod   uint
	ModelPath   string
	ModelScript string
	NFSPath     string
	Port        string
}

type ODC struct {
	IP string
	ID string
}

func Load() (*Config, error) {
	return &Config{
		ID:         getEnv("UDC_ID", "udc_id"),
		ResultDir:  "/tmp/result/",
		ServerHost: getEnv("UDC_HOST", ""),
		ODC: ODC{
			IP: getEnv("ODC_IP", ""),
			ID: getEnv("ODC_ID", ""),
		},
		// default 180 sec
		OTPPeriod:   uint(getEnvInt("OTP_PERIOD", 180)),
		ModelPath:   getEnv("MODEL_PATH", "/model/model.tar"),
		ModelScript: getEnv("MODEL_SCRIPT", "/model/script.sh"),
		// In this example, default environment variable values are defined in the Dockerfile
		NFSPath: getEnv("NFS_EXPORT_PATH", "/nfs/share"),
		Port:    ":" + getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := getEnv(key, strconv.Itoa(defaultValue))
	if num, err := strconv.Atoi(value); err == nil {
		return num
	}

	log.Printf("config: (%s) invalid value(%s) use default value(%d)", key, value, defaultValue)
	return defaultValue
}
