package optruck

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type VersionFlag bool

func (v VersionFlag) BeforeApply(ctx *kong.Context) error {
	tmpl := `optruck version %s
  commit: %s
  date: %s
  built by: %s
`
	fmt.Printf(tmpl, buildInfo.Version, buildInfo.Commit, buildInfo.Date, buildInfo.BuiltBy)
	os.Exit(0)
	return nil
}
