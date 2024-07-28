package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type AlbumType string

const (
	Cheptels  AlbumType = "cheptel_albums"
	HiveNotes AlbumType = "hive_note_albums"
)

type CheptelAlbum struct {
	gorm.Model
	Album     `gorm:"embedded"`
	CheptelID uint           `gorm:"index:,unique,composite:key;not null"`
	Photos    []CheptelPhoto `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a CheptelAlbum) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Album, validation.Required),
		validation.Field(&a.CheptelID, validation.Required),
	)
}

type HiveNoteAlbum struct {
	gorm.Model
	Album      `gorm:"embedded"`
	HiveNoteID uint            `gorm:"index:,unique,composite:key;not null"`
	Photos     []HiveNotePhoto `gorm:"constraint:OnDelete:CASCADE;"`
}

func (a HiveNoteAlbum) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Album, validation.Required),
		validation.Field(&a.HiveNoteID, validation.Required),
	)
}

type Album struct {
	Name        string `gorm:"index:,unique,composite:key;not null"`
	Observation *string
}

// Validate Album structure.
func (a Album) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Observation, validation.NilOrNotEmpty),
	)
}

type CheptelPhoto struct {
	gorm.Model
	Path           string `gorm:"not null"`
	CheptelAlbumID uint   `gorm:"not null"`
}

// Validate Photo structure.
func (p CheptelPhoto) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Path, validation.Required),
		validation.Field(&p.CheptelAlbumID, validation.Required),
	)
}

type HiveNotePhoto struct {
	gorm.Model
	Path            string `gorm:"not null"`
	HiveNoteAlbumID uint   `gorm:"not null"`
}

// Validate Photo structure.
func (p HiveNotePhoto) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Path, validation.Required),
		validation.Field(&p.HiveNoteAlbumID, validation.Required),
	)
}
