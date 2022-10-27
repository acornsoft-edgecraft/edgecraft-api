package main

import (
	_ "github.com/acornsoft-edgecraft/edgecraft-api/docs"

	"github.com/acornsoft-edgecraft/edgecraft-api/cmd"
)

// @title EdgeCraft Swagger API
// @version 0.1.0
// @BasePath /api/v1
// @schemes http https
func main() {
	cmd.Execute()
}
