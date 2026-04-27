package uuid

import (
	"github.com/google/uuid"
	"github.com/tengolang/tengo/v3"
)

// Module is the Tengo "uuid" module.
//
//	uuid := import("uuid")
//	uuid.v4()          -> string
//	uuid.v1()          -> string
//	uuid.parse(s)      -> string | error   (normalises to lowercase canonical form)
//	uuid.valid(s)      -> bool
//	uuid.nil           -> string           ("00000000-0000-0000-0000-000000000000")
var Module = map[string]tengo.Object{
	"nil": &tengo.String{Value: uuid.Nil.String()},

	"v4": &tengo.UserFunction{
		Name: "v4",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return nil, tengo.ErrWrongNumArguments
			}
			return &tengo.String{Value: uuid.New().String()}, nil
		},
	},

	"v1": &tengo.UserFunction{
		Name: "v1",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return nil, tengo.ErrWrongNumArguments
			}
			id, err := uuid.NewUUID()
			if err != nil {
				return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
			}
			return &tengo.String{Value: id.String()}, nil
		},
	},

	"parse": &tengo.UserFunction{
		Name: "parse",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			s, ok := tengo.ToString(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "s",
					Expected: "string",
					Found:    args[0].TypeName(),
				}
			}
			id, err := uuid.Parse(s)
			if err != nil {
				return &tengo.Error{Value: &tengo.String{Value: err.Error()}}, nil
			}
			return &tengo.String{Value: id.String()}, nil
		},
	},

	"valid": &tengo.UserFunction{
		Name: "valid",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			s, ok := tengo.ToString(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "s",
					Expected: "string",
					Found:    args[0].TypeName(),
				}
			}
			_, err := uuid.Parse(s)
			if err != nil {
				return tengo.FalseValue, nil
			}
			return tengo.TrueValue, nil
		},
	},
}
