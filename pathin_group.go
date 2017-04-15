package pathin

type destGroup interface {
	destTarget
	Root() *root
	AddDestGroup(string, ...handlerFunc) *group
	AddDest(string, ...handlerFunc)
}

type group struct {
	root        *root
	name        string
	parentGroup destGroup
	handlers    []handlerFunc
}

func (g group) Name() string {
	return g.name
}

func (g group) ParentGroup() destGroup {
	return g.parentGroup
}

func (g group) Handlers() []handlerFunc {
	return g.handlers
}

func (g *group) AddDestGroup(name string, handlersChain ...handlerFunc) *group {
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
func (g *group) AddDest(name string, handlersChain ...handlerFunc) {
	if g.root == nil {
		panic("Whoops")
	}

	g.root.typeHandlers[name] = &target{
		name:        name,
		parentGroup: g,
		handlers:    handlersChain,
	}
}
