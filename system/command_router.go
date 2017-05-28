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
	CurrentCategory string
	Prefix          string
	Routes          []*CommandRoute
	Subrouters      []*SubCommandRouter
}

// NewCommandRouter ..,
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		Prefix:     "",
		Routes:     []*CommandRoute{},
		Subrouters: []*SubCommandRouter{},
	}
}

// On adds a command router to the list of routes.
//		matcher: The regular expression to use when searching for this route.
//		handler: The handler function for this command route.
func (c *CommandRouter) On(matcher string, handler HandlerFunc) *CommandRoute {

	// Specify that the matched text must be at the beginning of the command
	// And include the router prefix
	matcher = c.Prefix + matcher + `(\s|$)`

	route := &CommandRoute{
		Matcher:  regexp.MustCompile(matcher),
		Handler:  handler,
		Name:     matcher,
		Category: c.CurrentCategory,
	}

	c.Lock()
	c.Routes = append(c.Routes, route)
	c.Unlock()

	return route
}

// SetCategory sets the routers current category
//		name: the name of the category to add new routes to by default
func (c *CommandRouter) SetCategory(name string) {
	c.Lock()
	c.CurrentCategory = name
	c.Unlock()
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

// AddSubrouter adds a subrouter to the list of subrouters.
func (c *CommandRouter) AddSubrouter(subrouter *SubCommandRouter) *SubCommandRouter {
	c.Lock()
	c.Subrouters = append(c.Subrouters, subrouter)
	c.Unlock()

	return subrouter
}

// FindMatch returns the first match found
//		name: The name of the route to find
func (c *CommandRouter) FindMatch(name string) *CommandRoute {

	for _, route := range c.Routes {
		if route.Matcher.MatchString(name) {
			return route
		}
	}

	for _, v := range c.Subrouters {
		if loc := v.Matcher.FindStringIndex(name); loc != nil {
			if match := v.Router.FindMatch(name[loc[1]:]); match != nil {
				return match
			}

			// Return the subrouters command route if nothing is found
			return v.CommandRoute
		}
	}

	return nil
}

// FindMatches will return all commands matching the given string
//		name: The name of the route to find
func (c *CommandRouter) FindMatches(name string) []*CommandRoute {
	matches := []*CommandRoute{}

	// Search routes
	for _, route := range c.Routes {
		if route.Matcher.MatchString(name) {
			matches = append(matches, route)
		}
	}

	// Search subrouters
	for _, v := range c.Subrouters {
		if v.Matcher.MatchString(name) {
			if route := v.Router.FindMatch(name); route != nil {
				matches = append(matches, route)
			} else if v.CommandRoute != nil {
				matches = append(matches, v.CommandRoute)
			}
		}
	}

	return matches
}

// GetAllRoutes returns all routes including the routes
// of this routers subrouters.
func (c *CommandRouter) GetAllRoutes() []*CommandRoute {

	routes := []*CommandRoute{}

	var find func(router *CommandRouter)
	find = func(router *CommandRouter) {

		for _, v := range router.Routes {
			routes = append(routes, v)
		}

		for _, v := range router.Subrouters {
			if v.CommandRoute != nil {
				routes = append(routes, v.CommandRoute)
			}
		}

		for _, v := range router.Subrouters {
			find(v.Router)
		}

	}

	find(c)
	return routes
}

//////////////////////////////////
// 		SUB COMMAND ROUTER
/////////////////////////////////

// SubCommandRouter is a subrouter for commands
type SubCommandRouter struct {
	Matcher *regexp.Regexp
	Router  *CommandRouter
	Name    string

	// CommandRoute is retrieved when there are no matching routes found under the subrouter,
	// But the subrouter was found.s
	CommandRoute *CommandRoute
}

// NewSubCommandRouter returns a pointer to a new SubCommandRouter
//		matcher: The regular expression to use when matching for commands.
//				 Use the expression '.' and set the prefix to an empty string.
//				 to match everything..
func NewSubCommandRouter(matcher string) (*SubCommandRouter, error) {
	reg, err := regexp.Compile(matcher)
	if err != nil {
		return nil, err
	}

	return &SubCommandRouter{
		Matcher: reg,
		Router: &CommandRouter{
			Prefix: " ",
		},
		Name:         matcher,
		CommandRoute: nil,
	}, nil
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

//////////////////////////////////
// 		SORTING BY CATEGORY
/////////////////////////////////

// CommandRoutersByCategory implements the sort.Sortable interface
// To allow CommandRouters to be sorted in alphabetical order based on their
// Category.
type CommandRoutersByCategory []*CommandRoute

func (c CommandRoutersByCategory) Swap(a, b int) {
	c[a], c[b] = c[b], c[a]
}

// Len implements the sorter.Sortable interface
func (c CommandRoutersByCategory) Len() int {
	return len(c)
}

// Less implements the sorter.Sortable interface
func (c CommandRoutersByCategory) Less(a, b int) bool {
	return c[a].Category < c[b].Category
}

// Group splits the CommandRouters into separate slices according to group
func (c CommandRoutersByCategory) Group() [][]*CommandRoute {
	var (
		groups       = [][]*CommandRoute{}
		lastCategory string
		currentGroup = []*CommandRoute{}
	)

	for _, v := range c {

		if v.Category != lastCategory {
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
				currentGroup = []*CommandRoute{}
			}
		}

		currentGroup = append(currentGroup, v)
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}
