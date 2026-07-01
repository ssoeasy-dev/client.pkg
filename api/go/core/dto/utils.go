package dto

import "strings"

const BearerPrefix = "Bearer "

func ExtractBearer(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], BearerPrefix) {
		return parts[1]
	}
	return ""
}
