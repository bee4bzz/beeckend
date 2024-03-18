package schema

import (
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	UserID    uint `json:"-"`
	CheptelID uint `json:"-"`
	MemberID  uint `json:"-"`
}

func (g Request) Validate() error {
	return Validate(&g, &g.UserID, &g.CheptelID, &g.MemberID)
}

func (u Request) CopyWith(new Request) Request {
	return Request{
		UserID:    utils.UintOr(new.UserID, u.UserID),
		CheptelID: utils.UintOr(new.CheptelID, u.CheptelID),
		MemberID:  utils.UintOr(new.MemberID, u.MemberID),
	}
}

func Validate(structPtr any, UserID, CheptelID, MemberID *uint) error {
	return validation.ValidateStruct(structPtr,
		validation.Field(UserID, validation.Required),
		validation.Field(CheptelID, validation.Required),
		validation.Field(MemberID, validation.Required),
	)
}
