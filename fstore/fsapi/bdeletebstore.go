
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const DeleteBStoreMethod string = "deleteBStore"

type DeleteBStoreParams struct {
    Login      string           `json:"login"`
}

type DeleteBStoreResult struct {
}

func NewDeleteBStoreResult() *DeleteBStoreResult {
    return &DeleteBStoreResult{}
}
func NewDeleteBStoreParams() *DeleteBStoreParams {
    return &DeleteBStoreParams{}
}
