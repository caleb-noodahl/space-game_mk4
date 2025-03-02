package assets

import (
	"embed"
	"io/fs"
)

//go:embed fonts/*
var fonts embed.FS

func LoadDigital() []byte {
	content, err := fs.ReadFile(fonts, "fonts/scp.ttf")
	if err != nil {
		panic(err)
	}
	return content
}

func LoadText() []byte {
	content, err := fs.ReadFile(fonts, "fonts/dogica.ttf")
	if err != nil {
		panic(err)
	}
	return content
}
