package optionalvalidation

import (
	"errors"

	"4d63.com/optional"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrRequired     = validation.NewError("validation_optional_required", "The optional value is required")
	StringRequired  = By[string](validation.Required)
	BoolIsPresent   = validation.By(IsPresent[bool])
	IntIsPresent    = validation.By(IsPresent[int])
	StringIsPresent = validation.By(IsPresent[string])
)

func Length(min int, max int) OptionalRule[string] {
	return By[string](validation.Length(min, max))
}

func In[T any](values ...interface{}) OptionalRule[T] {
	return By[T](validation.In(values...))
}

func Min[T any](min T) OptionalRule[T] {
	return By[T](validation.Min(min))
}

func Max[T any](max T) OptionalRule[T] {
	return By[T](validation.Max(max))
}

func By[T any](rules ...validation.Rule) OptionalRule[T] {
	return OptionalRule[T]{Rules: rules}
}

func Of[T any](rules ...validation.Rule) OptionalRule[T] {
	return OptionalRule[T]{Rules: rules}
}

// OptionalRule is a validation rule that validates a value if it is empty even with optional.Optional values.
type OptionalRule[T any] struct {
	Rules []validation.Rule
}

// Validate checks if the given value is valid or not.
func (r OptionalRule[T]) Validate(value interface{}) error {
	s, ok := value.(optional.Optional[T])

	if !ok {
		return ErrRequired
	}

	if !s.IsPresent() {
		return nil
	}

	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	underlyingValue, ok := s.Get()

	if !ok {
		return ErrRequired
	}

	for _, rule := range r.Rules {
		if err := rule.Validate(underlyingValue); err != nil {
			return err
		}
	}
	return nil
}

// Validate checks if the given value is valid or not.
func IsPresent[T any](value interface{}) error {
	s, ok := value.(optional.Optional[T])

	if !ok {
		return errors.New("value is not an optional.Optional")
	}

	if !s.IsPresent() {
		return errors.New("cannot be blank")
	}
	return nil
}
