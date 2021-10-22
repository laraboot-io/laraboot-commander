// Package larabootcommander .
package larabootcommander

import (
	_ "embed" // required
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
	"github.com/laraboot-io/shared"
	"github.com/paketo-buildpacks/packit"
	"gopkg.in/yaml.v2"
)

//go:embed commander.yml
var commanderYml string

// Commander command structure.
type Commander struct {
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
}

// Build .
func Build(logger shared.LogEmitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("Building %s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		thisLayer, layerErr := context.Layers.Get("laraboot-commander")
		if layerErr != nil {
			fmt.Printf("	--> An error occurred getting or creating layer: %s", layerErr)
			return packit.BuildResult{}, layerErr
		}
		_, errDecoding := shared.NewFromFile(filepath.Join(context.WorkingDir, "laraboot.json"))
		if errDecoding != nil {
			fmt.Printf("	--> An error occurred while parsing laraboot file: '%s'", errDecoding)
		}
		var m struct {
			Commander `yaml:"laraboot-commander"`
		}
		unmarshallErr := yaml.Unmarshal([]byte(commanderYml), &m)
		if unmarshallErr != nil {
			fmt.Printf("	--> An error occurred reading YML: %s", unmarshallErr)
			return packit.BuildResult{}, unmarshallErr
		}
		if _, err := os.Stat(thisLayer.Path); os.IsNotExist(err) {
			err := os.Mkdir(thisLayer.Path, 0600) //nolint:gomnd //ignore
			if err != nil {
				return packit.BuildResult{}, err
			}
		}
		logger.Subprocess("Creating sandbox")
		commandsLen := len(m.Commander.Commands)
		commandsDir := context.WorkingDir
		for k, v := range m.Commander.Commands {
			fileName := fmt.Sprintf("%s/command-%d.sh", commandsDir, k)
			body := fmt.Sprintf("#!/usr/bin/env bash \n %s", v.Run)
			logger.Action("Running command [%d/%d]: %s", k+1, commandsLen, v.Name)
			err := ioutil.WriteFile(fileName, []byte(body), 0777) //nolint:gosec // we need exe permissions
			if err != nil {
				fmt.Printf("	--> An error occurred writing file: %s", err)
				return packit.BuildResult{}, err
			}
			p := script.Exec(fmt.Sprintf("bash -c '%s'", fileName))
			output, err := p.String()
			logger.Detail(output)
			if err != nil {
				fmt.Printf("	--> An error occurred : %s", err)
				return packit.BuildResult{}, err
			}
			if p.ExitStatus() != 0 {
				return packit.BuildResult{}, errors.New("build failed: command exited with a non-zero status")
			}
			logger.Break()
		}
		logger.Break()
		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
