package validator

import (
	"github.com/end-r/guardian/ast"
)

func (v *Validator) validateStatement(node ast.Node) {
	switch node.Type() {
	case ast.AssignmentStatement:
		v.validateAssignment(node.(ast.AssignmentStatementNode))
		break
	case ast.ForStatement:
		v.validateForStatement(node.(ast.ForStatementNode))
		break
	case ast.IfStatement:
		v.validateIfStatement(node.(ast.IfStatementNode))
		break
	case ast.ReturnStatement:
		v.validateReturnStatement(node.(ast.ReturnStatementNode))
		break
	case ast.SwitchStatement:
		v.validateSwitchStatement(node.(ast.SwitchStatementNode))
		break
		// TODO: handle call statement
	}
}

func (v *Validator) validateAssignment(node ast.AssignmentStatementNode) {
	// TODO: fundamental: define operator or not
	// valid assignments must have
	// 1. valid left hand expression (cannot be a call, literal, slice)
	// 2. types of left assignable to right
	// 3. right all valid expressions

	for _, r := range node.Right {
		v.validateExpression(r)
	}

	for _, l := range node.Left {
		switch l.Type() {
		case ast.CallExpression, ast.Literal, ast.MapLiteral,
			ast.ArrayLiteral, ast.SliceExpression, ast.FuncLiteral:
			v.addError(errInvalidExpressionLeft)
		}
	}

	leftTuple := v.ExpressionTuple(node.Left)
	rightTuple := v.ExpressionTuple(node.Right)

	if len(leftTuple.types) > len(rightTuple.types) && len(rightTuple.types) == 1 {
		right := rightTuple.types[0]
		for _, left := range leftTuple.types {
			if !assignableTo(right, left) {
				v.addError(errInvalidAssignment, WriteType(left), WriteType(right))
			}
		}

		for i, left := range node.Left {
			if leftTuple.types[i] == standards[Unknown] {
				if id, ok := left.(ast.IdentifierNode); ok {
					v.DeclareVarOfType(id.Name, rightTuple.types[0])
				}
			}
		}

	} else {
		if !rightTuple.compare(leftTuple) {
			v.addError(errInvalidAssignment, WriteType(leftTuple), WriteType(rightTuple))
		}

		// length of left tuple should always equal length of left
		// this is because tuples are not first class types
		// cannot assign to tuple expressions
		if len(node.Left) == len(leftTuple.types) {
			for i, left := range node.Left {
				if leftTuple.types[i] == standards[Unknown] {
					if id, ok := left.(ast.IdentifierNode); ok {
						//fmt.Printf("Declaring %s as %s\n", id.Name, WriteType(rightTuple.types[i]))
						v.DeclareVarOfType(id.Name, rightTuple.types[i])
					}
				}
			}
		}
	}
}

func (v *Validator) validateIfStatement(node ast.IfStatementNode) {

	if node.Init != nil {
		v.validateStatement(node.Init)
	}

	for _, cond := range node.Conditions {
		// condition must be of type bool
		v.requireType(standards[Bool], v.resolveExpression(cond.Condition))
		v.validateScope(cond.Body)
	}

	if node.Else != nil {
		v.validateScope(node.Else)
	}
}

func (v *Validator) validateSwitchStatement(node ast.SwitchStatementNode) {

	switchType := v.resolveExpression(node.Target)
	// target must be matched by all cases
	for _, node := range node.Cases.Sequence {
		if node.Type() == ast.CaseStatement {
			v.validateCaseStatement(switchType, node.(ast.CaseStatementNode))
		}
	}

}

func (v *Validator) validateCaseStatement(switchType Type, clause ast.CaseStatementNode) {
	for _, expr := range clause.Expressions {
		v.requireType(switchType, v.resolveExpression(expr))
	}
	v.validateScope(clause.Block)
}

func (v *Validator) validateReturnStatement(node ast.ReturnStatementNode) {
	// must be in the context of a function (checked during ast generation)
	// resolved tuple must match function expression

	// scope is now a func declaration

}

func (v *Validator) validateForStatement(node ast.ForStatementNode) {
	// init statement must be valid
	// if it exists
	if node.Init != nil {
		v.validateAssignment(*node.Init)
	}

	// cond statement must be a boolean
	v.requireType(standards[Bool], v.resolveExpression(node.Cond))

	// post statement must be valid
	if node.Post != nil {
		v.validateStatement(node.Post)
	}

	v.validateScope(node.Block)
}