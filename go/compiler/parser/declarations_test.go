package parser

import (
	"axia/guardian/go/compiler/ast"
	"axia/guardian/go/util"
	"testing"
)

func TestParseTypeDeclaration(t *testing.T) {
	p := createParser("type Bowl int")
	util.Assert(t, isTypeDeclaration(p), "not a type declaration")
	parseTypeDeclaration(p)
}

func TestParseFuncDeclaration(t *testing.T) {
	p := createParser("main(){}")
	parseFuncDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.FuncDeclaration, "wrong scope type")
}

func TestParseClassDeclaration(t *testing.T) {
	p := createParser("class Dog {}")
	parseClassDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.ClassDeclaration, "wrong scope type")
}

func TestParseInterfaceDeclaration(t *testing.T) {
	p := createParser("interface Dog {}")
	parseInterfaceDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.InterfaceDeclaration, "wrong scope type")
}

func TestParseContractDeclaration(t *testing.T) {
	p := createParser("contract Dog {}")
	parseContractDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.ContractDeclaration, "wrong scope type")
	p = createParser("contract Dog inherits Animal {}")
	parseContractDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.ContractDeclaration, "wrong scope type")
	p = createParser("contract Dog<T inherits Cat> inherits Animal<T> {}")
	parseContractDeclaration(p)
	util.Assert(t, p.scope.Type() == ast.ContractDeclaration, "wrong scope type")
}