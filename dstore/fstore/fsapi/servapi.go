
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const GetStatusMethod string = "getStatus"

type GetStatusParams struct {
}

type GetStatusResult struct {
    Uptime int64      `json:"uptime"    msgpack:"uptime"`
}

func NewGetStatusResult() *GetStatusResult {
    return &GetStatusResult{}
}
func NewGetStatusParams() *GetStatusParams {
    return &GetStatusParams{}
}
