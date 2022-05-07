
/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package ndapi

import (
    //"encoding/json"
)


const HelloMethod string = "hello"

type HelloParams struct {
    Message string      `json:"message" msgpack:"message" `
}

func NewHelloParams() *HelloParams {
    return &HelloParams{}
}

//func (this *HelloParams) JSON() []byte {
//    jBytes, _ := json.Marshal(this)
//    return jBytes
//}


type HelloResult struct {
    Message string      `json:"message" msgpack:"message" `
}

func NewHelloResult() *HelloResult {
    return &HelloResult{}
}

//func (this *HelloResult) JSON() []byte {
//    jBytes, _ := json.Marshal(this)
//    return jBytes
//}
