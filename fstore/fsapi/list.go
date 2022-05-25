
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "ndstore/dscom"
)

const ListMethod string = "list"

type ListParams struct {
    DirPath     string      `json:"dirPath"`
}

type ListResult struct {
    Files   []*dscom.DirEntry   `json:"cFiles"`
}

func NewListResult() *ListResult {
    return &ListResult{}
}
func NewListParams() *ListParams {
    return &ListParams{}
}
