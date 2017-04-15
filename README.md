PATHIN
======

[![Build Status](https://travis-ci.org/nxtvibe/pathin.svg)](https://travis-ci.org/nxtvibe/pathin) [![Go Report Card](https://goreportcard.com/badge/github.com/nxtvibe/pathin?branch=master)](https://goreportcard.com/report/github.com/nxtvibe/pathin?branch=master) [![Coverage Status](https://coveralls.io/repos/github/nxtvibe/pathin/badge.svg?branch=master)](https://coveralls.io/github/nxtvibe/pathin?branch=master) [![GoDoc](https://godoc.org/github.com/nxtvibe/pathin?status.svg)](https://godoc.org/github.com/nxtvibe/pathin)

Pathin is a path generator using predefined handlers that I wrote
as part of my abstract cloud storage package to store my files to the right location.

Let me give give you an example:

In your cloud storage, you probably have client data, that you only want to
store in their client bucket/directory, but this directory can also contain
multiple sub-directories, and so on. So you could end up with paths like these:

```go
    /my-product-bucket/clients/974/templates/999/images/thumbs/20170416/super_thumb.png
    /my-product-bucket/clients/974/profile/pictures/super_pic_997.png
    /my-product-bucket/public/public_file.pdf
```

Using Pathin, you could group rules, and create your target destination:

```go
    pg := pathin.New("my-product-bucket")
    pg.AddDest("public", pathin.RawHandler) // RawHandler uses the destination or group name

    clientGroup := pg.AddDestGroup("clients", clientHandler)
    clientGroup.AddDest("template-image", templatesHandler, imageHandler)
    clientGroup.AddDest("profile-picture", profilePicHandler)
```

And for the handlers:

```go
    // Useful data you want to pass in
    type MyData struct {
        clientId string
        templateId string
        fileName string
    }

    // eg. /clients/974
    func clientHandler(t string, values interface{}) (string, error) {
        myData, _ := values.(*myData)

        return "/clients/" + myData.clientId, nil
    }

    // eg. /templates/999
    func templatesHandler(t string, values interface{}) (string, error) {
        myData, _ := values.(*myData)

        return "/templates/" + myData.templateId, nil
    }

    func imageHandler(t string, values interface{}) (string, error) {
        thumbPreg := regexp.MustCompile(".*_thumb\\..+")

        path := ""
        if thumbPreg.Match(values.fileName) {
            path += "/thumbs"
        } else {
            path += "/misc"
        }

        return path + "/" + fmt.Sprint(time.Now().Format("Ymd")), nil
    }

    ... you get it
```

I would not recommend to use it in prod just yet, more tests will
need to be written.

TODO:

- More tests
- More docs

