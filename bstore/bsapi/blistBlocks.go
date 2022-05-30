
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

import (
    "ndstore/dscom"
)

const ListBlocksMethod string = "listBlocks"

type ListBlocksParams struct {
}

type ListBlocksResult struct {
    Blocks  []dscom.BlockDescr     `json:"blocks,omitempty"`
}

func NewListBlocksResult() *ListBlocksResult {
    return &ListBlocksResult{}
}
func NewListBlocksParams() *ListBlocksParams {
    return &ListBlocksParams{}
}
