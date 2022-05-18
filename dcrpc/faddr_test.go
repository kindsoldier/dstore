/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dcrpc

import (
    "net"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFConn0(t *testing.T) {
    var cConn, sConn net.Conn
    sConn, cConn = NewFConn()

    cData := []byte("qwerty")
    count := 10

    for i := 0; i < count; i++ {
        wc, err := cConn.Write(cData)
        if err != nil {
            t.Error(err)
        }
        assert.Equal(t, wc, len(cData))

        sData := make([]byte, len(cData))
        rc, err := sConn.Read(sData)
        assert.NoError(t, err)
        assert.Equal(t, rc, len(cData))
        assert.Equal(t, cData, sData)
    }
}

func TestFConn1(t *testing.T) {
    var cConn, sConn net.Conn
    cConn, sConn = NewFConn()

    cData := []byte("qwerty")
    count := 10

    for i := 0; i < count; i++ {
        wc, err := cConn.Write(cData)
        if err != nil {
            t.Error(err)
        }
        assert.Equal(t, wc, len(cData))

        sData := make([]byte, len(cData))
        rc, err := sConn.Read(sData)
        assert.NoError(t, err)
        assert.Equal(t, rc, len(cData))
        assert.Equal(t, cData, sData)
    }
}
