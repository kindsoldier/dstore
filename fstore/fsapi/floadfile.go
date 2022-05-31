
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const LoadFileMethod string = "loadFile"

type LoadFileParams struct {
    FilePath    string      `json:"filePath"`
}

type LoadFileResult struct {
}

func NewLoadFileResult() *LoadFileResult {
    return &LoadFileResult{}
}
func NewLoadFileParams() *LoadFileParams {
    return &LoadFileParams{}
}
