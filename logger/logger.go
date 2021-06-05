package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	Ok    = newlineReplacerLogger("[ok] ")
	Debug = newlineReplacerLogger("[ok.debug] ")

	toolLoggers = make(map[string]*log.Logger)
)

func Tool(toolName string) *log.Logger {
	if logger, loggerExists := toolLoggers[toolName]; loggerExists {
		return logger
	}

	logger := newlineReplacerLogger(fmt.Sprintf("[ok %s] ", toolName))
	toolLoggers[toolName] = logger
	return logger
}

func newlineReplacerLogger(prefix string) *log.Logger {
	return log.New(newReplacerWriter(os.Stdout, "\n", fmt.Sprintf("\n%s", prefix)), prefix, 0)
}

func newReplacerWriter(w io.Writer, old, new string) replacerWriter {
	return replacerWriter{
		old:    old,
		new:    new,
		writer: w,
	}
}

type replacerWriter struct {
	old    string
	new    string
	writer io.Writer
}

func (w replacerWriter) Write(p []byte) (int, error) {
	s := string(p)
	n := strings.Count(s, "\n")
	if strings.HasSuffix(s, "\n") {
		n--
	}

	return w.writer.Write([]byte(strings.Replace(s, w.old, w.new, n)))
}
