package main

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/cmd"
	_ "github.com/acornsoft-edgecraft/edgecraft-api/docs"
)

// @title EdgeCraft Swagger API
// @version 0.1.0
// --@host localhost:8100
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
