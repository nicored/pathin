package pathin

import (
	"errors"
	"fmt"
	"path/filepath"
)

type (
	handlerFunc  func(handlerName string, values interface{}) (string, error)
	typeHandlers map[string]destTarget
)

// Root holds the very top group, from which you can build your rules.
// For example
//     r := pathin.New("my-bucket-name")
//     r.AddDest("public", pathin.RawHandler)
//     clientGroup := r.AddDestGroup("client-bucket", clientHandler)
//     clientGroup.AddDest("pictures", picturesHandler)
type Root interface {
	destGroup
	GetPath(string, interface{}) (string, error)
}
type root struct {
	name         string
	typeHandlers typeHandlers
	mainGroup    *group
}

func New(name string) Root {
	newRoot := &root{
		name:         name,
		typeHandlers: typeHandlers{},
	}

	newRoot.mainGroup = &group{
		root:     newRoot,
		name:     name,
		handlers: []handlerFunc{},
	}

	return newRoot
}

func (r root) GetPath(targetName string, values interface{}) (string, error) {
	if handlers, ok := r.typeHandlers[targetName]; ok {
		pathChan, eChan := traverseHandlers(handlers, values)
		defer close(pathChan)
		defer close(eChan)

		select {
		case path := <-pathChan:
			return path, nil
		case err := <-eChan:
			return "", err
		}
	}

	return "", errors.New(fmt.Sprintf("target %s not found", targetName))
}

func traverseHandlers(dest destTarget, values interface{}) (chan string, chan error) {
	p1 := make(chan string)
	e1 := make(chan error)

	go func() {
		p2 := make(chan string)
		e2 := make(chan error)
		defer close(p2)
		defer close(e2)

		path, err := runHandlers(dest, values)
		if err != nil {
			e1 <- err
			return
		}

		if dest.ParentGroup() != nil {
			p2, e2 = traverseHandlers(dest.ParentGroup(), values)
		} else {
			go func() { p2 <- path }()
		}

		select {
		case data := <-p2:
			path = filepath.Join(string(data), path)
			p1 <- path
		case err := <-e2:
			e1 <- err
		}
	}()

	return p1, e1
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

func (r root) Name() string {
	return string(r.name)
}

func (r *root) AddDestGroup(name string, destHandlerChain ...handlerFunc) *group {
	return &group{
		name:        name,
		parentGroup: r.mainGroup,
		root:        r,
		handlers:    destHandlerChain,
	}
}

func (r *root) AddDest(name string, destHandlerChain ...handlerFunc) {
	r.typeHandlers[name] = &target{
		name:        name,
		parentGroup: r.mainGroup,
		handlers:    destHandlerChain,
	}
}

func (r *root) Handlers() []handlerFunc {
	return r.mainGroup.handlers
}

func (r *root) ParentGroup() destGroup {
	return nil
}

func (r *root) Root() *root {
	return r
}
