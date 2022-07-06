
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const LinkBlockMethod string = "linkBlock"

type LinkBlockParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    OldBlockVer    int64        `json:"oldBlockVer"`
    NewBlockVer    int64        `json:"newBlockVer"`
}

type LinkBlockResult struct {
}

func NewLinkBlockResult() *LinkBlockResult {
    return &LinkBlockResult{}
}
func NewLinkBlockParams() *LinkBlockParams {
    return &LinkBlockParams{}
}
