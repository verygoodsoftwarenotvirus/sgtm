package say

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaySpeaker_New(t *testing.T) {
	t.Parallel()

	validLanguage := "en"
	invalidLanguage := "ja"
	validVoice := "daniel"
	invalidVoice := "MrBad"

	assert := assert.New(t)

	// Test valid language, invalid voice combination
	ss := New(validLanguage, invalidVoice)
	assert.Equal(ss.VoiceID, "Alex", "They should be equal")

	// Test invalid language, valid voice combination
	ss = New(invalidLanguage, validVoice)
	assert.Equal(ss.Language, "en", "They should be equal")

	// Test valid language, valid voice combination
	ss = New(validLanguage, validVoice)
	assert.Equal(ss.VoiceID, "Daniel", "They should be equal")
	assert.Equal(ss.Language, validLanguage, "They should be equal")
}

func TestSaySpeaker_GenerateSpeech(t *testing.T) {
	t.Parallel()
	testText := "This is a test!"
	ss := New("daniel", "en")
	err := ss.GenerateSpeech(testText, "")
	assert.NoError(t, err)
}
