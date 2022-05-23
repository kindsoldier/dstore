
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdapi

const LoadMethod string = "load"

type LoadParams struct {
    FilePath    string      `json:"filePath"`
}

type LoadResult struct {
}

func NewLoadResult() *LoadResult {
    return &LoadResult{}
}
func NewLoadParams() *LoadParams {
    return &LoadParams{}
}
