package pathin

type destGroup interface {
	destTarget
	Root() *root
	AddDestGroup(handlerName, ...handlerFunc) *group
	AddDest(handlerName, ...handlerFunc)
}

type group struct {
	root        *root
	name        handlerName
	parentGroup destGroup
	handlers    []handlerFunc
}

func (g group) Name() handlerName {
	return g.name
}

func (g group) ParentGroup() destGroup {
	return g.parentGroup
}

func (g group) Handlers() []handlerFunc {
	return g.handlers
}

func (g *group) AddDestGroup(name handlerName, handlersChain ...handlerFunc) *group {
	return &group{
		name:        name,
		root:        g.root,
		parentGroup: g,
		handlers:    handlersChain,
	}
}

func (g group) Root() *root {
	return g.root
}
func (g *group) AddDest(name handlerName, handlersChain ...handlerFunc) {
	if g.root == nil {
		panic("Whoops")
	}

	g.root.typeHandlers[name] = &target{
		name:        name,
		parentGroup: g,
		handlers:    handlersChain,
	}
}
