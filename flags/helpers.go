// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
package flags

import (
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
)

var (
	OriginCommandHelpTemplate = `{{.Name}}{{if .Subcommands}} command{{end}}{{if .Flags}} [command options]{{end}} {{.ArgsUsage}}
{{if .Description}}{{.Description}}
{{end}}{{if .Subcommands}}
SUBCOMMANDS:
  {{range .Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}{{end}}{{if .Flags}}
OPTIONS:
{{range $.Flags}}   {{.}}
{{end}}
{{end}}`
)

// HelpData is a one shot struct to pass to the usage template
type HelpData struct {
	App        interface{}
	FlagGroups []FlagGroup
}

// FlagGroup is a collection of flags belonging to a single topic.
type FlagGroup struct {
	Name  string
	Flags []cli.Flag
}

// byCategory sorts an array of FlagGroup by Name in the order
// defined in AppHelpFlagGroups.
type ByCategory []FlagGroup

func (a ByCategory) Len() int      { return len(a) }
func (a ByCategory) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCategory) Less(i, j int) bool {
	iCat, jCat := a[i].Name, a[j].Name
	iIdx, jIdx := len(a), len(a) // ensure non categorized flags come last

	for i, group := range a {
		if iCat == group.Name {
			iIdx = i
		}
		if jCat == group.Name {
			jIdx = i
		}
	}

	return iIdx < jIdx
}

// NewApp creates an app with sane defaults.
func NewApp(gitCommit, gitDate, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	app.Email = ""
	app.Version = params.VersionWithCommit(gitCommit, gitDate)
	app.Usage = usage
	return app
}
