package laniakea

import "os"

const (
	VersionString = "0.1.4"
	VersionMajor  = 0
	VersionMinor  = 1
	VersionPatch  = 4
)

var GoVersion = os.Getenv("GoV")
