package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Szubie/dicomweb-go/dicomweb"
)

func UploadStudy(path string, serverUrl string) error {

	// Get URL from config if not provided
	if serverUrl == "" {
		config, err := LoadServerConfig()
		if err != nil {
			return err
		}
		serverUrl = config.STOWEndpoint
	}

	// Create client with the required options for STOW
	client := dicomweb.NewClient(
		dicomweb.ClientOption{
			STOWEndpoint: serverUrl,
		},
	)

	// Prepare multi-part upload
	var parts [][]byte

	// Read file
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	var bytesBuf bytes.Buffer
	_, err = io.Copy(&bytesBuf, f)
	if err != nil {
		return err
	}
	parts = append(parts, bytesBuf.Bytes())

	// Make request
	stow := dicomweb.STOWRequest{
		StudyInstanceUID: "",
		Parts:            parts,
	}
	_, err = client.Store(stow)
	if err != nil {
		return err
	}

	return nil
}

func DownloadStudy(studyId string, destPath string, serverUrl string) error {
	// Get URL from config if not provided
	if serverUrl == "" {
		config, err := LoadServerConfig()
		if err != nil {
			return err
		}
		serverUrl = config.STOWEndpoint
	}
	// Use studyId as destPath if destPath not provided
	if destPath == "" {
		destPath = fmt.Sprintf("%s.dcm", studyId)
	}

	// Create client with the required options for STOW
	client := dicomweb.NewClient(
		dicomweb.ClientOption{
			WADOEndpoint: serverUrl,
		},
	)

	wado := dicomweb.WADORequest{
		Type:             dicomweb.StudyRaw,
		StudyInstanceUID: studyId,
	}
	parts, err := client.Retrieve(wado)
	if err != nil {
		return err
	}

	if len(parts) == 1 {
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		_, err = f.Write(parts[0])
		if err != nil {
			return err
		}
	} else {
		for i, p := range parts {
			f, err := os.Create(fmt.Sprintf("%s_%d.dcm", destPath, i))
			if err != nil {
				return err
			}
			_, err = f.Write(p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func QueryStudy(studyId string, serverUrl string, destPath string) error {
	// Get URL from config if not provided
	if serverUrl == "" {
		config, err := LoadServerConfig()
		if err != nil {
			return err
		}
		serverUrl = config.QIDOEndpoint
	}
	// Use studyId as destPath if destPath not provided
	if destPath == "" {
		destPath = fmt.Sprintf("%s.json", studyId)
	}

	// Create client with the required options for STOW
	client := dicomweb.NewClient(
		dicomweb.ClientOption{
			QIDOEndpoint: serverUrl,
		},
	)

	qido := dicomweb.QIDORequest{
		Type:             dicomweb.Study,
		StudyInstanceUID: studyId,
	}
	resp, err := client.Query(qido)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	// Serialize response into JSON string
	respJsonBytes, err := json.MarshalIndent(resp[0], "", "  ")
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(respJsonBytes)

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	io.Copy(f, buffer)

	return nil
}
