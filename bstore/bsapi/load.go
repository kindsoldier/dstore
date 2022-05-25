
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const LoadMethod string = "load"

type LoadParams struct {
    FileId      int64           `json:"fileId"`
    BatchId     int64           `json:"batchId"`
    BlockId     int64           `json:"blockId"`
}

type LoadResult struct {
}

func NewLoadResult() *LoadResult {
    return &LoadResult{}
}
func NewLoadParams() *LoadParams {
    return &LoadParams{}
}
