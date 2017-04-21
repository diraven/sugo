package sugo

import (
	"os/exec"
	"fmt"
)

const versionMajor int = 0
const versionMinor int = 0
const versionBuild int = 12
var versionHash string
var Version string

func init() {
	result, _ := exec.Command("git", "rev-parse", " --short HEAD").Output()
	versionHash = fmt.Sprintf("%s", result)

	Version = fmt.Sprintf("%d.%d.%d.%s", versionMajor, versionMinor, versionBuild, versionHash)
}