package sql

import "regexp"

var (
	AllowedFieldAndTableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`) // To prevent SQL injections, only alphanumeric characters and dashes are allowed
)
