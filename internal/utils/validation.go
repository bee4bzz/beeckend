package utils

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var (
	// ErrNilUUID is the error that returns in case of nil uuid.
	ErrNilUUID = validation.NewError("validation_nil_uuid", "The nil uuid is not valid {{.nil}}").SetParams(map[string]interface{}{"nil": uuid.Nil})
	rxUUID     = regexp.MustCompile(uuid.Nil.String())
	NotNilUUID = validation.NewStringRuleWithError(func(str string) bool {
		return !rxUUID.MatchString(str)
	}, ErrNilUUID)
)
