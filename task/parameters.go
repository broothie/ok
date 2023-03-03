package task

import (
	"fmt"
	"strings"
)

type Parameters []Parameter

func (p Parameters) String() string {
	var fields []string
	for _, param := range p {
		if param.IsRequired() {
			fields = append(fields, fmt.Sprintf("<%s>", param.Name))
		} else {
			fields = append(fields, fmt.Sprintf("--%s=%s", param.Name, *param.Default))
		}
	}

	return strings.Join(fields, " ")
}
