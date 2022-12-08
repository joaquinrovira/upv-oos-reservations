package model

import "regexp"

var HHMMTimeRegex = regexp.MustCompile("[0-9]{2}:[0-9]{2}")
var NumberRegex = regexp.MustCompile("[0-9]+")
var AlreadyRegisteredRegex = regexp.MustCompile(regexp.QuoteMeta("Ya inscrito"))
