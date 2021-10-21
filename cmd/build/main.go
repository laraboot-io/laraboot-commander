package main

import (
	"os"

	"github.com/laraboot-io/shared"
	"github.com/paketo-buildpacks/packit"
	Larabootcommander "laraboot-buildpacks/laraboot-commander/laraboot"
)

func main() {
	logEmitter := shared.NewLogEmitter(os.Stdout)
	packit.Build(Larabootcommander.Build(logEmitter))
}
