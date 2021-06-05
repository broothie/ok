package ok

import (
	"os"

	"github.com/broothie/ok/cli"
	"github.com/broothie/ok/logger"
	"github.com/broothie/ok/tools"
	"github.com/pkg/errors"
)

func (ok *Ok) HandleOptions() (taskName string, halt bool, err error) {
	taskName, options, err := ok.Parser.ParseOptions()
	if err != nil {
		return "", true, err
	}

	ok.Options = options
	if options.Help {
		return "", true, cli.PrintHelp(os.Stdout, Version())
	} else if options.Version {
		return "", true, cli.PrintVersion(os.Stdout, Version())
	} else if options.ListTools {
		tools.List()
		return "", true, nil
	} else if options.Init != "" {
		return "", true, tools.Init(ok.Options.Init)
	} else if taskName == "" {
		if mountErrors := ok.Mount(); len(mountErrors) != 0 {
			for toolName, err := range mountErrors {
				logger.Tool(toolName).Print(errors.Wrapf(err, "failed to mount tool '%s'", toolName))
			}
		}

		return "", true, ok.List()
	}

	return taskName, false, err
}
