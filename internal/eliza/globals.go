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

// Copied from https://github.com/mattshiel/eliza-go and modified.
//
// See https://github.com/mattshiel/eliza-go/blob/master/LICENSE.

package eliza

import "regexp"

// Input statements which terminate the session.
var goodbyeInputSet = map[string]struct{}{
	"bye":     {},
	"exit":    {},
	"goodbye": {},
	"quit":    {},
}

// End-of-session responses.
var goodbyeResponses = []string{
	"Goodbye. It was nice talking to you.",
	"Thank you for talking with me.",
	"Thank you, that will be $150. Have a good day!",
	"Goodbye. This was really a nice talk.",
	"Goodbye. I'm looking forward to our next session.",
	"This was a good session, wasn't it – but time is over now. Goodbye.",
	"Maybe we could discuss this over more in our next session? Goodbye.",
	"Good-bye.",
}

// Request phrase to response phrases as a lookup table.
var requestInputRegexToResponseOptions = map[*regexp.Regexp][]string{
	regexp.MustCompile(`i need (.*)`): {
		"Why do you need %s?",
		"Would it really help you to get %s?",
		"Are you sure you need %s?",
	},
	regexp.MustCompile(`why don'?t you ([^\?]*)\??`): {
		"Do you really think I don't %s?",
		"Perhaps eventually I will %s.",
		"Do you really want me to %s?",
	},
	regexp.MustCompile(`why can'?t I ([^\?]*)\??`): {
		"Do you think you should be able to %s?",
		"If you could %s, what would you do?",
		"I don't know -- why can't you %s?",
		"Have you really tried?",
	},
	regexp.MustCompile(`i can'?t (.*)`): {
		"How do you know you can't %s?",
		"Perhaps you could %s if you tried.",
		"What would it take for you to %s?",
	},
	regexp.MustCompile(`i am (.*)`): {
		"Did you come to me because you are %s?",
		"How long have you been %s?",
		"How do you feel about being %s?",
	},
	regexp.MustCompile(`i'?m (.*)`): {
		"How does being %s make you feel?",
		"Do you enjoy being %s?",
		"Why do you tell me you're %s?",
		"Why do you think you're %s?",
	},
	regexp.MustCompile(`are you ([^\?]*)\??`): {
		"Why does it matter whether I am %s?",
		"Would you prefer it if I were not %s?",
		"Perhaps you believe I am %s.",
		"I may be %s -- what do you think?",
	},
	regexp.MustCompile(`what (.*)`): {
		"Why do you ask?",
		"How would an answer to that help you?",
		"What do you think?",
	},
	regexp.MustCompile(`how (.*)`): {
		"How do you suppose?",
		"Perhaps you can answer your own question.",
		"What is it you're really asking?",
	},
	regexp.MustCompile(`because (.*)`): {
		"Is that the real reason?",
		"What other reasons come to mind?",
		"Does that reason apply to anything else?",
		"If %s, what else must be true?",
	},
	regexp.MustCompile(`(.*) sorry (.*)`): {
		"There are many times when no apology is needed.",
		"What feelings do you have when you apologize?",
	},
	regexp.MustCompile(`^hello(.*)`): {
		"Hello...I'm glad you could drop by today.",
		"Hello there...how are you today?",
		"Hello, how are you feeling today?",
	},
	regexp.MustCompile(`^hi(.*)`): {
		"Hello...I'm glad you could drop by today.",
		"Hi there...how are you today?",
		"Hello, how are you feeling today?",
	},
	regexp.MustCompile(`^thanks(.*)`): {
		"You're welcome!",
		"Anytime!",
	},
	regexp.MustCompile(`^thank you(.*)`): {
		"You're welcome!",
		"Anytime!",
	},
	regexp.MustCompile(`^good morning(.*)`): {
		"Good morning...I'm glad you could drop by today.",
		"Good morning...how are you today?",
		"Good morning, how are you feeling today?",
	},
	regexp.MustCompile(`^good afternoon(.*)`): {
		"Good afternoon...I'm glad you could drop by today.",
		"Good afternoon...how are you today?",
		"Good afternoon, how are you feeling today?",
	},
	regexp.MustCompile(`I think (.*)`): {
		"Do you doubt %s?",
		"Do you really think so?",
		"But you're not sure %s?",
	},
	regexp.MustCompile(`(.*) friend (.*)`): {
		"Tell me more about your friends.",
		"When you think of a friend, what comes to mind?",
		"Why don't you tell me about a childhood friend?",
	},
	regexp.MustCompile(`yes`): {
		"You seem quite sure.",
		"OK, but can you elaborate a bit?",
	},
	regexp.MustCompile(`(.*) computer(.*)`): {
		"Are you really talking about me?",
		"Does it seem strange to talk to a computer?",
		"How do computers make you feel?",
		"Do you feel threatened by computers?",
	},
	regexp.MustCompile(`is it (.*)`): {
		"Do you think it is %s?",
		"Perhaps it's %s -- what do you think?",
		"If it were %s, what would you do?",
		"It could well be that %s.",
	},
	regexp.MustCompile(`it is (.*)`): {
		"You seem very certain.",
		"If I told you that it probably isn't %s, what would you feel?",
	},
	regexp.MustCompile(`can you ([^\?]*)\??`): {
		"What makes you think I can't %s?",
		"If I could %s, then what?",
		"Why do you ask if I can %s?",
	},
	regexp.MustCompile(`(.*)dream(.*)`): {
		"Tell me more about your dream.",
	},
	regexp.MustCompile(`can I ([^\?]*)\??`): {
		"Perhaps you don't want to %s.",
		"Do you want to be able to %s?",
		"If you could %s, would you?",
	},
	regexp.MustCompile(`you are (.*)`): {
		"Why do you think I am %s?",
		"Does it please you to think that I'm %s?",
		"Perhaps you would like me to be %s.",
		"Perhaps you're really talking about yourself?",
	},
	regexp.MustCompile(`you'?re (.*)`): {
		"Why do you say I am %s?",
		"Why do you think I am %s?",
		"Are we talking about you, or me?",
	},
	regexp.MustCompile(`i don'?t (.*)`): {
		"Don't you really %s?",
		"Why don't you %s?",
		"Do you want to %s?",
	},
	regexp.MustCompile(`i feel (.*)`): {
		"Good, tell me more about these feelings.",
		"Do you often feel %s?",
		"When do you usually feel %s?",
		"When you feel %s, what do you do?",
		"Feeling %s? Tell me more.",
	},
	regexp.MustCompile(`i have (.*)`): {
		"Why do you tell me that you've %s?",
		"Have you really %s?",
		"Now that you have %s, what will you do next?",
	},
	regexp.MustCompile(`i would (.*)`): {
		"Could you explain why you would %s?",
		"Why would you %s?",
		"Who else knows that you would %s?",
	},
	regexp.MustCompile(`is there (.*)`): {
		"Do you think there is %s?",
		"It's likely that there is %s.",
		"Would you like there to be %s?",
	},
	regexp.MustCompile(`my (.*)`): {
		"I see, your %s.",
		"Why do you say that your %s?",
		"When your %s, how do you feel?",
	},
	regexp.MustCompile(`you (.*)`): {
		"We should be discussing you, not me.",
		"Why do you say that about me?",
		"Why do you care whether I %s?",
	},
	regexp.MustCompile(`why (.*)`): {
		"Why don't you tell me the reason why %s?",
		"Why do you think %s?",
	},
	regexp.MustCompile(`i want (.*)`): {
		"What would it mean to you if you got %s?",
		"Why do you want %s?",
		"What would you do if you got %s?",
		"If you got %s, then what would you do?",
	},
	regexp.MustCompile(`(.*) mother(.*)`): {
		"Tell me more about your mother.",
		"What was your relationship with your mother like?",
		"How do you feel about your mother?",
		"How does this relate to your feelings today?",
		"Good family relations are important.",
	},
	regexp.MustCompile(`(.*) father(.*)`): {
		"Tell me more about your father.",
		"How did your father make you feel?",
		"How do you feel about your father?",
		"Does your relationship with your father relate to your feelings today?",
		"Do you have trouble showing affection with your family?",
	},
	regexp.MustCompile(`(.*) child(.*)`): {
		"Did you have close friends as a child?",
		"What is your favorite childhood memory?",
		"Do you remember any dreams or nightmares from childhood?",
		"Did the other children sometimes tease you?",
		"How do you think your childhood experiences relate to your feelings today?",
	},
	regexp.MustCompile(`(.*)\?`): {
		"Why do you ask that?",
		"Please consider whether you can answer your own question.",
		"Perhaps the answer lies within yourself?",
		"Why don't you tell me?",
	},
}

// Default responses when nothing more specific applies.
var defaultResponses = []string{
	"Please tell me more.",
	"Let's change focus a bit...Tell me about your family.",
	"Can you elaborate on that?",
	"I see.",
	"Very interesting.",
	"I see. And what does that tell you?",
	"How does that make you feel?",
	"How do you feel when you say that?",
}

// A table to reflect words in question fragments inside the response.
// For example, the phrase "your jacket" in "I want your jacket" should be
// reflected to "my jacket" in the response.
var reflectedWords = map[string]string{
	"am":     "are",
	"was":    "were",
	"i":      "you",
	"i'd":    "you would",
	"i've":   "you have",
	"i'll":   "you will",
	"my":     "your",
	"are":    "am",
	"you've": "I have",
	"you'll": "I will",
	"your":   "my",
	"yours":  "mine",
	"you":    "me",
	"me":     "you",
}
