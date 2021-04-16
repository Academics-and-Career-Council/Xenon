package internal

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mikespook/gorbac"
	"github.com/spf13/viper"
)

type PermissionsClinet struct {
	rbac        *gorbac.RBAC
	permissions gorbac.Permissions
}

// map[RoleId]PermissionIds
var jsonRoles map[string][]string

// map[RoleId]ParentIds
var jsonInher map[string][]string

var PermissionManager PermissionsClinet

// Load roles information
func PermissionsInit() {
	if err := LoadJson(viper.GetString("rbac.roles_path"), &jsonRoles); err != nil {
		log.Fatal(err)
	}
	// Load inheritance information
	if err := LoadJson(viper.GetString("rbac.inher_path"), &jsonInher); err != nil {
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
	PermissionManager.rbac = rbac
	PermissionManager.permissions = permissions
	if PermissionManager.rbac.IsGranted("admin", permissions["getCourseData"], nil) {
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
