package docker

import "strings"

func parseImageVersion(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		return parts[1]
	}

	return "latest"
}
