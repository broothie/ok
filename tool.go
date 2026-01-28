package ok

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bobg/errors"
	"github.com/broothie/ok/table"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type Tool struct {
	Name        string
	CommandName string
	FileGlobs   []string
	ProcessFile func(ctx context.Context, filePath string) ([]Task, error)
}

type toolInfo struct {
	Tool
	commandPath    string
	commandPathErr error
	filePaths      []string
	filePathsErr   error
}

func (i toolInfo) err() error {
	return errors.Join(i.commandPathErr, i.filePathsErr)
}

func (o *Ok) ListTools(w io.Writer) error {
	rows := lo.Map(o.tools, func(tl toolInfo, _ int) []string {
		executable := fmt.Sprintf("✔ %s", tl.commandPath)
		if tl.commandPathErr != nil {
			executable = fmt.Sprintf("✘ %v", tl.commandPathErr)
		}

		files := strings.Join(tl.filePaths, ",")
		if tl.filePathsErr != nil {
			files = fmt.Sprintf("✘ %v", tl.filePathsErr)
		}

		return []string{tl.Name, executable, files}
	})

	sort.Slice(rows, func(i, j int) bool { return rows[i][0] < rows[j][0] })

	rows = append([][]string{{"TOOL", "EXECUTABLE", "FILES"}}, rows...)
	return errors.Wrap(table.Write(w, rows), "writing table")
}

func (o *Ok) SetUpTools(ctx context.Context, tls []Tool) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, tl := range tls {
		group.Go(o.processTool(tl))
	}

	return group.Wait()
}

func (o *Ok) processTool(tl Tool) func() error {
	return func() error {
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
		o.tools = append(o.tools, toolInfo{
			Tool:           tl,
			commandPath:    commandPath,
			commandPathErr: commandPathErr,
			filePaths:      filePaths,
			filePathsErr:   filePathsErr,
		})
		o.toolsLock.Unlock()
		return nil
	}
}
