// configurations.go
package configurations

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

const jsonPropertyFileName = "config.json"

type GitConfig struct {
	Enabled      bool   `json:"enabled"`
	RepoURL      string `json:"repo-url"`
	RepoUser     string `json:"repo-user"`
	RepoPassword string `json:"repo-password"`
}

type ServerConfig struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip-address"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
}

type ClusterConfig struct {
	Name    string         `json:"name"`
	Servers []ServerConfig `json:"servers"`
}

type FrameworkConfig struct {
	Name string `json:"name"`
}

type LanguageConfig struct {
	Name      string          `json:"name"`
	Version   string          `json:"version"`
	Framework FrameworkConfig `json:"framework"`
}

type ArtifactConfig struct {
	Language LanguageConfig `json:"language"`
}

type Config struct {
	AppJetURL string `json:"appJetURL"`
	Plugins   struct {
		Git GitConfig `json:"git"`
	} `json:"plugins"`
	Cluster  ClusterConfig  `json:"cluster"`
	Artifact ArtifactConfig `json:"artifact"`
}

var AppConfig Config

func LoadConfig() {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		panic(err)
	}
}

func getConfigFilePath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get current file path")
	}

	return filepath.Join(filepath.Dir(filename), "", jsonPropertyFileName), nil
}

func CreateConfigFile(appJetURL string, gitEnabled bool, gitRepoURL string, gitRepoUser string, gitRepoPassword string,
	serverName string, serverIPAddress string, serverPort int, serverUser string, serverPassword string,
	clusterName string, languageName string, languageVersion string, frameworkName string) error {

	// Create a new Config instance
	newConfig := Config{
		AppJetURL: appJetURL,
		Plugins: struct {
			Git GitConfig `json:"git"`
		}{
			Git: GitConfig{
				Enabled:      gitEnabled,
				RepoURL:      gitRepoURL,
				RepoUser:     gitRepoUser,
				RepoPassword: gitRepoPassword,
			},
		},
		Cluster: ClusterConfig{
			Name: clusterName,
			Servers: []ServerConfig{
				{
					Name:      serverName,
					IPAddress: serverIPAddress,
					Port:      serverPort,
					User:      serverUser,
					Password:  serverPassword,
				},
			},
		},
		Artifact: ArtifactConfig{
			Language: LanguageConfig{
				Name: languageName,
				Version: languageVersion,
				Framework: FrameworkConfig{
					Name: frameworkName,
				},
			},
		},
	}

	// Get the path to the configuration file
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Create or overwrite the configuration file
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the configuration and write it to the file
	encoder := json.NewEncoder(file)
	err = encoder.Encode(&newConfig)
	if err != nil {
		return err
	}

	return nil
}