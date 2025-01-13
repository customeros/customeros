package utils

import (
	"github.com/lucasepe/codename"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

func toCamelCase(input string) string {
	splitName := strings.Split(input, "-")

	for i := 0; i < len(splitName); i++ {
		splitName[i] = cases.Title(language.English).String(splitName[i])
	}

	camelCaseName := strings.Join(splitName, "")
	return camelCaseName
}

func Sanitize(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	return reg.ReplaceAllString(input, "")
}

func GenerateName() string {
	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	name := codename.Generate(rng, 0)
	return toCamelCase(name)
}
