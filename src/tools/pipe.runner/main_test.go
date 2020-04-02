package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlesResultRequired(t *testing.T) {
	d := &workingData{}
	d.args = []string{"a", "b"}
	d.outPipe = "pOut"
	assert.Equal(t, "a b pOut", getWriteData(d))
}
