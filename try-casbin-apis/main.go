package main

import (
	"fmt"
	"log"
	"strings"

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

	getUserRoles(e, "app1", "doc", "doc1", "amir")
	getRoleUsers(e, "app1", "doc", "doc1", "owner")
	getObjectRoles(e, "app1", "doc", "doc1")

}

func getUserRoles(e *casbin.Enforcer, app string, object string, instance string, user string) {
	// How to get roles of a user
	// Given object ID, instance ID, and user ID
	rm := e.GetRoleManager()
	roles, err := rm.GetRoles(user, app)
	if len(roles) == 0 {
		fmt.Println("no roles")
	}
	exitIfError(err)

	roleNames := make([]string, 0)
	for _, r := range roles {
		s1 := strings.Split(r, ":")
		if s1[0] != object {
			continue
		}

		s2 := strings.Split(s1[1], "-")
		if s2[0] != instance {
			continue
		}

		roleNames = append(roleNames, s2[1])
	}

	if len(roleNames) == 0 {
		return
	}

	fmt.Printf("Roles of '%s' in '%s:%s' in '%s' are:\n", user, object, instance, app)
	for _, r := range roleNames {
		fmt.Printf("  - %s\n", r)
	}
	fmt.Println()
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
