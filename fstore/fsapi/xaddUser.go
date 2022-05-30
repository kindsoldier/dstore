
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const AddUserMethod string = "addUser"

type AddUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
}

type AddUserResult struct {
}

func NewAddUserResult() *AddUserResult {
    return &AddUserResult{}
}
func NewAddUserParams() *AddUserParams {
    return &AddUserParams{}
}
