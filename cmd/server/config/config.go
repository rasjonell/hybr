package config

import (
	"os"
	"strings"

	"github.com/a-h/templ"
)

var (
	urlHostPrefix string
	hostPrefix    string = os.Getenv("HYBR_CONSOLE_HOST")
)

func init() {
	if hostPrefix != "" {
		if !strings.HasPrefix(hostPrefix, "/") {
			hostPrefix = "/" + hostPrefix
		}
		hostPrefix = strings.TrimSuffix(hostPrefix, "/")
	} else {
		hostPrefix = "/"
	}

	urlHostPrefix = hostPrefix
	if hostPrefix == "/" {
		urlHostPrefix = ""
	}
}

func GetHostPrefix() string {
	return hostPrefix
}

func BuildHostURL(rest string) string {
	if strings.HasPrefix(rest, "/") {
		rest = strings.TrimPrefix(rest, "/")
	}

	return urlHostPrefix + "/" + rest
}

func BuildSafeURL(rest string) templ.SafeURL {
	return templ.SafeURL(BuildHostURL(rest))
}
