/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZipUnzip(t *testing.T) {
	iBytes := []byte("ЙЦУКЕНГqwerty1234567890")
	zipBytes, err := Zip(iBytes)
        assert.NoError(t, err)

	oBytes, err := Unzip(zipBytes)
        assert.NoError(t, err)
	assert.Equal(t, iBytes, oBytes, nil)
}
