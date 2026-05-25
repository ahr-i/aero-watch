package corsController

import (
	"strings"

	"github.com/rs/cors"
)

func SetCors(origins string, methods string, headers string, credentials bool) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   splitValues(origins),
		AllowedMethods:   splitValues(methods),
		AllowedHeaders:   splitValues(headers),
		AllowCredentials: credentials,
	})
}

func splitValues(values string) []string {
	result := []string{}
	for _, value := range strings.Split(values, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
