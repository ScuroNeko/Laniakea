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
	// ErrMessageNil reports that a required message value is nil.
	ErrMessageNil = errors.New("message is nil")
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
var (
	// ErrBindArgsTargetNotPointer reports that BindArgs received a nil or non-pointer destination.
	ErrBindArgsTargetNotPointer = errors.New("bind args: dst must be a non-nil pointer")
	// ErrBindArgsTargetNotStruct reports that BindArgs received a pointer to a non-struct value.
	ErrBindArgsTargetNotStruct = errors.New("bind args: dst must point to a struct")
	// ErrBindArgsUnsupportedFieldType reports that BindArgs encountered an unsupported field kind.
	ErrBindArgsUnsupportedFieldType = errors.New("bind args: unsupported field type")
	// ErrBindArgsConversion reports that BindArgs could not convert a string argument into a field type.
	ErrBindArgsConversion = errors.New("bind args: conversion failed")
	// ErrCantFindSession reports that no scene session matches the current context.
	ErrCantFindSession = errors.New("can't find session for this context")
	// ErrSceneNotFound reports that the requested scene is not registered.
	ErrSceneNotFound = errors.New("scene not found")
	// ErrSceneStepNotFound reports that the requested scene step is not registered.
	ErrSceneStepNotFound = errors.New("scene step not found")
	// ErrNotInScene reports that the current context has no active scene session.
	ErrNotInScene = errors.New("not in scene")
	// ErrSceneEntryNotSet reports that a scene has no configured entry step.
	ErrSceneEntryNotSet = errors.New("scene entry step not set")
	// ErrSceneRuntimeNil reports that scene APIs were used without an attached runtime.
	ErrSceneRuntimeNil = errors.New("scene runtime is nil")
)

var (
	// ErrNilBotWebhookOpts reports that a nil BotWebhookOpts was passed.
	ErrNilBotWebhookOpts = errors.New("nil BotWebhookOpts")
	// ErrNoBotWebhookOptsURL reports that BotWebhookOpts.URL is empty.
	ErrNoBotWebhookOptsURL = errors.New("empty BotWebhookOpts.URL")
	// ErrBotWebhookOptsMaxConnectionsRange reports that BotWebhookOpts.MaxConnections is out of range.
	ErrBotWebhookOptsMaxConnectionsRange = errors.New("BotWebhookOpts.MaxConnections must be between 1 and 100")
	// ErrBotUploaderWhenCertificate reports that a certificate was set without an uploader.
	ErrBotUploaderWhenCertificate = errors.New("bot uploader nil, but certificate set")
	// ErrStatusPathSecretRequired reports that UseStatusPath requires SecretToken to be set.
	ErrStatusPathSecretRequired = errors.New("SecretToken required when UseStatusPath is enabled")
	// ErrSetWebhookFailed reports that Telegram rejected the setWebhook request.
	ErrSetWebhookFailed = errors.New("failed to set webhook")
	// ErrBotAPINil reports that an operation requires an API client but none is set.
	ErrBotAPINil = errors.New("bot api is nil")
	// ErrBotWebhookOptsEmptyPath reports that BotWebhookOpts.Path is empty.
	ErrBotWebhookOptsEmptyPath = errors.New("empty BotWebhookOpts.Path")
	// ErrBotWebhookOptsPathNoSlash reports that BotWebhookOpts.Path does not start with '/'.
	ErrBotWebhookOptsPathNoSlash = errors.New("BotWebhookOpts.Path must start with '/'")
	// ErrBotWebhookOptsPathHasQueryOrFragment reports that BotWebhookOpts.Path contains a query or fragment.
	ErrBotWebhookOptsPathHasQueryOrFragment = errors.New("BotWebhookOpts.Path must not contain query or fragment")
	// ErrBotWebhookOptsPathCollidesStatus reports that BotWebhookOpts.Path collides with the reserved /status endpoint.
	ErrBotWebhookOptsPathCollidesStatus = errors.New("BotWebhookOpts.Path must not be '/status' when status path is enabled")
	// ErrBotWebhookTLSFilesIncomplete reports that only one of the two TLS files was provided.
	ErrBotWebhookTLSFilesIncomplete = errors.New("you must specify both private and public keys")
	// ErrBotWebhookTLSFilesTooMany reports that more than two TLS files were provided.
	ErrBotWebhookTLSFilesTooMany = errors.New("too many files; you must specify only private and public keys")
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
