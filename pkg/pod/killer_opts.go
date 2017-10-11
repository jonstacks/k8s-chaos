package pod

import (
	"regexp"
	"time"
)

// KillerOpts allows you to set options for the pod.Killer
type KillerOpts struct {
	Period    time.Duration
	Namespace string
	Regex     *regexp.Regexp // if specified, a regex to filter the pods on
	MaxKill   int            // Maximum number of pods to kill at a time.
}

// DefaultKillerOpts returns the default options to the pod killer.
func DefaultKillerOpts() *KillerOpts {
	return &KillerOpts{
		Period:    time.Minute,
		Namespace: "does-not-exist",
		Regex:     regexp.MustCompile("."), // match everything by default.
		MaxKill:   1,
	}
}
