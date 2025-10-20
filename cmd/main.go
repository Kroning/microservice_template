package main

import (
	"fmt"

	"{{index .App "git"}}/internal/app/container"
)

func main() {
	var appContainer = &container.Container{}
	{{if index .Modules "vault"}}
	fmt.Println("vault")
	{{end}}
	appContainer.RunApp()
}
