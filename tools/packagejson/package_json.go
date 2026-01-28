package packagejson

import (
	"encoding/json"
	"os"

	"github.com/bobg/errors"
)

const packageJSONFileName = "package.json"

type schema struct {
	Scripts map[string]string `json:"scripts"`
}

func read(path string) (_ schema, err error) {
	file, err := os.Open(path)
	if err != nil {
		return schema{}, errors.Wrapf(err, "opening %q", path)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	var payload schema
	if err := json.NewDecoder(file).Decode(&payload); err != nil {
		return schema{}, errors.Wrapf(err, "parsing %q", path)
	}

	return payload, nil
}
