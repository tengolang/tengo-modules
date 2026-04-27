package http

import (
	"bytes"
	"io"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/tengolang/tengo/v3"
)

var defaultClient = &gohttp.Client{Timeout: 30 * time.Second}

// Module is the Tengo "http" module.
//
//	http := import("http")
//	http.get(url string) => Response | error
//	http.post(url string, content_type string, body bytes) => Response | error
//	http.request(method string, url string, headers map, body bytes) => Response | error
//
// Response map keys: status_code int, status string, body bytes, headers map
var Module = map[string]tengo.Object{
	"get": &tengo.UserFunction{
		Name: "get",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 1); err != nil {
				return nil, err
			}
			url, err := tengo.ArgString(args, 0, "url")
			if err != nil {
				return nil, err
			}
			resp, e := defaultClient.Get(url) //nolint:gosec
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return responseObject(resp)
		},
	},

	"post": &tengo.UserFunction{
		Name: "post",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 3); err != nil {
				return nil, err
			}
			url, err := tengo.ArgString(args, 0, "url")
			if err != nil {
				return nil, err
			}
			contentType, err := tengo.ArgString(args, 1, "content_type")
			if err != nil {
				return nil, err
			}
			body, err := tengo.ArgBytes(args, 2, "body")
			if err != nil {
				return nil, err
			}
			resp, e := defaultClient.Post(url, contentType, bytes.NewReader(body)) //nolint:gosec
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return responseObject(resp)
		},
	},

	"request": &tengo.UserFunction{
		Name: "request",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if err := tengo.ArgCount(args, 4); err != nil {
				return nil, err
			}
			method, err := tengo.ArgString(args, 0, "method")
			if err != nil {
				return nil, err
			}
			url, err := tengo.ArgString(args, 1, "url")
			if err != nil {
				return nil, err
			}
			body, err := tengo.ArgBytes(args, 3, "body")
			if err != nil {
				return nil, err
			}

			req, e := gohttp.NewRequest(strings.ToUpper(method), url, bytes.NewReader(body)) //nolint:gosec
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}

			// headers arg — map or undefined
			switch h := args[2].(type) {
			case *tengo.Map:
				for k, v := range h.Value {
					s, ok := tengo.ToString(v)
					if ok {
						req.Header.Set(k, s)
					}
				}
			case *tengo.ImmutableMap:
				for k, v := range h.Value {
					if strings.HasPrefix(k, "__") {
						continue
					}
					s, ok := tengo.ToString(v)
					if ok {
						req.Header.Set(k, s)
					}
				}
			case *tengo.Undefined:
				// no headers
			default:
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "headers",
					Expected: "map or undefined",
					Found:    args[2].TypeName(),
				}
			}

			resp, e := defaultClient.Do(req)
			if e != nil {
				return &tengo.Error{Value: &tengo.String{Value: e.Error()}}, nil
			}
			return responseObject(resp)
		},
	},
}

func responseObject(resp *gohttp.Response) (tengo.Object, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
	}

	headers := make(map[string]tengo.Object, len(resp.Header))
	for k, vs := range resp.Header {
		headers[k] = &tengo.String{Value: strings.Join(vs, ", ")}
	}

	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"status_code": tengo.Int{Value: int64(resp.StatusCode)},
			"status":      &tengo.String{Value: resp.Status},
			"body":        &tengo.Bytes{Value: body},
			"headers":     &tengo.ImmutableMap{Value: headers},
		},
	}, nil
}
