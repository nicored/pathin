package pathin

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type chanResp struct {
	path chan string
	err  chan error
}

func TestNewFS(t *testing.T) {
	myFs := New("bucket-name")
	assert.Equal(t, myFs.Name(), "bucket-name")

	inBucketDest := myFs.AddDestGroup("companyBucket", groupHandler)
	inBucketDest.AddDestType("cad-files", rawHandler)

	inUserDest := inBucketDest.AddDestGroup("userBucket", userHandler)
	inUserDest.AddDestType("profile-picture", rawHandler)

	if handlers, ok := myFs.typeHandlers["profile-picture"]; ok {
		pathChan, eChan := fetchDestHandler(handlers, &bucketInfo{5, 941})
		defer close(pathChan)
		defer close(eChan)

		select {
		case path := <-pathChan:
			fmt.Println(path)
			break
		case err := <-eChan:
			fmt.Println(err)
		}
	}
}

func fetchDestHandler(dest destTarget, values interface{}) (chan string, chan error) {
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
			p2, e2 = fetchDestHandler(dest.ParentGroup(), values)
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

		output, err := handlers[i](dest.Name(), values)
		if err != nil {
			return "", err
		}

		path = filepath.Join(path, output)
	}

	return path, nil
}

type child struct {
	p *child
	n []string
}
type enum string

func TestIt(t *testing.T) {
	c1 := child{
		n: []string{"1", "2", "3"},
	}
	c2 := child{
		n: []string{"4", "5", "6"},
		p: &c1,
	}
	c3 := child{
		n: []string{"7", "8", "9"},
		p: &c2,
	}

	c, e := yup(&c3)
	defer close(c)
	defer close(e)

	select {
	case path := <-c:
		fmt.Println(path)
		break
	case err := <-e:
		fmt.Println(err)
	}
}

func yup(c *child) (chan enum, chan error) {
	en := make(chan enum)
	er := make(chan error)

	go func() {
		en2 := make(chan enum)
		defer close(en2)
		er2 := make(chan error)
		defer close(er2)

		path := strings.Join(c.n, "/")

		if c.p != nil {
			en2, er2 = yup(c.p)
		} else {
			go func() { en2 <- "" }()
		}

		select {
		case data := <-en2:
			path = filepath.Join(path, string(data))
			en <- data
			return
		case err := <-er2:
			er <- err
			return
		}
	}()

	return en, er
}
