// Copyright 2022 Buf Technologies, Inc.
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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplyToGoodbyes(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"bye", "quit", "exit", "goodbye"} {
		response, end := Reply(input)
		assert.True(t, end)
		assert.Contains(t, goodbyeResponses, response)
	}
}

func TestDefaultAnswers(t *testing.T) {
	t.Parallel()
	for i := 0; i < 3; i++ {
		response, _ := Reply("i have" + strings.Repeat(" ", i))
		assert.Contains(t, defaultResponses, response)
	}
}

func TestHello(t *testing.T) {
	t.Parallel()
	response, _ := Reply("hello eliza!")
	assert.Contains(t, response, "Hello")

	response, _ = Reply("hello there")
	assert.Contains(t, response, "Hello")
}

func TestReflectiveAnswers(t *testing.T) {
	t.Parallel()
	response, _ := Reply("i have a problem")
	assert.Contains(t, response, "a problem")

	response, _ = Reply("i have a problem with your tone")
	assert.Contains(t, response, "a problem with my tone")
}
