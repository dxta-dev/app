package util

import (
	"io"
	"log"
	"io/fs"

	"github.com/golang/freetype/truetype"
	"github.com/dxta-dev/app"
)

func getFont(path string) *truetype.Font {
	publicFS, err := fs.Sub(static.Public, "public")

	if err != nil {
		log.Println(err)
		return nil
	}

	file, err := publicFS.Open(path)

	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	fontBytes, err := io.ReadAll(file)
    if err != nil {
        panic(err)
    }

    font, err := truetype.Parse(fontBytes)
    if err != nil {
        log.Println(err)
        return nil
    }
    return font
}

func GetRegularFont() *truetype.Font {
	return getFont("fonts/GeistVariableVF.ttf")
}

func GetMonospaceFont() *truetype.Font {
	return getFont("fonts/GeistMonoVariableVF.ttf")
}
