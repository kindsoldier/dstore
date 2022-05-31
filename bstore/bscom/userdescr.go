/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bscom

type UserDescr struct {
    Login   string      `json:"login"   db:"login"`
    Pass    string      `json:"pass"    db:"pass"`
    State   string      `json:"state"   db:"state"`
    Role    string      `json:"role"    db:"role"`
}

func NewUserDescr() *UserDescr {
    return &UserDescr{}
}
