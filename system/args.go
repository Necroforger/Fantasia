package system

import "strings"

//////////////////////////////////
// 		ARGS
/////////////////////////////////

// Args ...
type Args []string

// Get returns the argument at position 'n' or an empty string
// If nothing is found
func (a Args) Get(n int) string {
	if n >= 0 && n < len(a) {
		return a[n]
	}
	return ""
}

// After is a shortcut to After(1)
func (a Args) After() string {
	return a.AfterN(1)
}

// AfterN returns the arguments after positi
func (a Args) AfterN(n int) string {
	if n >= 0 && n < len(a) {
		return strings.Join(a[n:], " ")
	}
	return ""
}
