/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package xtools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandHexCreate(t *testing.T) {
	for size := 0; size < 1024*16; size++ {
		randHex := RandBytesHex(size)
		assert.Equal(t, len(randHex), size)
	}
}

func TestRandIntNeg(t *testing.T) {
	min := -1024
	max := -128
	t.Log("min-max range", max-min)
	theRand, err := RandInt(min, max)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("rand", theRand)
}

func TestRandIntRange0(t *testing.T) {
	min := 12
	max := 12
	t.Log("min-max range", max-min)
	theRand, err := RandInt(min, max)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("rand", theRand)
}
