
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

import (
    "ndstore/dscom"
)

const ListFilesMethod string = "listFiles"

type ListFilesParams struct {
    DirPath     string              `json:"dirPath"`
}

type ListFilesResult struct {
    Entries   []*dscom.EntryDescr   `json:"entries"`
}

func NewListFilesResult() *ListFilesResult {
    return &ListFilesResult{}
}
func NewListFilesParams() *ListFilesParams {
    return &ListFilesParams{}
}
