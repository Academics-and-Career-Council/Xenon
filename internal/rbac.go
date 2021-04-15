package internal

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mikespook/gorbac"
)

type PermissionsClinet struct {
	*gorbac.RBAC
}

// map[RoleId]PermissionIds
var jsonRoles map[string][]string

// map[RoleId]ParentIds
var jsonInher map[string][]string

var PermissionManager PermissionsClinet

// Load roles information
func (p PermissionsClinet) Init() {
	if err := LoadJson("config/roles.json", &jsonRoles); err != nil {
		log.Fatal(err)
	}
	// Load inheritance information
	if err := LoadJson("config/inher.json", &jsonInher); err != nil {
		log.Fatal(err)
	}
	rbac := gorbac.New()
	permissions := make(gorbac.Permissions)
	for rid, pids := range jsonRoles {
		role := gorbac.NewStdRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = gorbac.NewStdPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac.Add(role)
	}
	for rid, parents := range jsonInher {
		if err := rbac.SetParents(rid, parents); err != nil {
			log.Fatal(err)
		}
	}
	p.RBAC = rbac

	if p.IsGranted("admin", permissions["getCourseData"], nil) {
		log.Println("Permissions Checked!")
	}
}

func LoadJson(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
