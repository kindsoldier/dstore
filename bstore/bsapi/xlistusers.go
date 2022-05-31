
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

import (
    "ndstore/bstore/bscom"
)

const ListUsersMethod string = "listUsers"

type ListUsersParams struct {
}

type ListUsersResult struct {
    Users  []*bscom.UserDescr     `json:"users,omitempty"`
}

func NewListUsersResult() *ListUsersResult {
    return &ListUsersResult{}
}
func NewListUsersParams() *ListUsersParams {
    return &ListUsersParams{}
}
