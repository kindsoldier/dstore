/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRaw2Hex2Raw(t *testing.T) {
    iRawBytes := []byte("ЙЦУКЕНГqwerty1234567890")
    hexBytes := Raw2HexBytes(iRawBytes)

    oRawBytes, err := Hex2Raw(hexBytes)
    assert.NoError(t, err)

    assert.Equal(t, iRawBytes, oRawBytes, nil)
}
