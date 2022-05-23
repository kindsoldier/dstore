
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const SaveMethod string = "save"

type SaveParams struct {
    FilePath    string      `json:"filePath"`
}

type SaveResult struct {
}

func NewSaveResult() *SaveResult {
    return &SaveResult{}
}
func NewSaveParams() *SaveParams {
    return &SaveParams{}
}
