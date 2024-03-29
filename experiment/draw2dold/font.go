// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff
package draw2d

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"io/ioutil"
	"log"
	"path"
)

var (
	fontFolder = "../resource/font/"
	fonts      = make(map[string]*truetype.Font)
)

type FontStyle byte

const (
	FontStyleNormal FontStyle = iota
	FontStyleBold
	FontStyleItalic
)

type FontFamily byte

const (
	FontFamilySans FontFamily = iota
	FontFamilySerif
	FontFamilyMono
)

type FontData struct {
	Name   string
	Family FontFamily
	Style  FontStyle
}

func GetFont(fontData FontData) *truetype.Font {
	fontFileName := fontData.Name
	switch fontData.Family {
	case FontFamilySans:
		fontFileName += "s"
	case FontFamilySerif:
		fontFileName += "r"
	case FontFamilyMono:
		fontFileName += "m"
	}
	if fontData.Style&FontStyleBold != 0 {
		fontFileName += "b"
	} else {
		fontFileName += "r"
	}

	if fontData.Style&FontStyleItalic != 0 {
		fontFileName += "i"
	}
	fontFileName += ".ttf"
	font := fonts[fontFileName]
	if font != nil {
		return font
	}
	fonts[fontFileName] = loadFont(fontFileName)
	return fonts[fontFileName]
}

func GetFontFolder() string {
	return fontFolder
}

func SetFontFolder(folder string) {
	fontFolder = folder
}

func loadFont(fontFileName string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(path.Join(fontFolder, fontFileName))
	if err != nil {
		log.Println(err)
		return nil
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}
