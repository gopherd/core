package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	name, version, branch, commit, datetime string
)

func Name() string {
	if name != "" {
		return name
	}
	s := filepath.Base(os.Args[0])
	if strings.HasSuffix(s, ".exe") {
		s = strings.TrimSuffix(s, ".exe")
	}
	return s
}

func Version() string {
	return fmt.Sprintf("%s %s(%s: %s) built at %s by %s", Name(), version, branch, commit, datetime, runtime.Version())
}

func Print() {
	fmt.Println(Version())
}
