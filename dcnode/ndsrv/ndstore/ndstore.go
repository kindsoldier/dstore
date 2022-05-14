/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndstore

import (
    "io"
    "dcstore/dcrpc"
)

type Store struct {
}

func NewStore() *Store {
    return &Store{}
}


func (store *Store) SaveBlock(blockReader io.Reader, blockSize int64) error {
    var err error

    _, err = dcrpc.ReadBytes(blockReader, blockSize)

    return err
}
