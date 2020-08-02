package db

import (
	"image/color"
)

const PaletteKey = "palette"

func (db *DB) SetPalette(palette []color.RGBA) error {
	return db.SetData(PaletteKey, palette)
}

func (db *DB) GetPalette() ([]color.RGBA, error) {
	var palette []color.RGBA
	err := db.GetData(PaletteKey, &palette)
	return palette, err
}
