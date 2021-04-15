package internal

import (
	"reflect"
	"sort"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

type gqlBody struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

func Introspect(q string) string {
	query, _ := parser.ParseQuery(&ast.Source{Input: q})
	a := query.Operations.ForName("").SelectionSet[0]
	e := reflect.ValueOf(a).Elem()
	n := e.FieldByName("Name").String()
	return n
}

func IntrospectGetReviews(email string, course string) (bool, error) {
	u, err := MongoClient.getUser(email)
	if err != nil {
		return false, err
	}
	i := sort.SearchStrings(u.CoursesUnlocked, course)
	if i < len(u.CoursesUnlocked) && u.CoursesUnlocked[i] == course {
		return true, err
	}
	return false, err
}
