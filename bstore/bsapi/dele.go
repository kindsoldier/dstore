
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const DeleteMethod string = "delete"

type DeleteParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
}

type DeleteResult struct {
}

func NewDeleteResult() *DeleteResult {
    return &DeleteResult{}
}
func NewDeleteParams() *DeleteParams {
    return &DeleteParams{}
}
