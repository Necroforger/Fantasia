package system

import (
	"errors"
	"regexp"
	"sort"
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

	// Prefix is appended to the beginning of each command added with On
	Prefix string
	// Suffix is appended to the end of each command added with On
	Suffix string

	Routes     []*CommandRoute
	Subrouters []*SubCommandRouter
}

// NewCommandRouter ..,
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		Routes:     []*CommandRoute{},
		Subrouters: []*SubCommandRouter{},
	}
}

// On adds a command router to the list of routes.
//		matcher: The regular expression to use when searching for this route.
//		handler: The handler function for this command route.
func (c *CommandRouter) On(matcher string, handler HandlerFunc) *CommandRoute {

	// Specify that the matched text must be at the beginning and end in a whitespace character
	// Or end of line.
	reg := "^" + c.Prefix + matcher + c.Suffix + `(\s|$)`

	route := &CommandRoute{
		Matcher:  regexp.MustCompile(reg),
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

// SetDisabled sets the specified command to disabled
func (c *CommandRouter) SetDisabled(name string, disabled bool) error {
	if route, _ := c.FindMatch(name); route != nil {
		route.Disabled = disabled
	}
	return errors.New("route not found")
}

// AddSubrouter adds a subrouter to the list of subrouters.
func (c *CommandRouter) AddSubrouter(subrouter *SubCommandRouter) *SubCommandRouter {

	// Set the default category to this routers current category.
	if subrouter.Category() == "" {
		subrouter.SetCategory(c.CurrentCategory)
	}

	c.Lock()
	c.Subrouters = append(c.Subrouters, subrouter)
	c.Unlock()

	return subrouter
}

// FindMatch returns the first match found
//		name: The name of the route to find
func (c *CommandRouter) findMatch(name string, skipDisabled bool) (*CommandRoute, []int) {

	for _, route := range c.Routes {
		if skipDisabled && route.Disabled == true {
			continue
		}
		if loc := route.Matcher.FindStringIndex(name); loc != nil {
			return route, loc
		}
	}

	for _, v := range c.Subrouters {
		if loc := v.Matcher.FindStringIndex(name); loc != nil {
			if match, loc2 := v.Router.findMatch(name[loc[1]:], skipDisabled); match != nil {
				return match, []int{loc[0], loc[1] + loc2[1]}
			}

			if skipDisabled && v.CommandRoute != nil && v.CommandRoute.Disabled == true {
				continue
			}

			// Return the subrouters command route if nothing is found
			return v.CommandRoute, loc
		}
	}

	return nil, nil
}

// FindMatch returns the first match that matches the given string
//		name: The name of the route to find
func (c *CommandRouter) FindMatch(name string) (*CommandRoute, []int) {
	return c.findMatch(name, false)
}

// FindEnabledMatch returns the first non-disabled route that matches the given string
//		name: The name of the route to find
func (c *CommandRouter) FindEnabledMatch(name string) (*CommandRoute, []int) {
	return c.findMatch(name, true)
}

// TODO Return an array of match locations

// FindMatches will return all commands matching the given string
//		name: The name of the route to find
// func (c *CommandRouter) FindMatches(name string) []*CommandRoute {
// 	matches := []*CommandRoute{}

// 	// Search routes
// 	for _, route := range c.Routes {
// 		if route.Matcher.MatchString(name) {
// 			matches = append(matches, route)
// 		}
// 	}

// 	// Search subrouters
// 	for _, v := range c.Subrouters {
// 		if v.Matcher.MatchString(name) {
// 			if route, _ := v.Router.FindMatch(name); route != nil {
// 				matches = append(matches, route)
// 			} else if v.CommandRoute != nil {
// 				matches = append(matches, v.CommandRoute)
// 			}
// 		}
// 	}

// 	return matches
// }

// GetAllRoutes returns all routes including the routes
// of this routers subrouters.
func (c *CommandRouter) GetAllRoutes() []*CommandRoute {

	var find func(router *CommandRouter) []*CommandRoute
	find = func(router *CommandRouter) []*CommandRoute {
		routes := []*CommandRoute{}

		for _, v := range router.Routes {
			routes = append(routes, v)
		}

		for _, v := range router.Subrouters {
			if v.CommandRoute != nil {
				routes = append(routes, v.CommandRoute)
			}
			routes = append(routes, find(v.Router)...)
		}

		return routes
	}

	return find(c)
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
	// But the subrouter was matched.
	CommandRoute *CommandRoute
}

// NewSubCommandRouter returns a pointer to a new SubCommandRouter
//		matcher: The regular expression to use when matching for commands.
//				 to match everything..
//
//		name: 	 The name to give the subrouter.
func NewSubCommandRouter(matcher string, name string) (*SubCommandRouter, error) {
	reg, err := regexp.Compile(matcher)
	if err != nil {
		return nil, err
	}

	router := NewCommandRouter()
	router.Prefix = " "
	// Set the prefix to be space separated by default.

	return &SubCommandRouter{
		Matcher:      reg,
		Router:       router,
		Name:         name,
		CommandRoute: nil,
	}, nil
}

// SetCategory sets the current category of the routers
func (s *SubCommandRouter) SetCategory(name string) {
	s.Router.SetCategory(name)
}

// Set sets the field values of the CommandRoute
// Accepts three fields:
//		1:	Name
//		2:  Description
//		3:  Category
func (s *SubCommandRouter) Set(values ...string) {
	if s.CommandRoute == nil {
		s.CommandRoute = &CommandRoute{}
	}

	switch {

	case len(values) > 2:
		if values[2] != "" {
			s.CommandRoute.Category = values[2]
		}
		fallthrough

	case len(values) > 1:
		if values[1] != "" {
			s.CommandRoute.Desc = values[1]
		}
		fallthrough

	case len(values) > 0:
		if values[0] != "" {
			s.CommandRoute.Name = values[0]
		}
	}
}

// Category returns the category of the subrouter
func (s *SubCommandRouter) Category() string {
	if s.Router != nil {
		return s.Router.CurrentCategory
	}
	return ""
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
	Disabled bool
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

// CommandRoutesByCategory implements the sort.Sortable interface
// To allow CommandRouters to be sorted in alphabetical order based on their
// Category.
type CommandRoutesByCategory []*CommandRoute

func (c CommandRoutesByCategory) Swap(a, b int) {
	c[a], c[b] = c[b], c[a]
}

// Len implements the sorter.Sortable interface
func (c CommandRoutesByCategory) Len() int {
	return len(c)
}

// Less implements the sorter.Sortable interface
func (c CommandRoutesByCategory) Less(a, b int) bool {
	return c[a].Category < c[b].Category
}

// Group splits the CommandRouters into separate slices according to group
func (c CommandRoutesByCategory) Group() [][]*CommandRoute {
	var (
		groups       = [][]*CommandRoute{}
		lastCategory = "__undefined__"
		currentGroup = []*CommandRoute{}
	)

	sort.Sort(c)

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
