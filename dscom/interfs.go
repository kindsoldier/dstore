/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package dscom

import (
    //"io"
    "os"
)

type IBSPool interface {
    SaveBlock(fileId, batchId, blockId, blockSize int64, blockReader *os.File, dataSize int64,
                                    blockType, hashAlg, hashInit, hashSum string) (int64, error)
}
