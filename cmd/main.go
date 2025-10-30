package main

import (
	"{{index .App "git"}}/internal/app/container"
)

func main() {
	var appContainer = &container.Container{}

	appContainer.RunApp()
}
