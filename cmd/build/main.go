package main

import (
	"os"

	"github.com/laraboot-io/shared"
	"github.com/paketo-buildpacks/packit"
	"laraboot-buildpacks/laraboot-commander/laraboot"
)

func main() {
	logEmitter := shared.NewLogEmitter(os.Stdout)
	packit.Build(laraboot.Build(logEmitter))
}
