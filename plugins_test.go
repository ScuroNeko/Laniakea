package laniakea

import (
	"errors"
	"testing"
)

func TestValidateArgsRequiresFullMatch(t *testing.T) {
	intCmd := NewCommand[NoData](func(ctx *MsgContext, db NoData) error { return nil }, "int", NewCommandArg("n").SetValueType(CommandValueIntType).SetRequired())
	if err := intCmd.validateArgs([]string{"123"}); err != nil {
		t.Fatalf("expected valid integer argument, got %v", err)
	}
	if err := intCmd.validateArgs([]string{"123abc"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial int match, got %v", err)
	}

	boolCmd := NewCommand[NoData](func(ctx *MsgContext, db NoData) error { return nil }, "bool", NewCommandArg("flag").SetValueType(CommandValueBoolType).SetRequired())
	if err := boolCmd.validateArgs([]string{"false"}); err != nil {
		t.Fatalf("expected valid bool argument, got %v", err)
	}
	if err := boolCmd.validateArgs([]string{"falsey"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial bool match, got %v", err)
	}
}

func TestValidateArgsEnforcesRequiredArgIndex(t *testing.T) {
	cmd := NewCommand[NoData](
		func(ctx *MsgContext, db NoData) error { return nil },
		"mixed",
		NewCommandArg("optional"),
		NewCommandArg("required").SetRequired(),
	)

	if err := cmd.validateArgs([]string{"only-optional"}); !errors.Is(err, ErrCmdArgCountMismatch) {
		t.Fatalf("expected ErrCmdArgCountMismatch when required second arg is missing, got %v", err)
	}
	if err := cmd.validateArgs([]string{"optional", "required"}); err != nil {
		t.Fatalf("expected both args to validate, got %v", err)
	}
}
