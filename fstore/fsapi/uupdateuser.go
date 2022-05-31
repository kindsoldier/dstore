
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const UpdateUserMethod string = "updateUser"

type UpdateUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
    State   string      `json:"state"`
}

type UpdateUserResult struct {
}

func NewUpdateUserResult() *UpdateUserResult {
    return &UpdateUserResult{}
}
func NewUpdateUserParams() *UpdateUserParams {
    return &UpdateUserParams{}
}
