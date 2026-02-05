package snapshot

import (
	"errors"
	"os"
	"path"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Context struct {
	Name     string
	Filename string
}

func (c *Context) Load() (*Snapshot, error) {
	c.Filename = path.Join(
		"testdata", "__snapshots__",
		strings.ReplaceAll(cases.Lower(language.Und).String(c.Name), " ", "_")+".txtar",
	)

	if strings.ToUpper(updateSnapshots) == "ALL" || strings.Contains(updateSnapshots, c.Name) {
		return &Snapshot{ctx: c}, nil
	}

	snapshotRaw, err := os.ReadFile(c.Filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Snapshot{ctx: c}, nil
		}
		return nil, err
	}

	return &Snapshot{
		ctx:   c,
		files: slices.Collect(FilesSeq(snapshotRaw)),
	}, nil
}
