package core

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type DICOMWebServer struct {
	QIDOEndpoint string
	WADOEndpoint string
	STOWEndpoint string
}

func LoadServerConfig() (config DICOMWebServer, err error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}

	configPath := filepath.Join(userHome, configDir, configFile)
	f, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	// Read config
	var byteBuffer bytes.Buffer
	_, err = io.Copy(&byteBuffer, f)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(byteBuffer.Bytes(), &config)
	return config, err
}

func SaveServerConfig(config DICOMWebServer) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(userHome, configDir, configFile)

	err = os.MkdirAll(filepath.Dir(configPath), 0777)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	// Write config
	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	_, err = f.Write(jsonBytes)
	return err
}

func ConfigExists() bool {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	configPath := filepath.Join(userHome, configDir, configFile)
	_, err = os.Stat(configPath)
	return err == nil
}
