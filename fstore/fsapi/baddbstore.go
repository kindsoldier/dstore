
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const AddBStoreMethod string = "addBStore"

type AddBStoreParams struct {
    Address string      `json:"address" db:"address"`
    Port    string      `json:"port"    db:"port"`
    Login   string      `json:"login"   db:"login"`
    Pass    string      `json:"pass"    db:"pass"`
    State   string      `json:"state"   db:"state"`
}

type AddBStoreResult struct {
}

func NewAddBStoreResult() *AddBStoreResult {
    return &AddBStoreResult{}
}
func NewAddBStoreParams() *AddBStoreParams {
    return &AddBStoreParams{}
}
