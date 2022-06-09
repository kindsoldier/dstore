
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const LoadBlockMethod string = "loadBlock"

type LoadBlockParams struct {
    BlockType   string          `json:"blockType"`
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
}

type LoadBlockResult struct {
}

func NewLoadBlockResult() *LoadBlockResult {
    return &LoadBlockResult{}
}
func NewLoadBlockParams() *LoadBlockParams {
    return &LoadBlockParams{}
}
