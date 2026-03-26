package laniakea

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	maxMessageTextLen    = 4096
	maxMessageCaptionLen = 1024
)

var (
	// ErrEmptyMessage reports that a required message text is empty.
	ErrEmptyMessage = errors.New("empty message")
	// ErrMessageTooLong reports that a message exceeds Telegram's text limit.
	ErrMessageTooLong = errors.New("message too long")
	// ErrCaptionTooLong reports that a caption exceeds Telegram's caption limit.
	ErrCaptionTooLong = errors.New("caption too long")
	// ErrMessageSplitImpossible reports that automatic message splitting cannot preserve semantics.
	ErrMessageSplitImpossible = errors.New("message split is impossible")
	// ErrPayloadTypeMismatch reports that callback payload encoding does not match bot policy.
	ErrPayloadTypeMismatch = errors.New("payload type mismatch")
	// ErrDraftChatIDZero reports that a draft has no target chat ID.
	ErrDraftChatIDZero = errors.New("zero draft chat ID")
	// ErrMessageContextNil reports that an operation requires ctx.Msg but none is set.
	ErrMessageContextNil = errors.New("message context is nil")
	// ErrEditTargetMissing reports that an edit operation has no message target.
	ErrEditTargetMissing = errors.New("edit target is missing")
	// ErrCallbackMessageMissing reports that a callback operation has no callback message target.
	ErrCallbackMessageMissing = errors.New("callback message is missing")
	// ErrDraftProviderNil reports that draft creation was requested without a draft provider.
	ErrDraftProviderNil = errors.New("draft provider is nil")
	// ErrAPIIsNil reports that an operation requires an API client but none is set.
	ErrAPIIsNil = errors.New("api is nil")
	// ErrMessageIDZero reports that an operation requires a non-zero message ID.
	ErrMessageIDZero = errors.New("message ID is zero")
)

func validateMessageText(text string) error {
	length := utf8.RuneCountInString(text)
	switch {
	case length == 0:
		return ErrEmptyMessage
	case length > maxMessageTextLen:
		return fmt.Errorf("%w: got %d, limit %d", ErrMessageTooLong, length, maxMessageTextLen)
	default:
		return nil
	}
}

func validateCaptionText(text string) error {
	length := utf8.RuneCountInString(text)
	if length > maxMessageCaptionLen {
		return fmt.Errorf("%w: got %d, limit %d", ErrCaptionTooLong, length, maxMessageCaptionLen)
	}
	return nil
}
