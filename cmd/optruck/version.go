package optruck

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var Version = "(dev)" // set in build with --ldflags "-X cmd/optruck/version.Version=v0.1.0"
type VersionFlag bool

func (v VersionFlag) BeforeApply(ctx *kong.Context) error {
	fmt.Printf("optruck version %s\n", Version)
	os.Exit(0)
	return nil
}
