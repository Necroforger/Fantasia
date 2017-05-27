package system

import (
	"regexp"
	"sync"
)

// HandlerFunc ...
type HandlerFunc func(*Context)

//////////////////////////////////
// 		COMMAND ROUTER
/////////////////////////////////

// CommandRouter ...
type CommandRouter struct {
	sync.Mutex
	Prefix string
	Routes []*CommandRoute
}

// On adds a command router to the list of routes.
//		matcher: The regular expression to use when searching for this route.
//		handler: The handler function for this command route.
func (c *CommandRouter) On(matcher string, handler HandlerFunc) *CommandRoute {

	// Specify that the matched text must be at the beginning of the command
	// And include the router prefix
	matcher = c.Prefix + matcher + `(\s|$)`

	route := &CommandRoute{
		Matcher: regexp.MustCompile(matcher),
		Handler: handler,
		Name:    matcher,
	}

	c.Lock()
	c.Routes = append(c.Routes, route)
	c.Unlock()

	return route
}

// OnReg allows you to supply a custom regular expression as the route matcher.
//		matcher: The regular expression to use when searching for this route
//		handler: The handler function for this command route.
func (c *CommandRouter) OnReg(matcher string, handler HandlerFunc) *CommandRoute {
	route := &CommandRoute{
		Matcher: regexp.MustCompile(matcher),
		Handler: handler,
		Name:    matcher,
	}

	c.Lock()
	c.Routes = append(c.Routes, route)
	c.Unlock()

	return route
}

// Off removes a CommandRoute from the list of routes and returns a pointer
// To the removed value.
//		name:	The regular expression to match against
func (c *CommandRouter) Off(name string) *CommandRoute {
	c.Lock()
	defer c.Unlock()

	for i, v := range c.Routes {
		if v.Matcher.MatchString(name) {
			c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
			return v
		}
	}

	return nil
}

// FindMatch returns the first match found
//		name: The name of the route to find
func (c *CommandRouter) FindMatch(name string) *CommandRoute {
	for _, route := range c.Routes {
		if route.Matcher.MatchString(name) {
			return route
		}
	}
	return nil
}

// FindMatches will return all commands matching the given string
//		name: The name of the route to find
func (c *CommandRouter) FindMatches(name string) []*CommandRoute {
	matches := []*CommandRoute{}

	for _, route := range c.Routes {
		if route.Matcher.MatchString(name) {
			matches = append(matches, route)
		}
	}

	return matches
}

//////////////////////////////////
// 		COMMAND ROUTE
/////////////////////////////////

// CommandRoute ...
type CommandRoute struct {
	Matcher  *regexp.Regexp
	Handler  HandlerFunc
	Name     string
	Desc     string
	Category string
}

// Set sets the field values of the CommandRoute
// Accepts three fields:
//		1:	Name
//		2:  Description
//		3:  Category
func (c *CommandRoute) Set(values ...string) {
	switch {

	case len(values) > 2:
		if values[2] != "" {
			c.Category = values[2]
		}
		fallthrough

	case len(values) > 1:
		if values[1] != "" {
			c.Desc = values[1]
		}
		fallthrough

	case len(values) > 0:
		if values[0] != "" {
			c.Name = values[0]
		}
	}

}
