// Lightly modified from https://github.com/mattshiel/eliza-go
package eliza

import (
	"fmt"
	"math/rand"
	"strings"
)

type Eliza interface {
	Introduction() string
	ReplyTo(input string) string
}

type eliza struct{}

// NewEliza get an instance to converse with Eliza
func NewEliza() Eliza {
	return &eliza{}
}

// Introduction will return a random introductory sentence for ELIZA.
func (e *eliza) Introduction() string {
	return randChoice(introductions)
}

// ReplyTo will construct a reply for a given input using ELIZA's rules.
func (e *eliza) ReplyTo(input string) string {
	input = preprocess(input)
	if isQuitStatement(input) {
		return randChoice(goodbyes)
	}

	if response, ok := lookupResponse(input); ok {
		return response
	}

	// If no patterns were matched, return a default response.
	return randChoice(defaultResponses)
}

// lookupResponse does a lookup with regex
func lookupResponse(input string) (string, bool) {
	// Look up responses from psychobabble mapping
	for re, responses := range psychobabble {
		matches := re.FindStringSubmatch(input)
		if len(matches) > 0 {
			var fragment string
			if len(matches) > 1 {
				fragment = reflect(matches[1])
			}
			response := randChoice(responses)
			if strings.Contains(response, "%s") {
				response = fmt.Sprintf(response, fragment)
			}
			return response, true
		}
	}
	return "", false
}

// isQuitStatement returns if the statement is a quit statement
func isQuitStatement(statement string) bool {
	statement = preprocess(statement)
	for _, quitStatement := range quitStatements {
		if statement == quitStatement {
			return true
		}
	}
	return false
}

// preprocess will do some normalization on a statement for better regex matching
func preprocess(input string) string {
	input = strings.TrimRight(input, "\n.!")
	input = strings.ToLower(input)
	return input
}

// reflect flips a few words in an input fragment (such as "I" -> "you").
func reflect(fragment string) string {
	words := strings.Split(fragment, " ")
	for i, word := range words {
		if reflectedWord, ok := reflectedWords[word]; ok {
			words[i] = reflectedWord
		}
	}
	return strings.Join(words, " ")
}

// randChoice returns a random element in the input array.
func randChoice(list []string) string {
	randIndex := rand.Intn(len(list))
	return list[randIndex]
}
