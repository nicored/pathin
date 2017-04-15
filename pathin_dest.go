package pathin

type destTarget interface {
	Name() string
	ParentGroup() destGroup
	Handlers() []handlerFunc
}

type target struct {
	name        string
	parentGroup *group
	handlers    []handlerFunc
}

func (t target) Name() string {
	return t.name
}

func (t target) Handlers() []handlerFunc {
	return t.handlers
}

func (t target) ParentGroup() destGroup {
	return destGroup(t.parentGroup)
}
