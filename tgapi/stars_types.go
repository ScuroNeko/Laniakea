package tgapi

// StarTransaction describes a Telegram Star transaction.
// See https://core.telegram.org/bots/api#startransaction
type StarTransaction struct {
	ID             string         `json:"id"`
	Amount         int            `json:"amount"`
	NanostarAmount int            `json:"nanostar_amount,omitempty"`
	Date           int            `json:"date"`
	Source         map[string]any `json:"source,omitempty"`
	Receiver       map[string]any `json:"receiver,omitempty"`
}

// StarTransactions contains a list of Telegram Star transactions.
// See https://core.telegram.org/bots/api#startransactions
type StarTransactions struct {
	Transactions []StarTransaction `json:"transactions"`
}
