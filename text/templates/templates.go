// Package templates provides utility functions for working with Go templates.
package templates

import (
	"bytes"
	"errors"
	"html/template"
	"regexp"
	"strings"
)

// New creates a new template with the default functions.
func New(name string) *template.Template {
	return template.New(name).Funcs(Funcs)
}

// Execute executes the default template with the given text and data.
func Execute(name, content string, data any, options ...string) (string, error) {
	var buf bytes.Buffer
	t := New(name)
	if len(options) > 0 {
		t = t.Option(options...)
	}
	if t, err := t.Parse(content); err != nil {
		return "", err
	} else if err := t.Execute(&buf, data); err != nil {
		return "", cleanTemplateError(err)
	} else {
		return buf.String(), nil
	}
}

func cleanTemplateError(err error) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	re := regexp.MustCompile(`^(template: .*?:\d+:\d+): executing "(.*?)" at .*?: (.*)$`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) == 4 && strings.Contains(matches[1], matches[2]) {
		cleaned := matches[1] + ": " + matches[3]
		return errors.New(cleaned)
	}
	return err
}
