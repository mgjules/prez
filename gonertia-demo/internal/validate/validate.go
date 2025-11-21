package validate

import (
	"strings"
)

type Error struct {
	fields map[string]string
}

func (e *Error) Add(field, err string) *Error {
	if e.fields == nil {
		e.fields = make(map[string]string)
	}

	e.fields[field] = err
	return e
}

func (e *Error) HasErrors() bool {
	return len(e.fields) > 0
}

func (e *Error) Error() string {
	var err strings.Builder
	for k, v := range e.fields {
		err.WriteString(k)
		err.WriteString(": ")
		err.WriteString(v)
		err.WriteString("\n")
	}
	return err.String()
}
