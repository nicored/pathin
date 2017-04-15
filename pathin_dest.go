package pathin

type destTarget interface {
	Name() handlerName
	ParentGroup() destGroup
	Handlers() []handlerFunc
}

type target struct {
	name        handlerName
	parentGroup *group
	handlers    []handlerFunc
}

func (t target) Name() handlerName {
	return t.name
}

func (t target) Handlers() []handlerFunc {
	return t.handlers
}

func (t target) ParentGroup() destGroup {
	return destGroup(t.parentGroup)
}
