package laniakea

import "os"

const (
	VersionString = "0.2.0"
	VersionMajor  = 0
	VersionMinor  = 2
	VersionPatch  = 0
)

var GoVersion = os.Getenv("GoV")
