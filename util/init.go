package util

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func InitFile(filename string, contents []byte) error {
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("%q already exists", filename)
	} else if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to stat file")
	}

	if err := os.WriteFile(filename, contents, 0666); err != nil {
		return errors.Wrap(err, "failed to write file")
	}

	fmt.Printf("created %q\n", filename)
	return nil
}
