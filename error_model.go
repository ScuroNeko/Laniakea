package laniakea

import "errors"

type classifiedError struct {
	err          error
	userVisible  bool
	internalOnly bool
}

func (e *classifiedError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *classifiedError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// AsUserError marks err as safe to show to the user through the centralized
// handler error flow.
func AsUserError(err error) error {
	if err == nil {
		return nil
	}
	return &classifiedError{err: err, userVisible: true}
}

// AsInternalError marks err as internal-only so it will be logged but not sent
// to the user through the centralized handler error flow.
func AsInternalError(err error) error {
	if err == nil {
		return nil
	}
	return &classifiedError{err: err, internalOnly: true}
}

// IsUserError reports whether err was explicitly marked as user-visible.
func IsUserError(err error) bool {
	var classified *classifiedError
	if !errors.As(err, &classified) {
		return false
	}
	return classified.userVisible
}

// IsInternalError reports whether err was explicitly marked as internal-only.
func IsInternalError(err error) bool {
	var classified *classifiedError
	if !errors.As(err, &classified) {
		return false
	}
	return classified.internalOnly
}
