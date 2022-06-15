/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package bsfunc

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/dscom"
    "ndstore/dserr"
)

const HelloMessage string = "hello"

func GetHello(uri string, auth *dsrpc.Auth) error {
    var err error

    params := bsapi.NewGetHelloParams()
    params.Message = HelloMessage
    result := bsapi.NewGetHelloResult()

    err = dsrpc.Exec(uri, bsapi.GetHelloMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func SaveBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId, blockSize int64,
                                            blockReader io.Reader, binSize int64, blockType,
                                                    hashAlg, hashInit, hashSum string) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.BlockSize = blockSize
    params.DataSize = binSize

    params.BlockType = blockType
    params.HashAlg  = hashAlg
    params.HashInit = hashInit
    params.HashSum  = hashSum
    result := bsapi.NewSaveBlockResult()

    err = dsrpc.Put(uri, bsapi.SaveBlockMethod, blockReader, binSize, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func LoadBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64,
                                                blockWriter io.Writer, blockType string) error {
    var err error
    params := bsapi.NewLoadBlockParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.BlockType = blockType

    result := bsapi.NewLoadBlockResult()
    err = dsrpc.Get(uri, bsapi.LoadBlockMethod, blockWriter, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func ListBlocks(uri string, auth *dsrpc.Auth) ([]*dscom.BlockDescr, error) {
    var err error
    blockDescrs := make([]*dscom.BlockDescr, 0)
    params := bsapi.NewListBlocksParams()
    result := bsapi.NewListBlocksResult()
    err = dsrpc.Exec(uri, bsapi.ListBlocksMethod, params, result, auth)
    if err != nil {
        return blockDescrs, dserr.Err(err)
    }
    blockDescrs = result.Blocks
    return blockDescrs, dserr.Err(err)
}

func DeleteBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string) error {
    var err error
    params := bsapi.NewDeleteBlockParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.BlockType = blockType
    result := bsapi.NewDeleteBlockResult()
    err = dsrpc.Exec(uri, bsapi.DeleteBlockMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func BlockExists(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string) (bool, error) {
    var err error
    var exists bool
    params := bsapi.NewBlockExistsParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.BlockType = blockType
    result := bsapi.NewBlockExistsResult()
    err = dsrpc.Exec(uri, bsapi.BlockExistsMethod, params, result, auth)
    exists = result.Exists
    if err != nil {
        return exists, dserr.Err(err)
    }
    return exists, dserr.Err(err)
}

func CheckBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId int64, blockType string) (bool, error) {
    var err error
    var correct bool
    params := bsapi.NewCheckBlockParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.BlockType = blockType
    result := bsapi.NewCheckBlockResult()
    err = dsrpc.Exec(uri, bsapi.CheckBlockMethod, params, result, auth)
    correct = result.Correct
    if err != nil {
        return correct, dserr.Err(err)
    }
    return correct, dserr.Err(err)
}

func PurgeAll(uri string, auth *dsrpc.Auth) error {
    var err error
    params := bsapi.NewPurgeAllParams()
    result := bsapi.NewPurgeAllResult()
    err = dsrpc.Exec(uri, bsapi.PurgeAllMethod, params, result, auth)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}
