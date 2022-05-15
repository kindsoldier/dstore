/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "encoding/hex"
    "errors"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func RandBytesRaw(size int) []byte {
    randBytes := make([]byte, size)
    rand.Read(randBytes)
    return randBytes
}

func RandBytesRaw64(size int64) []byte {
    randBytes := make([]byte, size)
    rand.Read(randBytes)
    return randBytes
}


func RandBytesHex(size int) []byte {
    randBuffer := make([]byte, size/2)
    rand.Read(randBuffer)
    hexBuffer := make([]byte, hex.EncodedLen(size))
    hex.Encode(hexBuffer, randBuffer)
    return hexBuffer[0:size]
}

func RandInt(min, max int) (int, error) {
    var err error
    var theRand int
    theRange := max - min
    if theRange < 0 {
        return theRand, errors.New("max less min")
    }
    if theRange == 0 {
        return min, err
    }

    return rand.Intn(theRange) + min, err
}

func RandInt64(min, max int64) (int64, error) {
    var err error
    var theRand int64
    theRange := max - min
    if theRange < 0 {
        return theRand, errors.New("max less min")
    }
    if theRange == 0 {
        return min, err
    }
    return rand.Int63n(theRange) + min, err
}
