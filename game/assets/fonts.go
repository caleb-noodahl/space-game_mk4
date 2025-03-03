package assets

import (
	"embed"
	"io/fs"
)

//go:embed fonts/*
var fonts embed.FS

func LoadDigital() []byte {
	content, err := fs.ReadFile(fonts, "fonts/d7.ttf")
	if err != nil {
		panic(err)
	}
	return content
}

func LoadText() []byte {
	content, err := fs.ReadFile(fonts, "fonts/ptm.ttf")
	if err != nil {
		panic(err)
	}
	return content
}
