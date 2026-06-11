package laniakea

import (
	"crypto/rand"
	"encoding/base64"

	"git.scuroneko.dev/scuroneko/laniakea/utils"
)

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T { return &v }

// Val returns dereferenced pointer value or def when p is nil.
func Val[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

const (
	// VersionString re-exports the module version string.
	VersionString = utils.VersionString
	// VersionMajor re-exports the module major version.
	VersionMajor = utils.VersionMajor
	// VersionMinor re-exports the module minor version.
	VersionMinor = utils.VersionMinor
	// VersionPatch re-exports the module patch version.
	VersionPatch = utils.VersionPatch
	// VersionBeta re-exports the module prerelease counter.
	VersionBeta = utils.VersionBeta
)

func generateToken(b int) (string, error) {
	bytes := make([]byte, b)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
