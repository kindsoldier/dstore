
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdapi

import (
    "ndstore/dscom"
)

const ListMethod string = "list"

type ListParams struct {
    ClusterId   int64           `json:"clusterId"`
}

type ListResult struct {
    Blocks      []dscom.Block   `json:"blocks"`
}

func NewListResult() *ListResult {
    return &ListResult{}
}
func NewListParams() *ListParams {
    return &ListParams{}
}
