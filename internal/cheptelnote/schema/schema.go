package schema

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	UserID    uint `json:"-"`
	CheptelID uint `json:"-"`
	NoteID    uint `json:"-"`
}

func (g Request) Validate() error {
	return Validate(&g, &g.UserID, &g.CheptelID, &g.NoteID)
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
	UserID           uint           `json:"-"`
	NoteID           uint           `json:"note_ID"`
	CheptelID        uint           `json:"cheptel_ID"`
	Name             string         `json:"name"`
	TemperatureDay   *float64       `json:"temperature_day"`
	TemperatureNight *float64       `json:"temperature_night"`
	Weather          entity.Weather `json:"weather"`
	Flora            string         `json:"flora"`
	State            *string        `json:"state"`
	Observation      *string        `json:"observation"`
}

func (u CreateRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserID, validation.Required),
		validation.Field(&u.CheptelID, validation.Required),
	)
}

func (u CreateRequest) CopyWith(new CreateRequest) CreateRequest {
	return CreateRequest{
		UserID:           utils.UintOr(new.UserID, u.UserID),
		CheptelID:        utils.UintOr(new.CheptelID, u.CheptelID),
		NoteID:           utils.UintOr(new.NoteID, u.NoteID),
		Name:             utils.StringOr(new.Name, u.Name),
		TemperatureDay:   utils.Float64Or(new.TemperatureDay, u.TemperatureDay),
		TemperatureNight: utils.Float64Or(new.TemperatureNight, u.TemperatureNight),
		Weather:          u.Weather,
		Flora:            utils.StringOr(new.Flora, u.Flora),
		State:            utils.Or(new.State, u.State),
		Observation:      utils.Or(new.Observation, u.Observation),
	}
}

type UpdateRequest struct {
	UserID              uint           `json:"-"`
	CheptelID           uint           `json:"-"`
	NoteID              uint           `json:"-"`
	NewCheptelID        uint           `json:"cheptel_ID"`
	NewName             string         `json:"name"`
	NewTemperatureDay   *float64       `json:"temperature_day"`
	NewTemperatureNight *float64       `json:"temperature_night"`
	NewWeather          entity.Weather `json:"weather"`
	NewFlora            string         `json:"flora"`
	NewState            *string        `json:"state"`
	NewObservation      *string        `json:"observation"`
}

func (u UpdateRequest) Validate() error {
	return Validate(&u, &u.UserID, &u.CheptelID, &u.NoteID)
}

func (u UpdateRequest) CopyWith(new UpdateRequest) UpdateRequest {
	return UpdateRequest{
		UserID:              utils.UintOr(new.UserID, u.UserID),
		CheptelID:           utils.UintOr(new.CheptelID, u.CheptelID),
		NoteID:              utils.UintOr(new.NoteID, u.NoteID),
		NewCheptelID:        utils.UintOr(new.NewCheptelID, u.NewCheptelID),
		NewName:             utils.StringOr(new.NewName, u.NewName),
		NewTemperatureDay:   utils.Float64Or(new.NewTemperatureDay, u.NewTemperatureDay),
		NewTemperatureNight: utils.Float64Or(new.NewTemperatureNight, u.NewTemperatureNight),
		NewWeather:          u.NewWeather,
		NewFlora:            utils.StringOr(new.NewFlora, u.NewFlora),
		NewState:            utils.Or(new.NewState, u.NewState),
		NewObservation:      utils.Or(new.NewObservation, u.NewObservation),
	}
}

func Validate(structPtr any, UserID, CheptelID, NoteID *uint) error {
	return validation.ValidateStruct(structPtr,
		validation.Field(UserID, validation.Required),
		validation.Field(CheptelID, validation.Required),
		validation.Field(NoteID, validation.Required),
	)
}
