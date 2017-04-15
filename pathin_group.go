package pathin

type destGroup interface {
	destTarget
	Root() *Root
	AddDestGroup(handlerName, ...handlerFunc) *Group
	AddDest(handlerName, ...handlerFunc)
}

type Group struct {
	root        *Root
	name        handlerName
	parentGroup destGroup
	handlers    []handlerFunc
}

func (g Group) Name() handlerName {
	return g.name
}

func (g Group) ParentGroup() destGroup {
	return g.parentGroup
}

func (g Group) Handlers() []handlerFunc {
	return g.handlers
}

func (g *Group) AddDestGroup(name handlerName, handlersChain ...handlerFunc) *Group {
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
func (g *Group) AddDest(name handlerName, handlersChain ...handlerFunc) {
	if g.root == nil {
		panic("Whoops")
	}

	g.root.typeHandlers[name] = &target{
		name:        name,
		parentGroup: g,
		handlers:    handlersChain,
	}
}
