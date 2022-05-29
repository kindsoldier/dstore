
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const CheckUserMethod string = "checkUser"

type CheckUserParams struct {
    Login   string      `json:"login"`
    Pass    string      `json:"pass"`
}

type CheckUserResult struct {
}

func NewCheckUserResult() *CheckUserResult {
    return &CheckUserResult{}
}
func NewCheckUserParams() *CheckUserParams {
    return &CheckUserParams{}
}
