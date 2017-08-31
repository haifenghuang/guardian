package parser

import (
	"testing"

	"github.com/end-r/guardian/go/compiler/ast"

	"github.com/end-r/goutil"
)

func TestParseInterfaceDeclarationEmpty(t *testing.T) {
	p := createContractParser(`interface Wagable {}`)
	goutil.Assert(t, isInterfaceDeclaration(p), "should detect interface decl")
	parseInterfaceDeclaration(p)
	goutil.Assert(t, len(p.Scope.Nodes["interface"]) == 1, "wrong node count")
}

func TestParseContractDeclarationEmpty(t *testing.T) {
	p := createParser(`contract Wagable {}`)
	goutil.Assert(t, isContractDeclaration(p), "should detect contract decl")
	parseContractDeclaration(p)
	goutil.Assert(t, len(p.Scope.Nodes["interface"]) == 1, "wrong node count")
}

func TestParseClassDeclarationEmpty(t *testing.T) {
	p := createParser(`class Wagable {}`)
	goutil.Assert(t, isClassDeclaration(p), "should detect class decl")
	parseClassDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.ClassDeclaration, "wrong node type")
}

func TestParseTypeDeclaration(t *testing.T) {
	p := createParser(`type Wagable int`)
	goutil.Assert(t, isTypeDeclaration(p), "should detect type decl")
	parseTypeDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.TypeDeclaration, "wrong node type")
}

func TestParseExplicitVarDeclaration(t *testing.T) {
	p := createParser(`x, y = 5, 3`)
	goutil.Assert(t, isTypeDeclaration(p), "should detect type decl")
	parseTypeDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.TypeDeclaration, "wrong node type")
}

func TestParseEventDeclarationEmpty(t *testing.T) {
	p := createParser(`event Notification{}`)
	goutil.Assert(t, isEventDeclaration(p), "should detect event decl")
	parseEventDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.EventDeclaration, "wrong node type")
}

func TestParseEventDeclarationSingle(t *testing.T) {
	p := createParser(`event Notification{string}`)
	goutil.Assert(t, isEventDeclaration(p), "should detect event decl")
	parseEventDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.EventDeclaration, "wrong node type")
}

func TestParseEventDeclarationMultiple(t *testing.T) {
	p := createParser(`event Notification{string, string}`)
	goutil.Assert(t, isEventDeclaration(p), "should detect event decl")
	parseEventDeclaration(p)
	goutil.Assert(t, p.Scope.Type() == ast.EventDeclaration, "wrong node type")
}
