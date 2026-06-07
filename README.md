# Fleur router extension

[![build badge](https://github.com/theobori/fleur-form/actions/workflows/build.yml/badge.svg)](https://github.com/theobori/fleur-form/actions/workflows/build.yml)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

This GitHub project is an extension for the router component of the [fleur](https://github.com/theobori/fleur) project. It allows you to create forms for users, with text fields and a submit button. The package exposes a `NewForm` constructor to create the `Form` object. Once created, you can also use the `Apply` method with the fleur router. You can also use the package’s `Apply` function, which will construct and then apply the extension to the router.

## Examples

You can check out the [ange](https://github.com/theobori/ange) project, which uses fleur-form. Below is a very simple example.

```go
package main

import (
	"fmt"
	"log"

	fleurform "github.com/theobori/fleur-form"
	"github.com/theobori/fleur/gophermap/evaluator"
	"github.com/theobori/fleur/server"
)

func main() {
	serverOptions, err := server.NewOptions(
		7070,
		"./",
		"localhost",
		true,
	)
	if err != nil {
		log.Fatalln(err)
	}

	evaluatorOptions := evaluator.Options{
		Port:                 serverOptions.Port,
		DirectoryPath:        serverOptions.DirectoryPath,
		Domain:               serverOptions.Domain,
		EnableAutoInlineText: true,
	}

	em := evaluator.RFC1436ItemsExtensionManager()

	router := server.NewRouter()
	err = fleurform.Apply(
		router,
		"/", // Here the form will be available at /form
		[]fleurform.ParameterMetadata{
			{
				Name:     "a",
				Required: true,
			},
			{
				Name:     "b",
				Required: false,
			},
			{
				Name:     "c",
				Required: true,
			},
		},
		func(p fleurform.Parameters, s *server.Server, ctx *server.RequestContext) error {
			fmt.Println(parameters)

			return nil
		},
	)

	evaluator := evaluator.NewEvaluator(&evaluatorOptions, em)
	server := server.NewServerWithRouter(serverOptions, evaluator, router)

	err = server.Serve()
	if err != nil {
		log.Fatalln(err)
	}
}
```

## Contribute

If you'd like to contribute to the project, please follow the instructions provided in the [CONTRIBUTING.md](./CONTRIBUTING.md) file.
