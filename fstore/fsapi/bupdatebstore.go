
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const UpdateBStoreMethod string = "updateBStore"

type UpdateBStoreParams struct {
    Address string      `json:"address" db:"address"`
    Port    string      `json:"port"    db:"port"`
    Login   string      `json:"login"   db:"login"`
    Pass    string      `json:"pass"    db:"pass"`
    State   string      `json:"state"   db:"state"`
}

type UpdateBStoreResult struct {
}

func NewUpdateBStoreResult() *UpdateBStoreResult {
    return &UpdateBStoreResult{}
}
func NewUpdateBStoreParams() *UpdateBStoreParams {
    return &UpdateBStoreParams{}
}
