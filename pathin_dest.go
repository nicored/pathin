package pathin

type destTarget interface {
	Name() string
	ParentGroup() DestGroup
	Handlers() []handlerFunc
}

type target struct {
	name        string
	parentGroup DestGroup
	handlers    []handlerFunc
}

func (t target) Name() string {
	return t.name
}

func (t target) Handlers() []handlerFunc {
	return t.handlers
}

func (t target) ParentGroup() DestGroup {
	return t.parentGroup
}
