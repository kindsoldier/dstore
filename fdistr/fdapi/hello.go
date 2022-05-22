
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdapi

const HelloMethod string = "hello"

type HelloParams struct {
    Message string      `json:"message" msgpack:"message" `
}

type HelloResult struct {
    Message string      `json:"message" msgpack:"message" `
}

func NewHelloResult() *HelloResult {
    return &HelloResult{}
}
func NewHelloParams() *HelloParams {
    return &HelloParams{}
}

