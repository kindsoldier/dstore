
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const DeleteMethod string = "delete"

type DeleteParams struct {
    FilePath    string      `json:"filePath"`
}

type DeleteResult struct {
}

func NewDeleteResult() *DeleteResult {
    return &DeleteResult{}
}
func NewDeleteParams() *DeleteParams {
    return &DeleteParams{}
}
