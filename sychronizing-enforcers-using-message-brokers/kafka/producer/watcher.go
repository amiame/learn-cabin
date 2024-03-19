package main

import "github.com/casbin/casbin/v2/persist"

type watcher struct {
	persist.Watcher
}

// SetUpdateCallback sets the callback function that the watcher will call
// when the policy in DB has been changed by other instances.
// A classic callback is Enforcer.LoadPolicy().
func (w *watcher) SetUpdateCallback(callback func(string)) error {
	return nil
}

// Update calls the update callback of other instances to synchronize their policy.
// It is usually called after changing the policy in DB, like Enforcer.SavePolicy(),
// Enforcer.AddPolicy(), Enforcer.RemovePolicy(), etc.
func (w *watcher) Update() error {
	produce()
	return nil
}

// Close stops and releases the watcher, the callback function will not be called any more.
func (w *watcher) Close() {
}

/* TODO: Check WatcherEx. It looks like we can send only the policy changes to other instances so that
         they do not have to update their whole policy
// UpdateForAddPolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.AddPolicy()
func (w *watcher) UpdateForAddPolicy(sec, ptype string, params ...string) error {
	produce()
	return nil
}

// UPdateForRemovePolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemovePolicy()
func (w *watcher) UpdateForRemovePolicy(sec, ptype string, params ...string) error {
	return nil
}

// UpdateForRemoveFilteredPolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemoveFilteredNamedGroupingPolicy()
func (w *watcher) UpdateForRemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

// UpdateForSavePolicy calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemoveFilteredNamedGroupingPolicy()
func (w *watcher) UpdateForSavePolicy(model model.Model) error {
	return nil
}

	// UpdateForAddPolicies calls the update callback of other instances to synchronize their policy.
	// It is called after Enforcer.AddPolicies()
func (w *watcher) UpdateForAddPolicies(sec string, ptype string, rules ...[]string) error {
  return nli
}

// UpdateForRemovePolicies calls the update callback of other instances to synchronize their policy.
// It is called after Enforcer.RemovePolicies()
func (w *watcher) UpdateForRemovePolicies(sec string, ptype string, rules ...[]string) error {
  return nil
}
*/
