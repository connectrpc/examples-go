// Copyright 2020-2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eliza

import (
	"fmt"
	"math/rand"
	"strings"
)

// ReplyTo will construct a reply for a given input using ELIZA's rules.
func ReplyTo(input string) (string, bool) {
	input = preprocess(input)
	if _, ok := goodbyeInputSet[input]; ok {
		return randomElementFrom(goodbyeResponses), true
	}
	return lookupResponse(input), false
}

// lookupResponse does a lookup with regex
func lookupResponse(input string) string {
	// Look up responses from requestInputRegexToResponseOptions mapping
	for re, responses := range requestInputRegexToResponseOptions {
		matches := re.FindStringSubmatch(input)
		if len(matches) < 1 {
			continue
		}
		// Select a random response
		response := randomElementFrom(responses)
		// We attempt to reflect a response phrase, when the response has an entry point
		if !strings.Contains(response, "%s") {
			return response
		}
		if len(matches) > 1 {
			var fragment string
			fragment = reflect(matches[1])
			response = fmt.Sprintf(response, fragment)
			return response
		}
	}
	return randomElementFrom(defaultResponses)
}

// preprocess will do some normalization on a statement for better regex matching
func preprocess(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)
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

// randomElementFrom returns a random element in the input array.
func randomElementFrom(list []string) string {
	random := rand.Intn(len(list))
	return list[random]
}
