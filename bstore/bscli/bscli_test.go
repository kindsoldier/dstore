package main

import (
    "fmt"
    //"path/filepath"
    "testing"
    //"math/rand"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)

func BenchmarkHello(b *testing.B) {
    //var err error

    util := NewUtil()
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)

    auth := dsrpc.CreateAuth([]byte(util.ALogin), []byte(util.APass))
    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            _, err := util.GetHelloCmd(auth)
            assert.NoError(b, err)
        }
    }
    b.SetParallelism(10)
    b.RunParallel(pBench)
}
