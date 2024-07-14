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
	Album `gorm:"embedded"`
}

type HiveNoteAlbum struct {
	gorm.Model
	Album `gorm:"embedded"`
}

type Album struct {
	Name        string `gorm:"index:,unique,composite:key;not null"`
	Observation *string
	OwnerID     uint    `gorm:"index:,unique,composite:key;not null"`
	Photos      []Photo `gorm:"constraint:OnDelete:CASCADE;"`
}

// Validate Album structure.
func (a Album) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Observation, validation.NilOrNotEmpty),
		validation.Field(&a.OwnerID, validation.Required),
	)
}

type Photo struct {
	gorm.Model
	Path    string `gorm:"not null"`
	AlbumID uint   `gorm:"not null"`
}

// Validate Photo structure.
func (p Photo) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Path, validation.Required),
		validation.Field(&p.AlbumID, validation.Required),
	)
}
