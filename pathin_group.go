package pathin

type DestGroup interface {
	destTarget
	Root() *Root
	AddDestGroup(string, ...HandlerFunc) DestGroup
	AddDest(string, ...HandlerFunc)
}

type Group struct {
	root        *Root
	name        string
	parentGroup DestGroup
	handlers    []HandlerFunc
}

func (g Group) Name() string {
	return g.name
}

func (g Group) ParentGroup() DestGroup {
	return g.parentGroup
}

func (g Group) Handlers() []HandlerFunc {
	return g.handlers
}

func (g *Group) AddDestGroup(name string, handlersChain ...HandlerFunc) DestGroup {
	return &Group{
		name:        name,
		root:        g.root,
		parentGroup: g,
		handlers:    handlersChain,
	}
}

func (g Group) Root() *Root {
	return g.root
}

func (g *Group) AddDest(name string, handlersChain ...HandlerFunc) {
	if g.root == nil {
		panic("Whoops")
	}

	g.root.typeHandlers[name] = &target{
		name:        name,
		parentGroup: g,
		handlers:    handlersChain,
	}
}
