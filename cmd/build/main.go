package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	laraboot "laraboot-buildpacks/laraboot-commander/laraboot"
	"log"
	"os"
)

func init() {
	log.Println("::init")
}

func main() {

	logEmitter := laraboot.NewLogEmitter(os.Stdout)

	packit.Build(laraboot.Build(
		logEmitter,
		chronos.DefaultClock))
}
