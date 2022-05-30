
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const SaveFileMethod string = "saveFile"

type SaveFileParams struct {
    FilePath    string      `json:"filePath"`
}

type SaveFileResult struct {
}

func NewSaveFileResult() *SaveFileResult {
    return &SaveFileResult{}
}
func NewSaveFileParams() *SaveFileParams {
    return &SaveFileParams{}
}
