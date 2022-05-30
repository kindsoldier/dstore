
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsapi

const GetHelloMethod string = "getHello"

type GetHelloParams struct {
    Message string      `json:"message" msgpack:"message" `
}

type GetHelloResult struct {
    Message string      `json:"message" msgpack:"message" `
}

func NewGetHelloResult() *GetHelloResult {
    return &GetHelloResult{}
}
func NewGetHelloParams() *GetHelloParams {
    return &GetHelloParams{}
}
