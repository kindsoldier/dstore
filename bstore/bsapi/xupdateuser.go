
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const UpdateUserMethod string = "updateUser"

type UpdateUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
    Role    string      `json:"role"`
}

type UpdateUserResult struct {
}

func NewUpdateUserResult() *UpdateUserResult {
    return &UpdateUserResult{}
}
func NewUpdateUserParams() *UpdateUserParams {
    return &UpdateUserParams{}
}
