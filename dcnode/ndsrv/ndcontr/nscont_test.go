package ndcontr

import (
    "testing"
    "dcstore/dcnode/ndapi"
    "dcstore/dcrpc"

    "github.com/stretchr/testify/assert"
)


func TestHello(t *testing.T) {
    var err error
    helloResp := "hello!"

    params := ndapi.NewHelloParams()
    params.Message = "hello server!"

    result := ndapi.NewHelloResult()

    contr := NewContr()

    err = dcrpc.LocalExec(ndapi.HelloMethod, params, result, nil, contr.HelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)
}
