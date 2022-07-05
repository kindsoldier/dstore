
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const EraseAllMethod string = "eraseAll"

type EraseAllParams struct {
}

type EraseAllResult struct {
}

func NewEraseAllResult() *EraseAllResult {
    return &EraseAllResult{}
}
func NewEraseAllParams() *EraseAllParams {
    return &EraseAllParams{}
}
