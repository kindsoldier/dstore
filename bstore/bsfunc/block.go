/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */


package bsfunc

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/dscom"
)

const HelloMessage string = "hello"

func GetHello(uri string, auth *dsrpc.Auth) error {
    var err error

    params := bsapi.NewGetHelloParams()
    params.Message = HelloMessage
    result := bsapi.NewGetHelloResult()

    err = dsrpc.Exec(uri, bsapi.GetHelloMethod, params, result, auth)
    if err != nil {
        return err
    }
    return err
}

func SaveBlock(uri string, auth *dsrpc.Auth, fileId, batchId, blockId, blockSize,
                                dataSize int64, blockReader io.Reader, binSize int64, blockType,
                                hashAlg, hashInit, hashSum string) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    params.FileId   = fileId
    params.BatchId  = batchId
    params.BlockId  = blockId
    params.DataSize = binSize

    params.BlockType = blockType
    params.HashAlg  = hashAlg
    params.HashInit = hashInit
    params.HashSum  = hashSum
    result := bsapi.NewSaveBlockResult()

    err = dsrpc.Put(uri, bsapi.SaveBlockMethod, blockReader, binSize, params, result, auth)
    if err != nil {
        return err
    }
    return err
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
        return err
    }
    return err
}

func ListBlocks(uri string, auth *dsrpc.Auth) ([]*dscom.BlockDescr, error) {
    var err error
    blockDescrs := make([]*dscom.BlockDescr, 0)
    params := bsapi.NewListBlocksParams()
    result := bsapi.NewListBlocksResult()
    err = dsrpc.Exec(uri, bsapi.ListBlocksMethod, params, result, auth)
    if err != nil {
        return blockDescrs, err
    }
    blockDescrs = result.Blocks
    return blockDescrs, err
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
        return err
    }
    return err
}
