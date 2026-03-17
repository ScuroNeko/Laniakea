package laniakea

import (
	"errors"
	"testing"
)

func TestValidateArgsRequiresFullMatch(t *testing.T) {
	intCmd := NewCommand[NoDB](func(ctx *MsgContext, db *NoDB) {}, "int", *NewCommandArg("n").SetValueType(CommandValueIntType).SetRequired())
	if err := intCmd.validateArgs([]string{"123"}); err != nil {
		t.Fatalf("expected valid integer argument, got %v", err)
	}
	if err := intCmd.validateArgs([]string{"123abc"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial int match, got %v", err)
	}

	boolCmd := NewCommand[NoDB](func(ctx *MsgContext, db *NoDB) {}, "bool", *NewCommandArg("flag").SetValueType(CommandValueBoolType).SetRequired())
	if err := boolCmd.validateArgs([]string{"false"}); err != nil {
		t.Fatalf("expected valid bool argument, got %v", err)
	}
	if err := boolCmd.validateArgs([]string{"falsey"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial bool match, got %v", err)
	}
}
