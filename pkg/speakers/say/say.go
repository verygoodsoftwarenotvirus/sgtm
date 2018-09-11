package say

import (
	"os/exec"
)

const (
	defaultLanguage = "en"
	defaultVoice    = "alex"
)

var (
	// acceptableVoices is the set of voices currently supported by SGTM. To see all voices available using Say, execute "say --voice=?"
	englishVoices = map[string]string{
		defaultVoice: "Alex",
		"daniel":     "Daniel",
		"fiona":      "Fiona",
		"fred":       "Fred",
		"samantha":   "Samantha",
		"victoria":   "Victoria",
	}
	spanishVoices = map[string]string{
		"diego":   "Diego",
		"jorge":   "Jorge",
		"juan":    "Juan",
		"monica":  "Monica",
		"paulina": "Paulina",
	}
	acceptableVoices = map[string]map[string]string{
		"en": englishVoices,
		"es": spanishVoices,
	}
)

type SaySpeaker struct {
	Language string
	VoiceID  string
}

// New takes in a language and voice to create a SaySpeaker instance
func New(language, voiceName string) *SaySpeaker {
	var voice string
	voice, ok := acceptableVoices[language][voiceName]
	if !ok {
		language = defaultLanguage
		voice = acceptableVoices[defaultLanguage][defaultVoice]
	}
	return &SaySpeaker{
		Language: language,
		VoiceID:  voice,
	}
}

// GenerateSpeech takes in a string and outputs Synthesized Speech and possibly an error.
func (ss *SaySpeaker) GenerateSpeech(text, fileName string) error {
	cmd := exec.Command("say", text)
	error := cmd.Run()
	if error != nil {
		return error
	}
	return nil
}
