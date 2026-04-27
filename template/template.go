package template

import (
	"bytes"
	htmltmpl "html/template"
	"strings"
	texttmpl "text/template"

	"github.com/tengolang/tengo/v3"
)

// Module is the Tengo "template" module.
//
//	tmpl := import("template")
//	tmpl.text(src string, data map) => string | error
//	tmpl.html(src string, data map) => string | error
var Module = map[string]tengo.Object{
	// text renders a text/template with the given data map.
	"text": &tengo.UserFunction{
		Name: "text",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			src, err := tengo.ArgString(args, 0, "src")
			if err != nil {
				return nil, err
			}
			data, err := toMap(args[1])
			if err != nil {
				return nil, err
			}
			t, e := texttmpl.New("").Parse(src)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			var buf bytes.Buffer
			if e = t.Execute(&buf, data); e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return &tengo.String{Value: buf.String()}, nil
		},
	},

	// html renders an html/template with the given data map.
	// Values are automatically HTML-escaped.
	"html": &tengo.UserFunction{
		Name: "html",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 2); err != nil {
				return nil, err
			}
			src, err := tengo.ArgString(args, 0, "src")
			if err != nil {
				return nil, err
			}
			data, err := toMap(args[1])
			if err != nil {
				return nil, err
			}
			t, e := htmltmpl.New("").Parse(src)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			var buf bytes.Buffer
			if e = t.Execute(&buf, data); e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return &tengo.String{Value: buf.String()}, nil
		},
	},

	// text_files parses and renders a named text/template from a glob pattern.
	"text_files": &tengo.UserFunction{
		Name: "text_files",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 3); err != nil {
				return nil, err
			}
			pattern, err := tengo.ArgString(args, 0, "pattern")
			if err != nil {
				return nil, err
			}
			name, err := tengo.ArgString(args, 1, "name")
			if err != nil {
				return nil, err
			}
			data, err := toMap(args[2])
			if err != nil {
				return nil, err
			}
			t, e := texttmpl.ParseGlob(pattern)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			var buf bytes.Buffer
			if e = t.ExecuteTemplate(&buf, name, data); e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return &tengo.String{Value: buf.String()}, nil
		},
	},

	// html_files parses and renders a named html/template from a glob pattern.
	"html_files": &tengo.UserFunction{
		Name: "html_files",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 3); err != nil {
				return nil, err
			}
			pattern, err := tengo.ArgString(args, 0, "pattern")
			if err != nil {
				return nil, err
			}
			name, err := tengo.ArgString(args, 1, "name")
			if err != nil {
				return nil, err
			}
			data, err := toMap(args[2])
			if err != nil {
				return nil, err
			}
			t, e := htmltmpl.ParseGlob(pattern)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			var buf bytes.Buffer
			if e = t.ExecuteTemplate(&buf, name, data); e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return &tengo.String{Value: buf.String()}, nil
		},
	},
}

// toMap converts a Tengo Map or ImmutableMap to map[string]any for template execution.
func toMap(o tengo.Object) (map[string]any, error) {
	var raw map[string]tengo.Object
	switch v := o.(type) {
	case *tengo.Map:
		raw = v.Value
	case *tengo.ImmutableMap:
		raw = v.Value
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "map",
			Found:    o.TypeName(),
		}
	}
	out := make(map[string]any, len(raw))
	for k, v := range raw {
		if strings.HasPrefix(k, "__") {
			continue
		}
		out[k] = tengoToAny(v)
	}
	return out, nil
}

func tengoToAny(o tengo.Object) any {
	switch v := o.(type) {
	case *tengo.String:
		return v.Value
	case tengo.Int:
		return v.Value
	case tengo.Float:
		return v.Value
	case tengo.Bool:
		return !v.IsFalsy()
	case *tengo.Bytes:
		return v.Value
	case *tengo.Array:
		out := make([]any, len(v.Value))
		for i, e := range v.Value {
			out[i] = tengoToAny(e)
		}
		return out
	case *tengo.Map:
		out := make(map[string]any, len(v.Value))
		for k, e := range v.Value {
			out[k] = tengoToAny(e)
		}
		return out
	case *tengo.ImmutableMap:
		out := make(map[string]any, len(v.Value))
		for k, e := range v.Value {
			if !strings.HasPrefix(k, "__") {
				out[k] = tengoToAny(e)
			}
		}
		return out
	default:
		return o.String()
	}
}
