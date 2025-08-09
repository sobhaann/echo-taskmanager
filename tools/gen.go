package main

import (
	"github.com/sobhaann/echo-taskmanager/models"
	"gorm.io/gen"
)

//go:generate go run gen.go

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: ".././dao", // Output directory for generated code
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(&models.Task{})
	g.ApplyBasic()

	g.Execute()
}
