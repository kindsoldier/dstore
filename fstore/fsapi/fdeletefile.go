
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const DeleteFileMethod string = "deleteFile"

type DeleteFileParams struct {
    FilePath    string      `json:"filePath"`
}

type DeleteFileResult struct {
}

func NewDeleteFileResult() *DeleteFileResult {
    return &DeleteFileResult{}
}
func NewDeleteFileParams() *DeleteFileParams {
    return &DeleteFileParams{}
}
