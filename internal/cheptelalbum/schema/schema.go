package schema

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	UserID    uint `json:"-"`
	CheptelID uint `json:"-"`
	AlbumID   uint `json:"-"`
}

func (g Request) Validate() error {
	return Validate(&g, &g.UserID, &g.CheptelID, &g.AlbumID)
}

type QueryRequest struct {
	UserID    uint             `json:"-"`
	AlbumType entity.AlbumType `json:"-"`
}

func (q QueryRequest) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.UserID, validation.Required),
		validation.Field(&q.AlbumType, validation.In(entity.Cheptels)),
	)
}

type CreateRequest struct {
	UserID      uint    `json:"-"`
	CheptelID   uint    `json:"-"`
	AlbumID     uint    `json:"album_ID"`
	Name        string  `json:"name"`
	Observation *string `json:"observation"`
}

func (u CreateRequest) Validate() error {
	err := Validate(&u, &u.UserID, &u.CheptelID, &u.AlbumID)
	if err != nil {
		return err
	}
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Observation, validation.NilOrNotEmpty),
	)
}

func (u CreateRequest) CopyWith(new CreateRequest) CreateRequest {
	return CreateRequest{
		UserID:      utils.UintOr(new.UserID, u.UserID),
		CheptelID:   utils.UintOr(new.CheptelID, u.CheptelID),
		AlbumID:     utils.UintOr(new.AlbumID, u.AlbumID),
		Name:        utils.StringOr(new.Name, u.Name),
		Observation: utils.Or(new.Observation, u.Observation),
	}
}

type UpdateRequest struct {
	UserID         uint    `json:"-"`
	CheptelID      uint    `json:"-"`
	AlbumID        uint    `json:"-"`
	NewCheptelID   uint    `json:"cheptel_ID"`
	NewName        string  `json:"name"`
	NewObservation *string `json:"observation"`
}

func (u UpdateRequest) Validate() error {
	err := Validate(&u, &u.UserID, &u.CheptelID, &u.AlbumID)
	if err != nil {
		return err
	}
	return validation.ValidateStruct(&u,
		validation.Field(&u.NewObservation, validation.NilOrNotEmpty),
	)
}

func (u UpdateRequest) CopyWith(new UpdateRequest) UpdateRequest {
	return UpdateRequest{
		UserID:         utils.UintOr(new.UserID, u.UserID),
		CheptelID:      utils.UintOr(new.CheptelID, u.CheptelID),
		AlbumID:        utils.UintOr(new.AlbumID, u.AlbumID),
		NewName:        utils.StringOr(new.NewName, u.NewName),
		NewObservation: utils.Or(new.NewObservation, u.NewObservation),
	}
}

func Validate(structPtr any, UserID, CheptelID, AlbumID *uint) error {
	return validation.ValidateStruct(structPtr,
		validation.Field(UserID, validation.Required),
		validation.Field(CheptelID, validation.Required),
		validation.Field(AlbumID, validation.Required),
	)
}
