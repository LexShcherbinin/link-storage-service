package main

import (
	"link-storage-service/internal/config/app"
)

func main() {
	app := app.NewApp()
	app.Run()
}
