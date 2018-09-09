package main

import "github.com/mjibson/esc/embed"

func main() {
	embed.Run(&embed.Config{
		Package:    "main",
		OutputFile: "templates_generated.go",
		Private:    true,
		Files:      []string{"templates"},
		Ignore:     "gen.go",
	})
}
