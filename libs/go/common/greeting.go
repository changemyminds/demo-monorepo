// Package common holds small, stable helpers shared across the Go services.
// It is consumed live-at-head via go.work — no version, no release.
package common

import "fmt"

// Greeting builds a standard greeting used by all Go services, so a change here
// (e.g. wording) instantly affects every consumer without a version bump.
func Greeting(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s! (from libs/go/common)", name)
}
