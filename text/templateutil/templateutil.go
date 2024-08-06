package templateutil

import (
	"fmt"
	"html/template"
	"strings"
)

var (
	stringsFuncs = template.FuncMap{
		"contains":    strings.Contains,
		"count":       strings.Count,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"index":       strings.Index,
		"join":        strings.Join,
		"lastIndex":   strings.LastIndex,
		"repeat":      strings.Repeat,
		"replace":     strings.Replace,
		"replaceN":    strings.ReplaceAll,
		"split":       strings.Split,
		"toLower":     strings.ToLower,
		"toUpper":     strings.ToUpper,
		"toTitle":     strings.ToTitle,
		"toValidUTF8": strings.ToValidUTF8,
		"trim":        strings.Trim,
		"trimLeft":    strings.TrimLeft,
		"trimRight":   strings.TrimRight,
		"trimPrefix":  strings.TrimPrefix,
		"trimSuffix":  strings.TrimSuffix,
		"trimSpace":   strings.TrimSpace,
	}
	fmtFuncs = template.FuncMap{
		"format": fmt.Sprintf,
	}
)

func withFuncs(dst template.FuncMap, funcs ...template.FuncMap) template.FuncMap {
	for _, f := range funcs {
		for k, v := range f {
			dst[k] = v
		}
	}
	return dst
}

// DefaultFuncs returns the default template functions.
func DefaultFuncs() template.FuncMap {
	return withFuncs(template.FuncMap{}, stringsFuncs, fmtFuncs)
}

func DefaultTemplate(name string) *template.Template {
	return template.New(name).Funcs(DefaultFuncs())
}
