// Package laraboot .
package laraboot

import (
	_ "embed" // required
	"errors"
	"fmt"
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
		logger.Process("Reading laraboot.json")

		thisLayer, blueprintGenErr := context.Layers.Get("laravel-model")
		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		_, errDecoding := shared.NewFromFile(filepath.Join(context.WorkingDir, "laraboot.json"))

		if errDecoding != nil {
			fmt.Printf("	--> An error occurred while parsing laraboot file: '%s'", errDecoding)
		}

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

		unmarshallErr := yaml.Unmarshal([]byte(commanderYml), &m)

		if unmarshallErr != nil {
			fmt.Printf("unmarshallErr: %v", unmarshallErr)
			return packit.BuildResult{}, unmarshallErr
		}

		commandsLen := len(m.Commander.Commands)
		for k, v := range m.Commander.Commands {
			logger.Process("Running command [%d/%d]: %s", k+1, commandsLen, v)
			p := script.Exec(fmt.Sprintf("bash -c '%s'", v))
			output, _ := p.String()
			fmt.Println(output)
			var exit = p.ExitStatus()
			if exit != 0 {
				err1 := errors.New("build failed: command exited with a non-zero status")
				return packit.BuildResult{}, err1
			}
		}

		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
