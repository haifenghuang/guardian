package validator

import (
	"github.com/end-r/guardian/typing"

	"github.com/end-r/guardian/ast"
)

func (v *Validator) resolveType(node ast.Node) typing.Type {
	if node == nil {
		// ?
		return typing.Invalid()
	}
	switch node.Type() {
	case ast.PlainType:
		r := node.(*ast.PlainTypeNode)
		return v.resolvePlainType(r)
	case ast.MapType:
		m := node.(*ast.MapTypeNode)
		return v.resolveMapType(m)
	case ast.ArrayType:
		a := node.(*ast.ArrayTypeNode)
		return v.resolveArrayType(a)
	case ast.FuncType:
		f := node.(*ast.FuncTypeNode)
		return v.resolveFuncType(f)
	}
	return typing.Invalid()
}

func (v *Validator) resolvePlainType(node *ast.PlainTypeNode) typing.Type {
	return v.getNamedType(node.Names...)
}

func (v *Validator) resolveArrayType(node *ast.ArrayTypeNode) typing.Array {
	a := typing.Array{}
	a.Value = v.resolveType(node.Value)
	return a
}

func (v *Validator) resolveMapType(node *ast.MapTypeNode) typing.Map {
	m := typing.Map{}
	m.Key = v.resolveType(node.Key)
	m.Value = v.resolveType(node.Value)
	return m
}

func (v *Validator) resolveFuncType(node *ast.FuncTypeNode) typing.Func {
	f := typing.Func{}
	f.Params = v.resolveTuple(node.Parameters)
	f.Results = v.resolveTuple(node.Results)
	return f
}

func (v *Validator) resolveTuple(nodes []ast.Node) typing.Tuple {
	t := typing.Tuple{}
	t.Types = make([]typing.Type, len(nodes))
	for i, n := range nodes {
		t.Types[i] = v.resolveType(n)
	}
	return t
}

func (v *Validator) resolveExpression(e ast.ExpressionNode) typing.Type {
	resolvers := map[ast.NodeType]resolver{
		ast.Literal:          resolveLiteralExpression,
		ast.MapLiteral:       resolveMapLiteralExpression,
		ast.ArrayLiteral:     resolveArrayLiteralExpression,
		ast.FuncLiteral:      resolveFuncLiteralExpression,
		ast.IndexExpression:  resolveIndexExpression,
		ast.CallExpression:   resolveCallExpression,
		ast.SliceExpression:  resolveSliceExpression,
		ast.BinaryExpression: resolveBinaryExpression,
		ast.UnaryExpression:  resolveUnaryExpression,
		ast.Reference:        resolveReference,
		ast.Identifier:       resolveIdentifier,
		ast.CompositeLiteral: resolveCompositeLiteral,
	}
	return resolvers[e.Type()](v, e)
}

type resolver func(v *Validator, e ast.ExpressionNode) typing.Type

func resolveIdentifier(v *Validator, e ast.ExpressionNode) typing.Type {
	i := e.(*ast.IdentifierNode)
	// look up the identifier in scope
	t := v.findVariable(i.Name)
	i.Resolved = t
	return t
}

func resolveLiteralExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	l := e.(*ast.LiteralNode)
	literalResolver, ok := v.literals[l.LiteralType]
	if ok {

		t := literalResolver(v, l.Data)
		l.Resolved = t
		return l.Resolved
	} else {
		v.addError(errStringLiteralUnsupported)
		l.Resolved = typing.Invalid()
		return l.Resolved
	}

}

func resolveArrayLiteralExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	m := e.(*ast.ArrayLiteralNode)
	keyType := v.resolveType(m.Signature.Value)
	arrayType := typing.Array{
		Value:    keyType,
		Length:   m.Signature.Length,
		Variable: m.Signature.Variable,
	}
	return arrayType
}

func resolveCompositeLiteral(v *Validator, e ast.ExpressionNode) typing.Type {
	c := e.(*ast.CompositeLiteralNode)
	c.Resolved = v.getNamedType(c.TypeName)
	if c.Resolved == typing.Unknown() {
		c.Resolved = v.getDeclarationNode([]string{c.TypeName})
	}
	return c.Resolved
}

func resolveFuncLiteralExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be func literal
	f := e.(*ast.FuncLiteralNode)
	var params, results []typing.Type
	for _, p := range f.Parameters {
		typ := v.resolveType(p.DeclaredType)
		for _ = range p.Identifiers {
			params = append(params, typ)
		}
	}

	for _, r := range f.Results {
		results = append(results, v.resolveType(r))
	}
	f.Resolved = typing.Func{
		Params:  typing.NewTuple(params...),
		Results: typing.NewTuple(results...),
	}

	return f.Resolved
}

func resolveMapLiteralExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	m := e.(*ast.MapLiteralNode)
	keyType := v.resolveType(m.Signature.Key)
	valueType := v.resolveType(m.Signature.Value)
	mapType := typing.Map{Key: keyType, Value: valueType}
	m.Resolved = mapType
	return m.Resolved
}

func resolveIndexExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	i := e.(*ast.IndexExpressionNode)
	exprType := v.resolveExpression(i.Expression)
	// enforce that this must be an array/map type
	switch t := exprType.(type) {
	case typing.Array:
		i.Resolved = t.Value
		break
	case typing.Map:
		i.Resolved = t.Value
		break
	default:
		i.Resolved = typing.Invalid()
		break
	}
	return i.Resolved
}

/*
// attempts to resolve an expression component as a type name
// used in constructors e.g. Dog()
func (v *Validator) attemptToFindType(e ast.ExpressionNode) typing.Type {
	var names []string
	switch res := e.(type) {
	case ast.IdentifierNode:
		names = append(names, res.Name)
	case ast.ReferenceNode:
		var current ast.ExpressionNode
		for current = res; current != nil; current = res.Reference {
			switch a := current.(type) {
			case ast.ReferenceNode:
				if p, ok := a.Parent.(ast.IdentifierNode); ok {
					names = append(names, p.Name)
				} else {
					return typing.Unknown()
				}
				break
			case ast.IdentifierNode:
				names = append(names, a.Name)
				break
			default:
				return typing.Unknown()
			}
		}
		break
	default:
		return typing.Unknown()
	}
	return v.getNamedType(names...)
}*/
/*
func (v *Validator) resolveInContext(t typing.Type, property string) typing.Type {
	switch r := t.(type) {
	case typing.Class:
		t, ok := r.Types[property]
		if ok {
			return t
		}
		t, ok = r.Properties[property]
		if ok {
			return t
		}
		break
	case typing.Contract:
		t, ok := r.Types[property]
		if ok {
			return t
		}
		t, ok = r.Properties[property]
		if ok {
			return t
		}
		break
	case typing.Interface:
		t, ok := r.Funcs[property]
		if ok {
			return t
		}
		break
	case typing.Enum:
		for _, item := range r.Items {
			if item == property {
				return v.SmallestNumericType(typing.BitsNeeded(len(r.Items)), false)
			}
		}
		break
	}
	return typing.Unknown()
}*/

func resolveCallExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be call expression
	c := e.(*ast.CallExpressionNode)
	// return type of a call expression is always a tuple
	// tuple may be empty or single-valued
	call := v.resolveExpression(c.Call)
	var under typing.Type
	if call.Compare(typing.Unknown()) {
		// try to resolve as a type name
		switch n := c.Call.(type) {
		case *ast.IdentifierNode:
			under = v.getNamedType(n.Name)
		}

	} else {
		under = typing.ResolveUnderlying(call)
	}
	switch ctwo := under.(type) {
	case typing.Func:
		c.Resolved = ctwo.Results
		return c.Resolved
	case typing.Class:
		c.Resolved = ctwo
		return c.Resolved
	}
	v.addError(errCallExpressionNoFunc, typing.WriteType(call))
	c.Resolved = typing.Invalid()
	return c.Resolved

}

func resolveSliceExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	s := e.(*ast.SliceExpressionNode)
	exprType := v.resolveExpression(s.Expression)
	// must be an array
	switch t := exprType.(type) {
	case typing.Array:
		s.Resolved = t
		return s.Resolved
	}
	s.Resolved = typing.Invalid()
	return s.Resolved
}

func resolveBinaryExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be literal
	b := e.(*ast.BinaryExpressionNode)
	// rules for binary Expressions
	leftType := v.resolveExpression(b.Left)
	rightType := v.resolveExpression(b.Right)
	operatorFunc, ok := v.operators[b.Operator]
	if !ok {
		b.Resolved = typing.Invalid()
		return b.Resolved
	}
	t := operatorFunc(v, leftType, rightType)
	b.Resolved = t
	return b.Resolved
	/*
		switch b.Operator {
		case token.Add:
			// can be numeric or a string
			// string = type user has defined as string literal
			getStrType, ok := v.literals[token.String]
			if ok && v.resolveExpression(b.Left).Compare(getStrType(v)) {
				return getStrType(v)
			} else {
				return v.resolveNumericType()
			}
		case token.Sub, token.Div, token.Mul, token.Mod:
			// must be numeric
			return v.resolveNumericType()
		case token.Geq, token.Leq, token.Lss, token.Gtr:
			// must be numeric
			return standards[boolean]
		case token.Eql, token.Neq:
			// don't have to be numeric
			return standards[boolean]
		case token.Shl, token.Shr, token.And, token.Or, token.Xor:
			// must be numeric
			return standards[Int]
		case token.LogicalAnd, token.LogicalOr:
			// must be boolean
			return standards[boolean]
		case token.As:
			// make sure this is a type

		}

		// else it is a type which is not defined for binary operators
		return typing.Invalid()
	*/
}

func resolveUnaryExpression(v *Validator, e ast.ExpressionNode) typing.Type {
	m := e.(*ast.UnaryExpressionNode)
	operandType := v.resolveExpression(m.Operand)
	m.Resolved = operandType
	return operandType
}

func (v *Validator) resolveContextualReference(context typing.Type, exp ast.ExpressionNode) typing.Type {
	// check if context is subscriptable
	if isSubscriptable(context) {
		if name, ok := getIdentifier(exp); ok {
			if _, ok := v.getPropertyType(context, name); ok {
				if exp.Type() == ast.Reference {
					a := exp.(*ast.ReferenceNode)
					context = v.resolveExpression(a.Parent)
					return v.resolveContextualReference(context, a.Reference)
				}
				return v.resolveExpression(exp)
			} else {
				v.addError(errPropertyNotFound, typing.WriteType(context), name)
			}
		} else {
			v.addError(errUnnamedReference)
		}
	} else {
		v.addError(errInvalidSubscriptable, typing.WriteType(context))
	}
	return typing.Invalid()
}

func resolveReference(v *Validator, e ast.ExpressionNode) typing.Type {
	// must be reference
	m := e.(*ast.ReferenceNode)
	context := v.resolveExpression(m.Parent)
	t := v.resolveContextualReference(context, m.Reference)
	m.Resolved = t
	return m.Resolved
}

func getIdentifier(exp ast.ExpressionNode) (string, bool) {
	switch exp.Type() {
	case ast.Identifier:
		i := exp.(*ast.IdentifierNode)
		return i.Name, true
	case ast.CallExpression:
		c := exp.(*ast.CallExpressionNode)
		return getIdentifier(c.Call)
	case ast.SliceExpression:
		s := exp.(*ast.SliceExpressionNode)
		return getIdentifier(s.Expression)
	case ast.IndexExpression:
		i := exp.(*ast.IndexExpressionNode)
		return getIdentifier(i.Expression)
	case ast.Reference:
		r := exp.(*ast.ReferenceNode)
		return getIdentifier(r.Parent)
	default:
		return "", false
	}
}

func (v *Validator) getPropertiesType(t typing.Type, names []string) (resolved typing.Type) {
	var working bool
	for _, name := range names {
		if !working {
			break
		}
		t, working = v.getPropertyType(t, name)
	}
	return t
}

func (v *Validator) getPropertyType(t typing.Type, name string) (typing.Type, bool) {
	// only classes, interfaces, contracts and enums are subscriptable
	switch c := t.(type) {
	case typing.Class:
		p, has := c.Properties[name]
		return p, has
	case typing.Contract:
		p, has := c.Properties[name]
		return p, has
	case typing.Interface:
		p, has := c.Funcs[name]
		return p, has
	case typing.Enum:
		for _, s := range c.Items {
			if s == name {
				return v.SmallestNumericType(len(c.Items), false), true
			}
		}
		// TODO: fix this
		return v.SmallestNumericType(len(c.Items), false), false
	}
	return typing.Invalid(), false
}

func isSubscriptable(t typing.Type) bool {
	// only classes, interfaces and enums are subscriptable
	switch t.(type) {
	case typing.Class, typing.Interface, typing.Enum, typing.Contract:
		return true
	}
	return false
}
