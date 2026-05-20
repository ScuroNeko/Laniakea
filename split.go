package laniakea

// SplitMessageText splits plain text into Telegram-safe message chunks.
//
// The function preserves the original text exactly: concatenating all returned
// chunks reconstructs text byte-for-byte. It prefers splitting at newlines or
// spaces within the Telegram message limit and falls back to hard rune-based
// splits when no separator is available.
func SplitMessageText(text string) []string {
	return splitTextByLimit(text, maxMessageTextLen)
}

func splitTextByLimit(text string, limit int) []string {
	if text == "" {
		return nil
	}

	runes := []rune(text)
	chunks := make([]string, 0, len(runes)/limit+1)

	for start := 0; start < len(runes); {
		end := start + limit
		if end >= len(runes) {
			chunks = append(chunks, string(runes[start:]))
			break
		}

		splitAt := -1
		for i := end - 1; i > start; i-- {
			if runes[i] == '\n' || runes[i] == ' ' {
				splitAt = i + 1
				break
			}
		}
		if splitAt == -1 {
			splitAt = end
		}

		chunks = append(chunks, string(runes[start:splitAt]))
		start = splitAt
	}

	return chunks
}
