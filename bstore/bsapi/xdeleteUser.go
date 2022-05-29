
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const DeleteUserMethod string = "deleteUser"

type DeleteUserParams struct {
    Login      string           `json:"login"`
    Pass     string           `json:"pass"`
    State     string           `json:"state"`
}

type DeleteUserResult struct {
}

func NewDeleteUserResult() *DeleteUserResult {
    return &DeleteUserResult{}
}
func NewDeleteUserParams() *DeleteUserParams {
    return &DeleteUserParams{}
}
