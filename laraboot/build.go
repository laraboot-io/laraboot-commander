package laraboot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitfield/script"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var larabootStruct struct {
	Version   string `json:"version"`
	ProjectID string `json:"project_id"`
	Php       struct {
		Version string `json:"version"`
	} `json:"php"`
	Framework struct {
		Config struct {
			Overrides []struct {
				Key     string `json:"key"`
				Envs    string `json:"envs"`
				Default string `json:"default"`
			} `json:"overrides"`
		} `json:"config"`
		Auth struct {
			Stack string `json:"stack"`
		} `json:"auth"`
		Models []struct {
			Name    string `json:"name"`
			Columns []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"columns"`
		} `json:"models"`
	} `json:"Framework"`
}

func Build(logger LogEmitter, clock chronos.Clock) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Process("Reading laraboot.json")

		thisLayer, blueprintGenErr := context.Layers.Get("laravel-model")
		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		lfile, errReadingJsonFile := os.Open(filepath.Join(context.WorkingDir, "laraboot.json"))

		if errReadingJsonFile != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		errDecoding := json.NewDecoder(lfile).Decode(&larabootStruct)

		if errDecoding != nil {
			fmt.Printf("	--> An error ocurred while parsing laraboot file: '%s'", blueprintGenErr)
			return packit.BuildResult{}, errDecoding
		}

		// ---- Buildpack read process
		yamlFile, yamlError := ioutil.ReadFile(filepath.Join(context.WorkingDir, "commander.yml"))

		var m struct {
			Commander struct {
				WorkingDir string   `yaml:"directory"`
				Commands   []string `yaml:"commands"`
				Git        struct {
					Enabled bool `yaml:"enabled"`
					Commit  bool `yaml:"commit"`
				} `yaml:"git"`
				Clean bool `yaml:"cleanup"`
			} `yaml:"laraboot-commander"`
		}
		//_, blueprintGenErr = toml.DecodeReader(file, &m)
		if yamlError != nil {
			return packit.BuildResult{}, yamlError
		}

		unmarshallErr := yaml.Unmarshal(yamlFile, &m)

		if unmarshallErr != nil {
			fmt.Printf("unmarshallErr: %v", unmarshallErr)
			return packit.BuildResult{}, yamlError
		}

		commandsLen := len(m.Commander.Commands)
		for k, v := range m.Commander.Commands {
			logger.Process("Running command [%d/%d]: %s", k+1, commandsLen, v)
			p := script.Exec(fmt.Sprintf("bash -c '%s'", v))
			output, _ := p.String()
			fmt.Println(output)
			var exit int = p.ExitStatus()
			if exit != 0 {
				err1 := errors.New("Build failed: command exited with a non-zero status")
				return packit.BuildResult{}, err1
			}
		}

		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
