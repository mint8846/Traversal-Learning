package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mint8846/Traversal-Learning/odc/internal/utils"
)

type Config struct {
	ID              string
	SessionKey      string
	NFSPath         string
	ResultFile      string
	UDC             UDC
	OTP             OTP
	ModelScript     string
	ModelOutputPath string
}

type UDC struct {
	Url               string
	ID                string
	HTTPClientTimeout time.Duration
}

type OTP struct {
	Seed   string
	Period uint
}

func Load() (*Config, error) {
	seed := getEnv("OTP_SEED", uuid.NewString())
	log.Printf("config: seed key(%s)", seed)
	sessionKey := utils.HashB64([]byte(seed))
	log.Printf("config: session key(%s)", sessionKey)

	return &Config{
		SessionKey: sessionKey,
		ResultFile: "/tmp/odc/result.txt",
		ID:         getEnv("ODC_ID", "odc_id"),
		NFSPath:    getEnv("MOUNT_PATH", "/tmp/nfs"),
		UDC: UDC{
			Url: getEnv("UDC_SERVER", "http://udc:8080"),
			ID:  getEnv("UDC_ID", ""),
			// default 1 hour
			HTTPClientTimeout: time.Duration(getEnvInt("UDC_HTTP_TIMEOUT", 3600)),
		},
		OTP: OTP{
			Seed: seed,
			// default 180 sec
			Period: uint(getEnvInt("OTP_PERIOD", 180)),
		},
		ModelScript:     getEnv("MODEL_SCRIPT", "/model/script.sh"),
		ModelOutputPath: getEnv("MODEL_OUTPUT_FILE_PATH", "/tmp/output"),
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
