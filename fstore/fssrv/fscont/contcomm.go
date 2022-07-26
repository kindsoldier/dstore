/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fscont

import (
    "dstore/fstore/fssrv/fstore"
)


type Contr struct {
    store  *fstore.Store
}

func NewContr(store *fstore.Store) (*Contr, error) {
    var err error
    var contr Contr
    contr.store = store
    return &contr, err
}

