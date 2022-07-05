
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const CheckBlockMethod string = "checkBlock"

type CheckBlockParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    BlockVer    int64           `json:"blockVer"`

}

type CheckBlockResult struct {
    Correct      bool           `json:"correct"`
}

func NewCheckBlockResult() *CheckBlockResult {
    return &CheckBlockResult{}
}
func NewCheckBlockParams() *CheckBlockParams {
    return &CheckBlockParams{}
}
