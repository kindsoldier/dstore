
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const DeleteBStoreMethod string = "deleteBStore"

type DeleteBStoreParams struct {
    Address string      `json:"address" db:"address"`
    Port    string      `json:"port"    db:"port"`
}

type DeleteBStoreResult struct {
}

func NewDeleteBStoreResult() *DeleteBStoreResult {
    return &DeleteBStoreResult{}
}
func NewDeleteBStoreParams() *DeleteBStoreParams {
    return &DeleteBStoreParams{}
}
