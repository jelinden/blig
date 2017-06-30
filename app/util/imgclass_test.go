package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImgClass(t *testing.T) {
	test := ImgClass("test string without image")
	assert.False(t, strings.Contains(test, "class"), "should not contain class")

	test2 := ImgClass("test string with image <img src=\"test\"/> ")
	assert.True(t, strings.Contains(test2, "class=\"pure-img\""), "should contain class")

	test3 := ImgClass("test string with 2 images <img src=\"test\"/> <img src=\"test\"/>")
	assert.True(t, strings.Contains(test3, "class=\"pure-img\"") && strings.Count(test3, "class=\"pure-img\"") == 2, "should contain class")
}
