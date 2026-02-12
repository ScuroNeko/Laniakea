package laniakea

import "git.nix13.pw/scuroneko/laniakea/utils"

func Ptr[T any](v T) *T { return &v }

func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

const VersionString = utils.VersionString

var EscapeMarkdown = utils.EscapeMarkdown
