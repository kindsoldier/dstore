
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const LoadMethod string = "load"

type LoadParams struct {
    ClusterId   int64           `json:"clusterId"`
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
