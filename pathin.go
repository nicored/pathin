package pathin

import (
	"fmt"
	"path/filepath"
)

type (
	HandlerFunc  func(handlerName string, values interface{}) (string, error)
	typeHandlers map[string]destTarget
)

// Root holds the very top group, from which you can build your rules.
// For example
//     r := pathin.New("my-bucket-name")
//     r.AddDest("public", pathin.RawHandler)
//     clientGroup := r.AddDestGroup("client-bucket", clientHandler)
//     clientGroup.AddDest("pictures", picturesHandler)
type DestRoot interface {
	DestGroup
	GetPath(string, interface{}) (string, error)
}
type Root struct {
	name         string
	typeHandlers typeHandlers
	mainGroup    DestGroup
}

func New(name string) *Root {
	newRoot := &Root{
		name:         name,
		typeHandlers: typeHandlers{},
	}

	newRoot.mainGroup = &Group{
		root:     newRoot,
		name:     name,
		handlers: []HandlerFunc{},
	}

	return newRoot
}

func (r Root) GetPath(targetName string, values interface{}) (string, error) {
	if handlers, ok := r.typeHandlers[targetName]; ok {
		path, err := traverseHandlers(handlers, values)
		if err != nil {
			return "", fmt.Errorf("Error getting path: %s", err)
		}

		return path, nil
	}

	return "", fmt.Errorf("target %s not found", targetName)
}

func traverseHandlers(dest destTarget, values interface{}) (string, error) {
	var path2 string
	var err2 error

	if dest.ParentGroup() != nil {
		path2, err2 = traverseHandlers(dest.ParentGroup(), values)
	} else {
		path2 = ""
	}

	if err2 != nil {
		return "", fmt.Errorf("Error traversing handlers: %s", err2)
	}

	path, err := runHandlers(dest, values)
	path = filepath.Join(path2, path)

	return path, err
}

func runHandlers(dest destTarget, values interface{}) (string, error) {
	handlers := dest.Handlers()

	path := ""
	for i := 0; i < len(handlers); i++ {

		output, err := handlers[i](string(dest.Name()), values)
		if err != nil {
			return "", err
		}

		path = filepath.Join(path, output)
	}

	return path, nil
}

func (r Root) Name() string {
	return string(r.name)
}

func (r *Root) AddDestGroup(name string, destHandlerChain ...HandlerFunc) DestGroup {
	return &Group{
		name:        name,
		parentGroup: r.mainGroup,
		root:        r,
		handlers:    destHandlerChain,
	}
}

func (r *Root) AddDest(name string, destHandlerChain ...HandlerFunc) {
	r.typeHandlers[name] = &target{
		name:        name,
		parentGroup: r.mainGroup,
		handlers:    destHandlerChain,
	}
}

func (r *Root) Handlers() []HandlerFunc {
	return r.mainGroup.Handlers()
}

func (r *Root) ParentGroup() DestGroup {
	return nil
}

func (r *Root) Root() *Root {
	return r
}
