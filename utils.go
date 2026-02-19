package laniakea

import "git.nix13.pw/scuroneko/laniakea/utils"

func Ptr[T any](v T) *T { return &v }
func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}
func EscapeMarkdown(s string) string   { return utils.EscapeMarkdown(s) }
func EscapeMarkdownV2(s string) string { return utils.EscapeMarkdownV2(s) }

const VersionString = utils.VersionString
