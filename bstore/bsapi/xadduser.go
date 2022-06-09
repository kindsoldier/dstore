
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const AddUserMethod string = "addUser"

type AddUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
    Role    string      `json:"role"`
}

type AddUserResult struct {
}

func NewAddUserResult() *AddUserResult {
    return &AddUserResult{}
}
func NewAddUserParams() *AddUserParams {
    return &AddUserParams{}
}
