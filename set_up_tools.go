package ok

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/bobg/errors"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

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
			CommandPath:    commandPath,
			CommandPathErr: commandPathErr,
			FilePaths:      filePaths,
			FilePathsErr:   filePathsErr,
		})
		o.toolsLock.Unlock()
		return nil
	}
}
