package laraboot_commander

import (
	"bytes"
	"path/filepath"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

import (
	"testing"
)

func TestUnitGoBuild(t *testing.T) {
	suite := spec.New("go-build", spec.Report(report.Terminal{}))
	suite("Build", testBuild, spec.Sequential())
	suite.Run(t)
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	bash := pexec.NewExecutable("bash")
	buffer := bytes.NewBuffer(nil)
	err := bash.Execute(pexec.Execution{
		Args:   []string{"-c", "scripts/package.sh --version 1.2.3"},
		Stdout: buffer,
		Stderr: buffer,
	})
	Expect(err).NotTo(HaveOccurred(), buffer.String)

	_, err = filepath.Abs("../dist/laraboot-rector.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite.Run(t)
}