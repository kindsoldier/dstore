/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dcrpc

type Empty struct {}

func NewEmpty() *Empty {
    return &Empty{}
}
