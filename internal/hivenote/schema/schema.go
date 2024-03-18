package schema

import (
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	UserID     uint `json:"-"`
	CheptelID  uint `json:"-"`
	HiveID     uint `json:"-"`
	HiveNoteID uint `json:"-"`
}

func (g Request) Validate() error {
	return Validate(&g, &g.UserID, &g.CheptelID, &g.HiveID, &g.HiveNoteID)
}

type QueryRequest struct {
	UserID uint `json:"-"`
}

func (q QueryRequest) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.UserID, validation.Required),
	)
}

type CreateRequest struct {
	UserID     uint   `json:"-"`
	HiveID     uint   `json:"hive_ID"`
	HiveNoteID uint   `json:"hive_note_ID"`
	CheptelID  uint   `json:"cheptel_ID"`
	Name       string `json:"name"`
}

func (u CreateRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserID, validation.Required),
		validation.Field(&u.CheptelID, validation.Required),
		validation.Field(&u.HiveID, validation.Required),
	)
}

func (u CreateRequest) CopyWith(new CreateRequest) CreateRequest {
	return CreateRequest{
		UserID:     utils.UintOr(new.UserID, u.UserID),
		CheptelID:  utils.UintOr(new.CheptelID, u.CheptelID),
		HiveID:     utils.UintOr(new.HiveID, u.HiveID),
		HiveNoteID: utils.UintOr(new.HiveNoteID, u.HiveNoteID),
		Name:       utils.StringOr(new.Name, u.Name),
	}
}

type UpdateRequest struct {
	UserID     uint   `json:"-"`
	CheptelID  uint   `json:"-"`
	HiveID     uint   `json:"-"`
	HiveNoteID uint   `json:"-"`
	NewHiveID  uint   `json:"hive_ID"`
	NewName    string `json:"name"`
}

func (u UpdateRequest) Validate() error {
	return Validate(&u, &u.UserID, &u.CheptelID, &u.HiveID, &u.HiveNoteID)
}

func (u UpdateRequest) CopyWith(new UpdateRequest) UpdateRequest {
	return UpdateRequest{
		UserID:     utils.UintOr(new.UserID, u.UserID),
		CheptelID:  utils.UintOr(new.CheptelID, u.CheptelID),
		HiveNoteID: utils.UintOr(new.HiveNoteID, u.HiveNoteID),
		HiveID:     utils.UintOr(new.HiveID, u.HiveID),
		NewName:    utils.StringOr(new.NewName, u.NewName),
	}
}

func Validate(structPtr any, UserID, CheptelID, HiveID, HiveNoteID *uint) error {
	return validation.ValidateStruct(structPtr,
		validation.Field(UserID, validation.Required),
		validation.Field(CheptelID, validation.Required),
		validation.Field(HiveID, validation.Required),
		validation.Field(HiveNoteID, validation.Required),
	)
}
