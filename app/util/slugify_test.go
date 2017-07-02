package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	test := Slugify("test string")
	assert.True(t, strings.EqualFold(test, "test-string"), "should be test-string")

	test2 := Slugify("Cache on the background with golang, Part 2")
	assert.True(t, strings.EqualFold(test2, "cache-on-the-background-with-golang-part-2"), "cache-on-the-background-with-golang-part-2")
}
