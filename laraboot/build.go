// Package laraboot .
package laraboot

import (
	_ "embed" // required
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
	"github.com/laraboot-io/shared"
	"github.com/paketo-buildpacks/packit"
	"gopkg.in/yaml.v2"
)

//go:embed commander.yml
var commanderYml string

// Build .
func Build(logger shared.LogEmitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		thisLayer, blueprintGenErr := context.Layers.Get("laraboot-commander")
		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		_, errDecoding := shared.NewFromFile(filepath.Join(context.WorkingDir, "laraboot.json"))
		if errDecoding != nil {
			fmt.Printf("	--> An error occurred while parsing laraboot file: '%s'", errDecoding)
		}

		var m struct {
			Commander struct {
				WorkingDir string `yaml:"directory"`
				Commands   []struct {
					Name string `yaml:"name"`
					Run  string `yaml:"run"`
				} `yaml:"commands"`
				Git struct {
					Enabled bool `yaml:"enabled"`
					Commit  bool `yaml:"commit"`
				} `yaml:"git"`
				Clean bool `yaml:"cleanup"`
			} `yaml:"laraboot-commander"`
		}

		unmarshallErr := yaml.Unmarshal([]byte(commanderYml), &m)
		if unmarshallErr != nil {
			return packit.BuildResult{}, unmarshallErr
		}

		commandsLen := len(m.Commander.Commands)
		for k, v := range m.Commander.Commands {
			fileName := fmt.Sprintf("%s/command-%d.sh", thisLayer.Path, k)
			body := fmt.Sprintf("#!/usr/bin/env bash \n %s", v.Run)
			logger.Process("Running command [%d/%d]: %s", k+1, commandsLen, v.Name)
			_, err := os.Create(fileName)
			if err != nil {
				return packit.BuildResult{}, err
			}
			writeFile, err := script.Echo(body).WriteFile(fileName)
			fmt.Println(writeFile)
			if err != nil {
				return packit.BuildResult{}, err
			}
			p := script.Exec(fmt.Sprintf("bash -c '%s'", fileName))
			output, _ := p.String()
			fmt.Println(output)
			if p.ExitStatus() != 0 {
				return packit.BuildResult{}, errors.New("build failed: command exited with a non-zero status")
			}
		}
		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
