package engine

import "path/filepath"

type (
	groups map[string]*Group

	Group struct {
		prefix string
		engine *Engine
		parent *Group
		HttpStatuses
	}
)

func (group *Group) pathFor(path string) string {
	joined := filepath.ToSlash(filepath.Join(group.prefix, path))
	// Append a '/' if the last component had one, but only if it's not there already
	if len(path) > 0 && path[len(path)-1] == '/' && joined[len(joined)-1] != '/' {
		return joined + "/"
	}
	return joined
}

// NewGroup creates a group with no parent and the provided prefix.
func NewGroup(prefix string, engine *Engine) *Group {
	if group, exists := engine.groups[prefix]; exists {
		return group
	} else {
		newgroup := &Group{prefix: prefix,
			engine:       engine,
			HttpStatuses: defaultHttpStatuses()}
		engine.groups[prefix] = newgroup
		return newgroup
	}
}

// New creates a group from an existing group using the component string as a
// prefix. The existing group will be the nominal parent of the new group.
func (group *Group) New(component string) *Group {
	prefix := group.pathFor(component)
	newgroup := NewGroup(prefix, group.engine)
	newgroup.parent = group
	return newgroup
}

// Handle provides a route, method, and Manage to the router, and creates
// a function using the handler when the router matches the route and method.
func (group *Group) Handle(route string, method string, handler Manage) {
	group.engine.Manage(method, group.pathFor(route), func(c *Ctx) {
		c.group = group
		handler(c)
	})
}
