
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

import (
    "ndstore/dscom"
)

const ListMethod string = "list"

type ListParams struct {
}

type ListResult struct {
    Blocks  []dscom.BlockMI     `json:"blocks"`
}

func NewListResult() *ListResult {
    return &ListResult{}
}
func NewListParams() *ListParams {
    return &ListParams{}
}
