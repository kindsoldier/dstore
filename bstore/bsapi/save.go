
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const SaveMethod string = "save"

type SaveParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
}

type SaveResult struct {
}

func NewSaveResult() *SaveResult {
    return &SaveResult{}
}
func NewSaveParams() *SaveParams {
    return &SaveParams{}
}
