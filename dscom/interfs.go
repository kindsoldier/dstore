/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    "io"
)

type IBSPool interface {
    SaveBlock(fileId, batchId, blockId, blockSize int64, blockReader io.Reader, dataSize int64,
                                    blockType, hashAlg, hashInit, hashSum string) (int64, error)
}
