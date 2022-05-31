
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const UpdateBStoreMethod string = "updateBStore"

type UpdateBStoreParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
}

type UpdateBStoreResult struct {
}

func NewUpdateBStoreResult() *UpdateBStoreResult {
    return &UpdateBStoreResult{}
}
func NewUpdateBStoreParams() *UpdateBStoreParams {
    return &UpdateBStoreParams{}
}
