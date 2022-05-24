// Lightly modified from https://github.com/mattshiel/eliza-go

// Package eliza is a simple (and not very convincing) simulation of a
// psychotherapist. It emulates the DOCTOR script written for Joseph
// Weizenbaum's 1966 ELIZA natural language processing system.
package eliza

import (
	"fmt"
	"math/rand"
	"strings"
)

// Reply responds to a statement as a pyschotherapist might.
func Reply(input string) (string, bool) {
	input = preprocess(input)
	if _, ok := goodbyeInputSet[input]; ok {
		return randomElementFrom(goodbyeResponses), true
	}
	return lookupResponse(input), false
}

func lookupResponse(input string) string {
	for re, responses := range requestInputRegexToResponseOptions {
		matches := re.FindStringSubmatch(input)
		if len(matches) < 1 {
			continue
		}
		response := randomElementFrom(responses)
		// If the response has an entry point, reflect the input phrase (so "I"
		// becomes "you").
		if !strings.Contains(response, "%s") {
			return response
		}
		if len(matches) > 1 {
			fragment := reflect(matches[1])
			response = fmt.Sprintf(response, fragment)
			return response
		}
	}
	return randomElementFrom(defaultResponses)
}

func preprocess(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
	input = strings.Trim(input, `.!?'"`)
	return input
}

// reflect flips a few words in an input fragment (such as "I" -> "you").
func reflect(fragment string) string {
	words := strings.Fields(fragment)
	for i, word := range words {
		if reflectedWord, ok := reflectedWords[word]; ok {
			words[i] = reflectedWord
		}
	}
	return strings.Join(words, " ")
}

func randomElementFrom(list []string) string {
	random := rand.Intn(len(list)) // nolint:gosec
	return list[random]
}
