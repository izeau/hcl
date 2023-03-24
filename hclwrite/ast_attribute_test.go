package hclwrite

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func TestAttributeSetName(t *testing.T) {
	tests := []struct {
		src     string
		oldName string
		newName string
		want    Tokens
	}{
		{
			"foo = true",
			"foo",
			"bar",
			Tokens{
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`bar`),
					SpacesBefore: 0,
				},
				{
					Type:         hclsyntax.TokenEqual,
					Bytes:        []byte{'='},
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`true`),
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenEOF,
					Bytes:        []byte{},
					SpacesBefore: 0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s in %s", test.oldName, test.newName, test.src), func(t *testing.T) {
			f, diags := ParseConfig([]byte(test.src), "", hcl.Pos{Line: 1, Column: 1})
			if len(diags) != 0 {
				for _, diag := range diags {
					t.Logf("- %s", diag.Error())
				}
				t.Fatalf("unexpected diagnostics")
			}

			b := f.Body().GetAttribute(test.oldName)
			b.SetName(test.newName)
			got := f.BuildTokens(nil)
			format(got)
			if !reflect.DeepEqual(got, test.want) {
				diff := cmp.Diff(test.want, got)
				t.Errorf("wrong result\ngot:  %s\nwant: %s\ndiff:\n%s", spew.Sdump(got), spew.Sdump(test.want), diff)
			}
		})
	}
}

func TestAttributeSetExprRaw(t *testing.T) {
	tests := []struct {
		src    string
		name   string
		tokens Tokens
		want   Tokens
	}{
		{
			"foo = 123",
			"foo",
			Tokens{
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`true`),
					SpacesBefore: 0,
				},
			},
			Tokens{
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`foo`),
					SpacesBefore: 0,
				},
				{
					Type:         hclsyntax.TokenEqual,
					Bytes:        []byte{'='},
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`true`),
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenEOF,
					Bytes:        []byte{},
					SpacesBefore: 0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s = %s in %s", test.name, test.tokens.Bytes(), test.src), func(t *testing.T) {
			f, diags := ParseConfig([]byte(test.src), "", hcl.Pos{Line: 1, Column: 1})
			if len(diags) != 0 {
				for _, diag := range diags {
					t.Logf("- %s", diag.Error())
				}
				t.Fatalf("unexpected diagnostics")
			}

			b := f.Body().GetAttribute(test.name)
			b.SetExprRaw(test.tokens)
			got := f.BuildTokens(nil)
			format(got)
			if !reflect.DeepEqual(got, test.want) {
				diff := cmp.Diff(test.want, got)
				t.Errorf("wrong result\ngot:  %s\nwant: %s\ndiff:\n%s", spew.Sdump(got), spew.Sdump(test.want), diff)
			}
		})
	}
}

func TestAttributeSetExprValue(t *testing.T) {
	tests := []struct {
		src  string
		name string
		val  cty.Value
		want Tokens
	}{
		{
			"foo = true",
			"foo",
			cty.StringVal("enabled"),
			Tokens{
				{
					Type:         hclsyntax.TokenIdent,
					Bytes:        []byte(`foo`),
					SpacesBefore: 0,
				},
				{
					Type:         hclsyntax.TokenEqual,
					Bytes:        []byte{'='},
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenOQuote,
					Bytes:        []byte{'"'},
					SpacesBefore: 1,
				},
				{
					Type:         hclsyntax.TokenQuotedLit,
					Bytes:        []byte(`enabled`),
					SpacesBefore: 0,
				},
				{
					Type:         hclsyntax.TokenCQuote,
					Bytes:        []byte{'"'},
					SpacesBefore: 0,
				},
				{
					Type:         hclsyntax.TokenEOF,
					Bytes:        []byte{},
					SpacesBefore: 0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s = %s in %s", test.name, test.val, test.src), func(t *testing.T) {
			f, diags := ParseConfig([]byte(test.src), "", hcl.Pos{Line: 1, Column: 1})
			if len(diags) != 0 {
				for _, diag := range diags {
					t.Logf("- %s", diag.Error())
				}
				t.Fatalf("unexpected diagnostics")
			}

			b := f.Body().GetAttribute(test.name)
			b.SetExprValue(test.val)
			got := f.BuildTokens(nil)
			format(got)
			if !reflect.DeepEqual(got, test.want) {
				diff := cmp.Diff(test.want, got)
				t.Errorf("wrong result\ngot:  %s\nwant: %s\ndiff:\n%s", spew.Sdump(got), spew.Sdump(test.want), diff)
			}
		})
	}
}
