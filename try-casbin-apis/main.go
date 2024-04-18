package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	e, err := casbin.NewEnforcer("model.conf", "policy.csv")
	exitIfError(err)

	//res, err := e.Enforce(enforceContext, "workbench", "alice", "survey-create", "survey/56")
	e.SetFieldIndex("p", "app", 1)

	getRoleUsers(e, "app1", "doc", "doc1", "owner")
	getObjectRoles(e, "app1", "doc", "doc1")

}

func getRoleUsers(e *casbin.Enforcer, app string, object string, instance string, role string) {
	// How to get users of a role
	// Given object ID, instance ID, role ID, app ID
	roleName := fmt.Sprintf("%s:%s-%s", object, instance, role)
	rm := e.GetRoleManager()
	users, err := rm.GetUsers(roleName, app)
	exitIfError(err)
	if len(users) == 0 {
		return
	}

	fmt.Printf("Users of '%s' in '%s' are:\n", roleName, app)
	for _, u := range users {
		fmt.Printf("  - %s\n", u)
	}
	fmt.Println()
}

func getObjectRoles(e *casbin.Enforcer, app string, object string, instance string) {
	// How to get roles of an instance
	// Given object ID, instance ID, app ID
	objectInstance := fmt.Sprintf("%s:%s", object, instance)
	policies := e.GetFilteredPolicy(3, objectInstance)
	temp := make([][]string, 0)
	for _, p := range policies {
		if p[0] == app {
			temp = append(temp, p)
		}
	}
	policies = temp

	if len(policies) == 0 {
		return
	}

	fmt.Printf("Roles of '%s' in '%s' are:\n", objectInstance, app)
	for _, p := range policies {
		fmt.Printf("  - %s\n", p[1])
	}
	fmt.Println()
}
