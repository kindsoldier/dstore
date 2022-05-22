/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package xtools

import (
    "errors"
    "github.com/google/uuid"
)

const UUIDSize int = 16

type UUIDBytes = [UUIDSize]byte

func UUIDStringToBytes(foo string) (UUIDBytes, error) {
    var err error
    var uuidBytes UUIDBytes
    uuid, err := uuid.Parse(foo)
    if err != nil {
        return uuidBytes, err
    }
    uuidSlice, err := uuid.MarshalBinary()
    if err != nil {
        return uuidBytes, err
    }
    if len(uuidSlice) < UUIDSize {
        return uuidBytes, errors.New("uuid array too short")
    }
    for i := 0; i < UUIDSize; i++ {
        uuidBytes[i] = uuidSlice[i]
    }
    return uuidBytes, err
}

func UUIDBytesToString(uuidBytes UUIDBytes) (string, error) {
    var err error
    var uuidString string
    uuidSlice := make([]byte, UUIDSize)
    for i := 0; i < UUIDSize; i++ {
        uuidSlice[i] = uuidBytes[i]
    }
    uuid, err := uuid.FromBytes(uuidSlice)
    if err != nil {
        return uuidString, err
    }
    uuidString = uuid.String()
    return uuidString, err
}
