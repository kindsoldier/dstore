/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "io"
)
type IFileSender interface {
}

type IFileReg interface {
}

type IBatch interface {
    Read(writer io.Writer) (int64, error)
    Write(reader io.Reader, need int64) (int64, error)
    Clean() error
    Close() error
}

type IBlock interface {
    Read(writer io.Writer) (int64, error)
    Write(reader io.Reader, need int64) (int64, error)
    Clean() error
    Close() error
}
