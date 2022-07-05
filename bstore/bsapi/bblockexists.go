
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const BlockExistsMethod string = "blockExists"

type BlockExistsParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    BlockVer    int64           `json:"blockVer"`

}

type BlockExistsResult struct {
    Exists      bool           `json:"exists"`
}

func NewBlockExistsResult() *BlockExistsResult {
    return &BlockExistsResult{}
}
func NewBlockExistsParams() *BlockExistsParams {
    return &BlockExistsParams{}
}
