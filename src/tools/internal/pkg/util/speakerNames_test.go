package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSpeakers(t *testing.T) {
	assert.Equal(t, map[string]string{}, ParseSpeakers(""))
	assert.Equal(t, map[string]string{}, ParseSpeakers("a=a;b=b"))
	assert.Equal(t, map[string]string{"1": "name oth"}, ParseSpeakers("1=audio_only_121_name_oth.mp4;b=b"))
	assert.Equal(t, map[string]string{"1": "name oth", "2": "name"},
		ParseSpeakers("1=audio_only_121_name_oth.mp4;2=audio_only_121_name.mp4"))
	assert.Equal(t, map[string]string{"1": "name oth"},
		ParseSpeakers("1=audio_only_121_name_oth.mp4;2=only_121_name.mp4"))
}

func TestGetSpeakerByPath(t *testing.T) {
	mp := ParseSpeakers("1=audio_only_121_name_oth.mp4;2=audio_only_121_name.mp4")
	assert.Equal(t, "", GetSpeakerByPath(mp, ""))
	assert.Equal(t, "name oth", GetSpeakerByPath(mp, "1/olia.lat"))
	assert.Equal(t, "name", GetSpeakerByPath(mp, "2/olia.lat"))
	assert.Equal(t, "", GetSpeakerByPath(mp, "3/olia.lat"))
}
