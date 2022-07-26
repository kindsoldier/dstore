/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bscont

import (
    "dstore/bstore/bssrv/bstore"
)


type Contr struct {
    store  *bstore.Store
}

func NewContr(store *bstore.Store) (*Contr, error) {
    var err error
    var contr Contr
    contr.store = store
    return &contr, err
}
