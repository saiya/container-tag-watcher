package config

import (
	"fmt"
	"strings"
)

type ConfigParseError struct {
	Details []ConfigParseErrorDetail
}

// Error implements Error interface
func (e *ConfigParseError) Error() string {
	buff := make([]string, len(e.Details))
	for i, d := range e.Details {
		buff[i] = d.String()
	}
	return strings.Join(buff, "\n")
}

func mergeParseError(into *ConfigParseError, e *ConfigParseError) *ConfigParseError {
	if into == nil {
		return e
	} else if e == nil {
		return into
	} else {
		into.Details = append(into.Details, e.Details...)
		return into
	}
}

func appendParseError(into *ConfigParseError, detail *ConfigParseErrorDetail) *ConfigParseError {
	if detail == nil {
		return into
	}
	if into == nil {
		into = &ConfigParseError{Details: make([]ConfigParseErrorDetail, 0, 1)}
	}
	into.Details = append(into.Details, *detail)
	return into
}

type ConfigParseErrorDetail struct {
	Location string
	Message  string
}

func (e *ConfigParseErrorDetail) String() string {
	return fmt.Sprintf("%s: %s", e.Location, e.Message)
}
