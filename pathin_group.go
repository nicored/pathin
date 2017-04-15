package pathin

type DestGroup interface {
	destTarget
	Root() *root
	AddDestGroup(string, ...handlerFunc) DestGroup
	AddDest(string, ...handlerFunc)
}

type Group struct {
	root        *root
	name        string
	parentGroup DestGroup
	handlers    []handlerFunc
}

func (g Group) Name() string {
	return g.name
}

func (g Group) ParentGroup() DestGroup {
	return g.parentGroup
}

func (g Group) Handlers() []handlerFunc {
	return g.handlers
}

func (g *Group) AddDestGroup(name string, handlersChain ...handlerFunc) DestGroup {
	return &Group{
		name:        name,
		root:        g.root,
		parentGroup: g,
		handlers:    handlersChain,
	}
}

func (g Group) Root() *root {
	return g.root
}

func (g *Group) AddDest(name string, handlersChain ...handlerFunc) {
	if g.root == nil {
		panic("Whoops")
	}

	g.root.typeHandlers[name] = &target{
		name:        name,
		parentGroup: g,
		handlers:    handlersChain,
	}
}
