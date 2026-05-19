package tgapi

// Game represents a game.
// Since: Bot API 2.2
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation      `json:"animation,omitempty"`
}

// CallbackGame is a placeholder for the future use of callback games.
// Since: Bot API 2.2
type CallbackGame struct{}

// GameHighScore represents one row in a game high score table.
// Since: Bot API 2.2
// See https://core.telegram.org/bots/api#gamehighscore
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}
