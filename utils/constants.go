package utils

import _ "embed"

var COMPILATIONTIME string
var BUILDCOUNT string
var COMMIT string
var BRAND string

//go:embed embeded/BUILDTOKEN.key
var BUILDTOKEN string

const (
	IDENTIFIER = "Speed"
	VERSION    = "4.0.2"
)
