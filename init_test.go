// Package Larabootcommander
package larabootcommander

import (
	"bytes"
	"path/filepath"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/pexec"
)

import (
	"testing"
)

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect //nolint:govet //ignore

	bash := pexec.NewExecutable("bash")
	buffer := bytes.NewBuffer(nil)
	err := bash.Execute(pexec.Execution{
		Args:   []string{"-c", "scripts/package.sh --version 1.2.3"},
		Stdout: buffer,
		Stderr: buffer,
	})
	Expect(err).NotTo(HaveOccurred(), buffer.String)

	_, err = filepath.Abs("dist/laraboot-rector.cnb")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)
}
