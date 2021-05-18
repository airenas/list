package util

import (
	"path/filepath"
	"strings"
)

func ParseSpeakers(s string) map[string]string {
	res := make(map[string]string)
	strs := strings.Split(s, ";")
	for _, idSp := range strs {
		sParsed := strings.Split(idSp, "=")
		if len(sParsed) == 2 {
			id := strings.TrimSpace(sParsed[0])
			sp := extractSpeaker(strings.TrimSpace(sParsed[1]))
			if id != "" && sp != "" {
				res[id] = sp
			}
		}
	}
	return res
}

func GetSpeakerByPath(idSpMap map[string]string, path string) string {
	for k, v := range idSpMap {
		if strings.Contains(path, k) {
			return v
		}
	}
	return ""
}

func extractSpeaker(f string) string {
	fb := filepath.Base(f)
	fb = strings.TrimSuffix(fb, filepath.Ext(fb))
	if strings.HasPrefix(fb, "audio_only_") {
		strs := strings.Split(fb, "_")
		if len(strs) > 3 {
			return strings.Join(strs[3:], " ")
		}
	}
	return ""
}
