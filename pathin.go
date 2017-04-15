package pathin

import (
	"errors"
	"fmt"
	"path/filepath"
)

type (
	rootName     string
	handlerName  string
	handlerFunc  func(handlerName string, values interface{}) (string, error)
	typeHandlers map[handlerName]destTarget
)

type Root struct {
	name         rootName
	typeHandlers typeHandlers
	mainGroup    *Group
}

func New(name rootName) *Root {
	newRoot := &Root{
		name:         name,
		typeHandlers: typeHandlers{},
	}

	newRoot.mainGroup = &Group{
		root:     newRoot,
		name:     handlerName(name),
		handlers: []handlerFunc{},
	}

	return newRoot
}

func (r Root) GetPath(targetName string, values interface{}) (string, error) {
	if handlers, ok := r.typeHandlers[handlerName(targetName)]; ok {
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

func (r Root) Name() string {
	return string(r.name)
}

func (r *Root) AddDestGroup(name handlerName, destHandlerChain ...handlerFunc) *Group {
	return &Group{
		name:        name,
		parentGroup: r.mainGroup,
		root:        r,
		handlers:    destHandlerChain,
	}
}

func (r *Root) AddDest(name handlerName, destHandlerChain ...handlerFunc) {
	r.typeHandlers[name] = &target{
		name:        name,
		parentGroup: r.mainGroup,
		handlers:    destHandlerChain,
	}
}
