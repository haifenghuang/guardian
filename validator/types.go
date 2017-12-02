package validator

import (
	"bytes"

	"github.com/end-r/guardian/lexer"

	"github.com/end-r/guardian/ast"
)

// There are 5 first-class guardian types:
// Literal: int, string etc.
// Array: arrays[Type]
// NOTE: array = golang's slice, there is no golang array equivalent
// Map: map[Type]Type
// Func: func(Tuple)Tuple

// There are 2 second-class guardian types:
// Tuple: (Type...)
// Aliased: string -> Type

type Type interface {
	write(*bytes.Buffer)
	compare(Type) bool
	inherits(Type) bool
	implements(Type) bool
}

type BaseType int

const (
	Invalid BaseType = iota
	Unknown
	Bool
)

type StandardType struct {
	name string
}

var standards = map[BaseType]StandardType{
	Invalid: StandardType{"invalid"},
	Unknown: StandardType{"unknown"},
	Bool:    StandardType{"bool"},
}

type Array struct {
	Length   int
	Value    Type
	Variable bool
}

func NewArray(value Type, length int, variable bool) Array {
	return Array{
		Value:  value,
		Length: length,
	}
}

type Map struct {
	Key   Type
	Value Type
}

func NewMap(key, value Type) Map {
	return Map{
		Key:   key,
		Value: value,
	}
}

type Func struct {
	name    string
	Params  Tuple
	Results Tuple
}

func NewFunc(params, results Tuple) Func {
	return Func{
		Params:  params,
		Results: results,
	}
}

type Tuple struct {
	types []Type
}

func NewTuple(types ...Type) Tuple {
	return Tuple{
		types: types,
	}
}

func (v *Validator) ExpressionTuple(exprs []ast.ExpressionNode) Tuple {
	var types []Type
	for _, expression := range exprs {
		typ := v.resolveExpression(expression)
		// expression tuples force inner tuples to just be lists of types
		// ((int, string)) --> (int, string)
		// ((int), string) --> (int, string)
		// this is to facilitate assignment comparisons
		if tuple, ok := typ.(Tuple); ok {
			types = append(types, tuple.types...)
		} else {
			types = append(types, typ)
		}
	}
	return NewTuple(types...)
}

type Aliased struct {
	alias      string
	underlying Type
}

func NewAliased(alias string, underlying Type) Aliased {
	return Aliased{
		alias:      alias,
		underlying: underlying,
	}
}

type Lifecycle struct {
	Type       lexer.TokenType
	Parameters []Type
}

func NewLifecycle(typ lexer.TokenType, params []Type) Lifecycle {
	return Lifecycle{
		Type:       typ,
		Parameters: params,
	}
}

// A Class is a collection of properties
type Class struct {
	Name       string
	Lifecycles map[lexer.TokenType][]Lifecycle
	Supers     []*Class
	Properties map[string]Type
	Types      map[string]Type
	Interfaces []*Interface
}

func NewClass(name string, supers []*Class, interfaces []*Interface, types, properties map[string]Type, lifecycles map[lexer.TokenType][]Lifecycle) Class {
	return Class{
		Name:       name,
		Supers:     supers,
		Properties: properties,
		Interfaces: interfaces,
		Types:      types,
		Lifecycles: lifecycles,
	}
}

type Enum struct {
	Name   string
	Supers []*Enum
	Items  []string
}

func NewEnum(name string, supers []*Enum, items []string) Enum {
	return Enum{
		Name:   name,
		Supers: supers,
		Items:  items,
	}
}

type Interface struct {
	Name   string
	Supers []*Interface
	Funcs  map[string]Func
}

func NewInterface(name string, supers []*Interface, funcs map[string]Func) Interface {
	return Interface{
		Name:   name,
		Supers: supers,
		Funcs:  funcs,
	}
}

// Contract ...
type Contract struct {
	Name       string
	Supers     []*Contract
	Interfaces []*Interface
	Lifecycles map[lexer.TokenType][]Lifecycle
	Types      map[string]Type
	Properties map[string]Type
}

func NewContract(name string, supers []*Contract, interfaces []*Interface, types, properties map[string]Type, lifecycles map[lexer.TokenType][]Lifecycle) Contract {
	return Contract{
		Name:       name,
		Supers:     supers,
		Interfaces: interfaces,
		Properties: properties,
		Types:      types,
		Lifecycles: lifecycles,
	}
}

// Event ...
type Event struct {
	Name       string
	Parameters Tuple
}

// NewEvent ...
func NewEvent(name string, params Tuple) Event {
	return Event{
		Name:       name,
		Parameters: params,
	}
}