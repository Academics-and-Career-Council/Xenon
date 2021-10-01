package gql

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"
)

type ACLMetaData struct {
	Relation  string
	Namespace string
	Object    *template.Template
}

type ketoACL map[string]ACLMetaData

var ACL ketoACL

func InitializeACL() {
	ACLDirectory := os.DirFS(viper.GetString("acl.directory"))
	namespaces := viper.GetStringSlice("acl.namespaces")
	ACL = ketoACL{}
	for _, v := range namespaces {
		aclSlice := GenerateACL(v, ACLDirectory)
		for _, acl := range aclSlice.Graphql {
			t, err := template.New(acl.Query).Parse(acl.Object)
			if err != nil {
				panic(err)
			}
			ACL[acl.Query] = ACLMetaData{Relation: acl.Relation, Namespace: v, Object: t}
		}
		for _, acl := range aclSlice.Rest {
			t, err := template.New(acl.Path).Parse(acl.Object)
			if err != nil {
				panic(err)
			}
			ACL[acl.Path] = ACLMetaData{Relation: acl.Relation, Namespace: v, Object: t}
		}
	}
	log.Println(ACL)
}

type ACLNamespace struct {
	Graphql []graphql `yaml:"graphql"`
	Rest    []rest    `yaml:"rest"`
}

type graphql struct {
	Query    string `yaml:"query"`
	Object   string `yaml:"object"`
	Relation string `yaml:"relation"`
}
type rest struct {
	Path     string `yaml:"path"`
	Object   string `yaml:"object"`
	Relation string `yaml:"relation"`
}

func GenerateACL(namespace string, ACLDirectory fs.FS) ACLNamespace {
	f, err := ACLDirectory.Open(fmt.Sprintf("%s.yml", namespace))
	if err != nil {
		panic(err)
	}
	yamlFile, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	config := new(ACLNamespace)
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return *config
}
