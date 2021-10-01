package gql

import (
	"log"
	"reflect"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type GqlBody struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func Introspect(g GqlBody) (string, map[string]string) {
	query, _ := parser.ParseQuery(&ast.Source{Input: g.Query})
	a := query.Operations.ForName("").SelectionSet[0]
	e := reflect.ValueOf(a).Elem()
	n := e.FieldByName("Name").String()
	m := ParseInputs(e, g.Variables)
	return n, m
}

func ParseInputs(e reflect.Value, v map[string]string) map[string]string {
	len := e.FieldByName("Arguments").Len()
	arguments := map[string]string{}
	for i := 0; i < len; i++ {
		_p := e.FieldByName("Arguments").Index(i).Elem().FieldByName("Name").String()
		_v := e.FieldByName("Arguments").Index(i).Elem().FieldByName("Value").Elem().Field(0).String()
		log.Print(_p, _v)
		_, ok := v[_v]
		if ok {
			arguments[_p] = v[_v]
		} else {
			arguments[_p] = _v
		}
	}
	return arguments
}
