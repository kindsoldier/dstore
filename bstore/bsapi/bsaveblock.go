
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const SaveBlockMethod string = "saveBlock"

type SaveBlockParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
    BlockType   string          `json:"blockType"`
    BlockVer    int64           `json:"blockVer"`

    BlockSize   int64           `json:"blockSize"`
    DataSize    int64           `json:"dataSize"`
    HashAlg     string          `json:"hashAlg"`
    HashSum     string          `json:"hashSum"`
    HashInit    string          `json:"hashInit"`
}

type SaveBlockResult struct {
}

func NewSaveBlockResult() *SaveBlockResult {
    return &SaveBlockResult{}
}
func NewSaveBlockParams() *SaveBlockParams {
    return &SaveBlockParams{}
}
