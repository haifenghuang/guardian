package ast

import (
	"testing"

	"github.com/end-r/guardian/go/compiler/parser"

	"github.com/end-r/firevm"
)

func TestBinaryExpressionBytecodeLiterals(t *testing.T) {
	p := parser.parser.ParseString("1 + 5")
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"PUSH", // push hash(x)
		"ADD",
	})
}

func TestBinaryExpressionBytecodeReferences(t *testing.T) {
	p := parser.ParseString("a + b")
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"PUSH", // push hash(x)
		"ADD",
	})
}

func TestBinaryExpressionBytecodeStringLiteralConcat(t *testing.T) {
	p := parser.ParseString(`"my name is" + " who knows tbh"`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"PUSH", // push hash(x)
		"ADD",
	})
}

func TestBinaryExpressionBytecodeStringReferenceConcat(t *testing.T) {
	p := parser.ParseString(`
		var (
			a = "hello"
			b = "world"
		)
		a + b
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"PUSH", // push hash(x)
		"ADD",
	})
}

func TestUnaryExpressionBytecodeLiteral(t *testing.T) {
	p := parser.ParseString("!1")
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"NOT",
	})
}

func TestCallExpressionBytecodeLiteral(t *testing.T) {
	p := parser.ParseString(`
		doSomething("data")
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"NOT",
	})
}

func TestCallExpressionBytecodeUseResult(t *testing.T) {
	p := parser.ParseString(`
		s := doSomething("data")
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		"NOT",
	})
}

func TestCallExpressionBytecodeUseMultipleResults(t *testing.T) {
	p := parser.ParseString(`
		s, a, p := doSomething("data")
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		// this is where the function code would go
		"PUSH", // push p
		"SET",  // set p
		"PUSH", // push a
		"SET",  // set a
		"PUSH", // push s
		"SET",  // set s
	})
}

func TestCallExpressionBytecodeIgnoredResult(t *testing.T) {
	p := parser.ParseString(`
		s, _, p := doSomething("data")
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		// this is where the function code would go
		// stack ends: {R3, R2, R1}
		"PUSH", // push p
		"SET",  // set p
		"POP",  // ignore second return value
		"PUSH", // push s
		"SET",  // set s
	})
}

func TestCallExpressionBytecodeNestedCall(t *testing.T) {
	p := parser.ParseString(`
		err := saySomething(doSomething("data"))
		`)
	vm := firevm.NewVM()
	Traverse(vm, p.Scope)
	checkMnemonics(t, vm.Instructions, []string{
		"PUSH", // push string data
		// this is where the function code would go
		// stack ends: {R1}
		// second function code goes here
		"PUSH", // push err
		"SET",  // set err
	})
}