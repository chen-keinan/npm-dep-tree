package utils

import (
	"net/url"
	"strings"
)

// TrimVersionSign func trim the following singes (^ / ~) from version prefix
func TrimVersionSign(version string) string {
	if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~") {
		return version[1:]
	}
	return version
}

// EscapePackageName escape package name if include `/` sign,example : async/one:1.0.1
func EscapePackageName(name string) string {
	return url.PathEscape(name)
}
