
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsapi

const GetHelloMethod string = "getHello"

type GetHelloParams struct {
    Message string      `json:"message,omitempty" msgpack:"message" `
}

type GetHelloResult struct {
    Message string      `json:"message,omitempty" msgpack:"message" `
}

func NewGetHelloResult() *GetHelloResult {
    return &GetHelloResult{}
}
func NewGetHelloParams() *GetHelloParams {
    return &GetHelloParams{}
}
