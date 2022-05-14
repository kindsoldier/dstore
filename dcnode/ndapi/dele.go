
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndapi

const DeleteMethod string = "delete"

type DeleteParams struct {
    ClusterId   int64           `json:"clusterId"`
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

