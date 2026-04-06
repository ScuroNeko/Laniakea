package tgapi

// ParseMode represents the text formatting mode for message parsing.
type ParseMode string

const (
	// ParseMDV2 enables MarkdownV2 style parsing.
	ParseMDV2 ParseMode = "MarkdownV2"
	// ParseHTML enables HTML style parsing.
	ParseHTML ParseMode = "HTML"
	// ParseMD enables legacy Markdown style parsing.
	ParseMD ParseMode = "Markdown"
	// ParseNone disables parse_mode and leaves plain-text requests unannotated.
	ParseNone ParseMode = ""
)

// EmptyParams is a placeholder for methods that take no parameters.
type EmptyParams struct{}

// NoParams is a convenient instance of EmptyParams.
var NoParams = EmptyParams{}
