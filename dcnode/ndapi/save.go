
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndapi

const SaveMethod string = "save"

type SaveParams struct {
    ClusterId   int64           `json:"clusterId"`
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
