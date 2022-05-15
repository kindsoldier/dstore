/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "encoding/hex"
    "time"
    "math/rand"
)

func Raw2HexBytes(rawBytes []byte) []byte {
    hexBytes := make([]byte, hex.EncodedLen(len(rawBytes)))
    hex.Encode(hexBytes, rawBytes)
    return hexBytes
}


func Raw2HexString(rawBytes []byte) string {
    hexBytes := make([]byte, hex.EncodedLen(len(rawBytes)))
    hex.Encode(hexBytes, rawBytes)
    return string(hexBytes)
}

func Hex2Raw(hexBytes []byte) ([]byte, error) {
    var err error
    rawBytes := make([]byte, hex.DecodedLen(len(hexBytes)))
    _, err = hex.Decode(rawBytes, hexBytes)
    if err != nil {
        return rawBytes, err
    }
    return rawBytes, err
}

func RandBytesHexString(size int) string {
    rand.Seed(time.Now().UnixNano())
    randBytes := make([]byte, size)
    rand.Read(randBytes)
    hexString := hex.EncodeToString(randBytes)
    return hexString
}
