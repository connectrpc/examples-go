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
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestReplyToGoodbyes(t *testing.T) {
	response, end := ReplyTo("bye")
	assert.True(t, end)
	assert.True(t, contains(goodbyeResponses, response))

	response, end = ReplyTo("quit")
	assert.True(t, end)
	assert.True(t, contains(goodbyeResponses, response))

	response, end = ReplyTo("exit")
	assert.True(t, end)
	assert.True(t, contains(goodbyeResponses, response))

	response, end = ReplyTo("goodbye")
	assert.True(t, end)
	assert.True(t, contains(goodbyeResponses, response))
}

func TestHello(t *testing.T) {
	rand.Seed(1234) // set random seed to pin responses
	response, _ := ReplyTo("hello eliza!")
	assert.Equal(t, "Hello, how are you feeling today?", response)

	response, _ = ReplyTo("hello there")
	assert.Equal(t, "Hello, how are you feeling today?", response)
}

func TestReflectiveAnswers(t *testing.T) {
	rand.Seed(1234) // set random seed to pin responses
	response, _ := ReplyTo("i have")
	assert.True(t, contains(defaultResponses, response))
	response, _ = ReplyTo("i have ")
	assert.True(t, contains(defaultResponses, response))
	response, _ = ReplyTo("i have  ")
	assert.True(t, contains(defaultResponses, response))
	response, _ = ReplyTo("i have a problem")
	assert.Equal(t, response, "Why do you tell me that you've a problem?")
	response, _ = ReplyTo("i have a problem with your tone")
	assert.Equal(t, response, "Have you really a problem with my tone?")
}

func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
