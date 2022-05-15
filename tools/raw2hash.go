/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "encoding/hex"
    "crypto/sha256"
)

func Raw2HashBytes(rawBytes []byte) []byte {
    hasher := sha256.New()
    hasher.Write(rawBytes)
    hashBytes := hasher.Sum(nil)
    hexBytes := make([]byte, hex.EncodedLen(len(hashBytes)))
    hex.Encode(hexBytes, hashBytes)
    return hexBytes
}

func Raw2HashString(rawBytes []byte) string {
    hasher := sha256.New()
    hasher.Write(rawBytes)
    hashBytes := hasher.Sum(nil)
    hexBytes := make([]byte, hex.EncodedLen(len(hashBytes)))
    hex.Encode(hexBytes, hashBytes)
    return string(hexBytes)
}
