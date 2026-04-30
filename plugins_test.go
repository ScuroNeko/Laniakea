package laniakea

import (
	"errors"
	"testing"
)

func TestValidateArgsRequiresFullMatch(t *testing.T) {
	intCmd := NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "int", NewCommandArg("n").SetValueType(CommandValueIntType).SetRequired())
	if err := intCmd.validateArgs([]string{"123"}); err != nil {
		t.Fatalf("expected valid integer argument, got %v", err)
	}
	if err := intCmd.validateArgs([]string{"123abc"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial int match, got %v", err)
	}

	boolCmd := NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "bool", NewCommandArg("flag").SetValueType(CommandValueBoolType).SetRequired())
	if err := boolCmd.validateArgs([]string{"false"}); err != nil {
		t.Fatalf("expected valid bool argument, got %v", err)
	}
	if err := boolCmd.validateArgs([]string{"falsey"}); !errors.Is(err, ErrCmdArgRegexpMismatch) {
		t.Fatalf("expected ErrCmdArgRegexpMismatch for partial bool match, got %v", err)
	}
}

func TestValidateArgsEnforcesRequiredArgIndex(t *testing.T) {
	cmd := NewCommand(
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

func TestCommandGroupBuildsPrefixedCommandsWithoutMutatingOriginal(t *testing.T) {
	groupMiddleware := NewMiddleware("group", func(ctx *MsgContext, db NoData) bool { return true })
	commandMiddleware := NewMiddleware("command", func(ctx *MsgContext, db NoData) bool { return true })
	cmd := NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "ban").
		SetDescription("Ban user").
		Use(commandMiddleware)

	group := NewCommandGroup[NoData]("admin").
		SetSeparator("_").
		Use(groupMiddleware).
		AddCommand(cmd)

	built := group.Build()
	if len(built) != 1 {
		t.Fatalf("expected one command, got %d", len(built))
	}

	grouped := built[0]
	if grouped.command != "admin_ban" {
		t.Fatalf("expected prefixed command name, got %q", grouped.command)
	}
	if grouped.description != "Ban user" {
		t.Fatalf("expected description to be copied, got %q", grouped.description)
	}
	if cmd.command != "ban" {
		t.Fatalf("expected original command name to stay unchanged, got %q", cmd.command)
	}
	if len(cmd.middlewares) != 1 || cmd.middlewares[0].name != "command" {
		t.Fatalf("expected original command middleware to stay unchanged, got %#v", cmd.middlewares)
	}
	if len(grouped.middlewares) != 2 {
		t.Fatalf("expected group and command middleware, got %d", len(grouped.middlewares))
	}
	if grouped.middlewares[0].name != "group" || grouped.middlewares[1].name != "command" {
		t.Fatalf("expected group middleware before command middleware, got %q then %q", grouped.middlewares[0].name, grouped.middlewares[1].name)
	}
}

func TestCommandGroupBuildIsRepeatable(t *testing.T) {
	group := NewCommandGroup[NoData]("admin").
		Use(NewMiddleware("group", func(ctx *MsgContext, db NoData) bool { return true })).
		AddCommand(NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "ban").
			Use(NewMiddleware("command", func(ctx *MsgContext, db NoData) bool { return true })))

	first := group.Build()
	second := group.Build()

	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("expected one command from each build, got %d and %d", len(first), len(second))
	}
	if len(first[0].middlewares) != 2 {
		t.Fatalf("expected first build to have two middlewares, got %d", len(first[0].middlewares))
	}
	if len(second[0].middlewares) != 2 {
		t.Fatalf("expected second build to have two middlewares, got %d", len(second[0].middlewares))
	}
	if first[0] == second[0] {
		t.Fatal("expected repeated Build calls to return distinct command copies")
	}
}

func TestPluginCommandGroupRegistersBuiltCommands(t *testing.T) {
	plugin := NewPlugin[NoData]("admin")

	plugin.CommandGroup("admin", func(group *CommandGroup[NoData]) {
		group.SetSeparator("_")
		group.AddCommand(NewCommand(func(ctx *MsgContext, db NoData) error { return nil }, "ban"))
	})

	if _, ok := plugin.commands["admin_ban"]; !ok {
		t.Fatal("expected plugin to register prefixed command")
	}
	if _, ok := plugin.commands["ban"]; ok {
		t.Fatal("expected plugin not to register unprefixed command")
	}

	plugin.CommandGroup("ignored", nil)
	plugin.AddCommandGroup(nil)
}
