package utils

import (
	"fmt"
	"os"
)

type RgxConfig struct {
	ServerUrl            string
	ArtifactRegistryBase string
	ArtifactRegistryAuth string
	ShowProgress         bool
	PackagesDir          string
	DownloadDir          string
	RcFileDir            string
}

type NexusArtifact struct {
	Repo       string
	Group      string
	Artifact   string
	Version    string
	Classifier string
	Extension  string
}

type Checksum struct {
	Algorithm string
	Hash      string
}

type RuntimeConfig struct {
	Machine string
	OS      string
	Arch    string
	User    string
}

type Dict map[string]any

func (d Dict) GetDict(k string) Dict {
	if _, ok := d[k]; !ok {
		fmt.Printf("No dictionary entry found for: %s\n", k)
		fmt.Println("Please check the configuration file")
		os.Exit(1)
	}
	return d[k].(map[string]any)
}

func (d Dict) GetString(k, fallback string) string {
	if d[k] == nil {
		return fallback
	} else {
		return d[k].(string)
	}
}

func (d Dict) GetInt(k string, fallback int) int {
	if d[k] == nil {
		return fallback
	} else {
		return int(d[k].(int64))
	}
}

func (d Dict) GetBool(k string, fallback bool) bool {
	if d[k] == nil {
		return fallback
	} else {
		return d[k].(bool)
	}
}
