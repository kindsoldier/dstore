/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "fmt"
    "errors"
)

func Err2Err(message string, err error) error {
    return errors.New(fmt.Sprintf("%s: %v", message, err))
}

func NewErr(message string) error {
    return errors.New(message)
}
