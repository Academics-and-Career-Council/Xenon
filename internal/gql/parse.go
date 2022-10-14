package gql

import (
	"log"
	"net/url"
	"reflect"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type GqlBody struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
type RestBody struct {
	Email string  `json:"email"`
	Path  url.URL `json:"path"`
}

func Introspect(g GqlBody) (string, map[string]string) {
	query, _ := parser.ParseQuery(&ast.Source{Input: g.Query})
	a := query.Operations.ForName("").SelectionSet[0]
	e := reflect.ValueOf(a).Elem()
	n := e.FieldByName("Name").String()
	v := map[string]string{}
	for key, element := range g.Variables {
		if s, ok := element.(string); ok {
			v[key] = s
		}
	}
	m := ParseInputs(e, v)
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
