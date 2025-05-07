package check

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/grafana/alloy/syntax/ast"
	"github.com/grafana/alloy/syntax/diag"
	"github.com/grafana/alloy/syntax/internal/syntaxtags"
)

type checkContext struct {
	tags map[string]syntaxtags.Field
}

func Block(b *ast.BlockStmt, args any) diag.Diagnostics {
	rv := reflect.ValueOf(args)

	c := checkContext{tags: getTags(rv.Type().Elem())}

	fmt.Printf("%+v\n", c)

	return nil
}

func block(b *ast.BlockStmt) diag.Diagnostics {
	return nil
}

func getTags(t reflect.Type) map[string]syntaxtags.Field {
	tags := syntaxtags.Get(t)

	m := make(map[string]syntaxtags.Field, len(tags))
	for _, tag := range tags {
		m[strings.Join(tag.Name, ".")] = tag
	}

	return m
}
