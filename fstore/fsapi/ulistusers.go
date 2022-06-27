
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "ndstore/dscom"
)

const ListUsersMethod string = "listUsers"

type ListUsersParams struct {
}

type ListUsersResult struct {
    Users  []*dscom.UserDescr     `json:"users"`
}

func NewListUsersResult() *ListUsersResult {
    return &ListUsersResult{}
}
func NewListUsersParams() *ListUsersParams {
    return &ListUsersParams{}
}
