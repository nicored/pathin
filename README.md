PATHIN
======

Pathin is a path generator using predefined handlers that I wrote
as part of my abstract cloud storage package to store my files to the right location.

Let me give give you an example:

In your cloud storage, you probably have client data, that you only want to
store in their client bucket/directory, but this directory can also contain
multiple sub-directories, and so on. So you could end up with paths like these:

```
    /my-product-bucket/clients/974/templates/999/images/thumbs/20170416/super_thumb.png
    /my-product-bucket/clients/974/profile/avatar/super_avatar.png
    /my-product-bucket/public/public_file.pdf
```

Using Pathin, you could group rules, and create your target destination:

```
    pg := pathin.New("my-product-bucket")
    pg.AddDest("public", pathin.RawHandler) // RawHandler uses the destination or group name

    clientGroup := pg.AddDestGroup("clients", clientHandler)
    clientGroup.AddDest("template-image", templatesHandler, imageHandler)
    clientGroup.AddDest("profile-picture", profilePicHandler)
```

And for the handlers:

```go
    type MyData struct {
        clientId string
        templateId string
        fileName string
    }

    // eg. /clients/974
    func clientHandler(t string, fileName string, values interface{}) {
        myData, _ := values.(*myData)

        return "/clients/" + myData.clientId
    }

    // eg. /templates/999
    func templatesHandler(t string, values interface{}) {
        myData, _ := values.(*myData)

        return "/templates/" + myData.templateId
    }

    func imageHandler(t string, values interface{}) {
        thumbPreg := regexp.MustCompile(".*_thumb\\..+")

        path := ""
        if thumbPreg.Match(values.fileName) {
            path += "/thumbs"
        } else {
            path += "/misc"
        }

        return path + "/" + fmt.Sprint(time.Now().Format("Ymd"))
    }

    ... you get it
```

