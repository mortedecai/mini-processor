package _test

import (
	"os"

	. "github.com/onsi/gomega"
)

// EnvMap provides the ability to register environment variables to create or remove before or after test execution.
// Note: This is for test functionality only.  DO NOT use in production code.
type EnvMap map[string]*string

// SetupEnv processes the EnvMap, `e`, and sets up the local environment according to its contents and returns an EnvMap to undo the changes later.
// Should a key have a `nil` value, the environment variable will be unset, otherwise the value will be set (or replaced).
//
// The returned EnvMap value can be used to restore the environment to its previous state in an AfterEach or deferred call.
func (e EnvMap) SetupEnv() EnvMap {
	restoreMap := EnvMap{}
	for k, v := range e {
		if ev, found := os.LookupEnv(k); found {
			restoreMap[k] = StringPointer(ev)
		}
		var err error
		if v != nil {
			err = os.Setenv(k, *v)
		} else {
			err = os.Unsetenv(k)
		}
		Expect(err).ToNot(HaveOccurred())
	}
	return restoreMap
}
