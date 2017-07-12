package ast

type TypeDeclarationNode struct {
	Identifier string
}

func (n TypeDeclarationNode) Type() NodeType { return TypeDeclaration }

func (n TypeDeclarationNode) Validate(t NodeType) bool {
	return true
}

func (n TypeDeclarationNode) Declare(key string, node Node) {

}

type FuncDeclarationNode struct {
	Identifier string
	Parameters []Node
	Results    []Node
	IsAbstract bool
}

func (n FuncDeclarationNode) Type() NodeType { return FuncDeclaration }

func (n FuncDeclarationNode) Validate(t NodeType) bool {
	return true
}

func (n FuncDeclarationNode) Declare(key string, node Node) {

}

type ClassDeclarationNode struct {
	Identifier   string
	IsAbstract   bool
	declarations map[string][]Node
}

func (n ClassDeclarationNode) Type() NodeType { return ClassDeclaration }

func (n ClassDeclarationNode) Validate(t NodeType) bool {
	switch t {
	case ClassDeclaration, InterfaceDeclaration, ContractDeclaration:
		return true
	}
	return false
}

func (n ClassDeclarationNode) Declare(key string, node Node) {

}

type InterfaceDeclarationNode struct {
	Identifier   string
	IsAbstract   bool
	declarations map[string][]Node
}

func (n InterfaceDeclarationNode) Type() NodeType { return InterfaceDeclaration }

func (n InterfaceDeclarationNode) Validate(t NodeType) bool {
	return true
}

func (n InterfaceDeclarationNode) Declare(key string, node Node) {

}

type ContractDeclarationNode struct {
	Identifier   string
	IsAbstract   bool
	declarations map[string][]Node
}

func (n ContractDeclarationNode) Type() NodeType { return ContractDeclaration }

func (n ContractDeclarationNode) Validate(t NodeType) bool {
	return true
}
func (n ContractDeclarationNode) Declare(key string, node Node) {

}

type MapTypeNode struct {
	Key   Node
	Value Node
}

func (n MapTypeNode) Type() NodeType { return MapType }

func (n MapTypeNode) Validate(t NodeType) bool {
	return true
}

func (n MapTypeNode) Declare(key string, node Node) {

}

type ArrayTypeNode struct {
	Value Node
}

func (n ArrayTypeNode) Type() NodeType { return ArrayType }

func (n ArrayTypeNode) Validate(t NodeType) bool {
	return true
}

func (n ArrayTypeNode) Declare(key string, node Node) {

}