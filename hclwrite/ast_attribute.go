// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hclwrite

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Attribute struct {
	inTree

	leadComments *node
	name         *node
	expr         *node
	lineComments *node
}

func newAttribute() *Attribute {
	return &Attribute{
		inTree: newInTree(),
	}
}

func (a *Attribute) init(name string, expr *Expression) {
	expr.assertUnattached()

	nameTok := newIdentToken(name)
	nameObj := newIdentifier(nameTok)
	a.leadComments = a.children.Append(newComments(nil))
	a.name = a.children.Append(nameObj)
	a.children.AppendUnstructuredTokens(Tokens{
		{
			Type:  hclsyntax.TokenEqual,
			Bytes: []byte{'='},
		},
	})
	a.expr = a.children.Append(expr)
	a.expr.list = a.children
	a.lineComments = a.children.Append(newComments(nil))
	a.children.AppendUnstructuredTokens(Tokens{
		{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		},
	})
}

func (a *Attribute) SetName(name string) {
	nameTok := newIdentToken(name)
	nameObj := newIdentifier(nameTok)
	a.name.ReplaceWith(nameObj)
}

func (a *Attribute) SetExprRaw(tokens Tokens) {
	expr := NewExpressionRaw(tokens)
	a.expr = a.expr.ReplaceWith(expr)
}

func (a *Attribute) SetExprValue(val cty.Value) {
	expr := NewExpressionLiteral(val)
	a.expr = a.expr.ReplaceWith(expr)
}

func (a *Attribute) Expr() *Expression {
	return a.expr.content.(*Expression)
}
