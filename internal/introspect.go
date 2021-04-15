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

func IntrospectGetReviews(u UserMetaData, course string) bool {
	i := sort.SearchStrings(u.CoursesUnlocked, course)
	if i < len(u.CoursesUnlocked) && u.CoursesUnlocked[i] == course {
		return true
	}
	return false
}

func getUser(email string) (*UserMetaData, error) {
	u := &UserMetaData{}
	err := BadgerDB.Get(email, u)
	if err != nil {
		err = MongoClient.getUser(email, u)
		BadgerDB.Save("permiso:"+email, u)
		if err != nil {
			return u, err
		}
	}
	return u, err
}
