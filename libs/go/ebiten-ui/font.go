package ebitenui

import (
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var (
	textFaceMu sync.RWMutex
	textFace   font.Face = basicfont.Face7x13
)

func TextFace() font.Face {
	textFaceMu.RLock()
	defer textFaceMu.RUnlock()
	if textFace == nil {
		return basicfont.Face7x13
	}
	return textFace
}

func SetTextFace(face font.Face) {
	textFaceMu.Lock()
	defer textFaceMu.Unlock()
	if face == nil {
		textFace = basicfont.Face7x13
		return
	}
	textFace = face
}
