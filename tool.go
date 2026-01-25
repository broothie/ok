package ok

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
	"github.com/broothie/ok/tool"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type Tool struct {
	tool.Tool
	CommandPath    string
	CommandPathErr error
	FilePaths      []string
	FilePathsErr   error
}

func (o *Ok) SetUpTools(ctx context.Context, tls []tool.Tool) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, tl := range tls {
		group.Go(func() error {
			commandPath, commandPathErr := exec.LookPath(tl.CommandName)

			var filePaths []string
			var filePathErrs []error
			for _, fileGlob := range tl.FileGlobs {
				globFilePaths, err := doublestar.Glob(os.DirFS("."), fileGlob)

				filePaths = append(filePaths, globFilePaths...)
				filePathErrs = append(filePathErrs, err)
			}

			filePaths = lo.UniqBy(filePaths, func(filePath string) string { return strings.ToLower(filePath) })
			filePathsErr := errors.Join(filePathErrs...)

			o.toolsLock.Lock()
			o.tools = append(o.tools, Tool{
				Tool:           tl,
				CommandPath:    commandPath,
				CommandPathErr: commandPathErr,
				FilePaths:      filePaths,
				FilePathsErr:   filePathsErr,
			})
			o.toolsLock.Unlock()
			return nil
		})
	}

	return group.Wait()
}

func (o *Ok) PrintTools(w io.Writer) error {
	rows := [][]string{{"TOOL", "EXECUTABLE", "FILES"}}

	for _, tl := range o.tools {
		executable := tl.CommandPath
		if tl.CommandPathErr != nil {
			executable = tl.CommandPathErr.Error()
		}

		files := strings.Join(tl.FilePaths, ",")
		if tl.FilePathsErr != nil {
			files = tl.FilePathsErr.Error()
		}

		rows = append(rows, []string{tl.Name, executable, files})
	}

	sort.Slice(rows[1:], func(i, j int) bool { return rows[1:][i][0] < rows[1:][j][0] })
	return errors.Wrap(table.Write(w, rows), "writing table")
}
