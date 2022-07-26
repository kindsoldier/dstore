
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

import (
    "dstore/dscomm/dsdescr"
)

const SaveBlockMethod string = "saveBlock"

type SaveBlockParams struct {
    FileId      int64           `msgpack:"fileId"       json:"fileId"`
    BatchId     int64           `msgpack:"batchId"      json:"batchId"`
    BlockType   int64           `msgpack:"blockType"    json:"blockType"`
    BlockId     int64           `msgpack:"blockId"      json:"blockId"`
    BlockSize   int64           `msgpack:"blockSize"    json:"blockSize"`
}

type SaveBlockResult struct {
}

func NewSaveBlockResult() *SaveBlockResult {
    return &SaveBlockResult{}
}

func NewSaveBlockParams() *SaveBlockParams {
    return &SaveBlockParams{}
}

const LoadBlockMethod string = "loadBlock"

type LoadBlockParams struct {
    FileId      int64           `msgpack:"fileId"       json:"fileId"`
    BatchId     int64           `msgpack:"batchId"      json:"batchId"`
    BlockType   int64           `msgpack:"blockType"    json:"blockType"`
    BlockId     int64           `msgpack:"blockId"      json:"blockId"`
}

type LoadBlockResult struct {
}

func NewLoadBlockResult() *LoadBlockResult {
    return &LoadBlockResult{}
}
func NewLoadBlockParams() *LoadBlockParams {
    return &LoadBlockParams{}
}




const ListBlocksMethod string = "listBlocks"

type ListBlocksParams struct {
    FileId      int64           `msgpack:"fileId"       json:"fileId"`
}

type ListBlocksResult struct {
    Blocks   []*dsdescr.Block   `msgpack:"blocks"       json:"blocks"`
}

func NewListBlocksResult() *ListBlocksResult {
    return &ListBlocksResult{}
}

func NewListBlocksParams() *ListBlocksParams {
    return &ListBlocksParams{}
}

const DeleteBlockMethod string = "deleteBlock"
type DeleteBlockParams struct {
    FileId      int64           `msgpack:"fileId"       json:"fileId"`
    BatchId     int64           `msgpack:"batchId"      json:"batchId"`
    BlockType   int64           `msgpack:"blockType"    json:"blockType"`
    BlockId     int64           `msgpack:"blockId"      json:"blockId"`
}

type DeleteBlockResult struct {
}

func NewDeleteBlockResult() *DeleteBlockResult {
    return &DeleteBlockResult{}
}

func NewDeleteBlockParams() *DeleteBlockParams {
    return &DeleteBlockParams{}
}
