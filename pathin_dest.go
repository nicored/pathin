package pathin

type destTarget interface {
	Name() string
	ParentGroup() DestGroup
	Handlers() []HandlerFunc
}

type target struct {
	name        string
	parentGroup DestGroup
	handlers    []HandlerFunc
}

func (t target) Name() string {
	return t.name
}

func (t target) Handlers() []HandlerFunc {
	return t.handlers
}

func (t target) ParentGroup() DestGroup {
	return t.parentGroup
}
