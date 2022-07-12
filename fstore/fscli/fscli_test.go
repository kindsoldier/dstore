package main

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
    "dstore/dsrpc"
)

func BenchmarkStatus(b *testing.B) {
    util := NewUtil()
    util.URI = fmt.Sprintf("%s:%s", util.Address, util.Port)

    auth := dsrpc.CreateAuth([]byte(util.aLogin), []byte(util.aPass))
    b.ResetTimer()

    pBench := func(pb *testing.PB) {
        for pb.Next() {
            _, err := util.GetStatusCmd(auth)
            require.NoError(b, err)
        }
    }
    b.SetParallelism(100)
    b.RunParallel(pBench)
}
