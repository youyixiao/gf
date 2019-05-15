// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gaes provides useful API for AES encryption/decryption algorithms.
package gaes

import (
    "bytes"
    "errors"
    "crypto/aes"
    "crypto/cipher"
)

const (
    ivDefValue = "I Love Go Frame!"
)

// AES加密, 使用CBC模式，注意key必须为16/24/32位长度，iv初始化向量为非必需参数
func Encrypt(plainText []byte, key []byte, iv...[]byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    plainText = PKCS5Padding(plainText, blockSize)
    ivValue := ([]byte)(nil)
    if len(iv) > 0 {
        ivValue = iv[0]
    } else {
        ivValue = []byte(ivDefValue)
    }
    blockMode  := cipher.NewCBCEncrypter(block, ivValue)
    ciphertext := make([]byte, len(plainText))
    blockMode.CryptBlocks(ciphertext, plainText)

    return ciphertext, nil
}

// AES解密, 使用CBC模式，注意key必须为16/24/32位长度，iv初始化向量为非必需参数
func Decrypt(cipherText []byte, key []byte, iv...[]byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    if len(cipherText) < blockSize {
        return nil, errors.New("cipherText too short")
    }
    ivValue := ([]byte)(nil)
    if len(iv) > 0 {
        ivValue = iv[0]
    } else {
        ivValue = []byte(ivDefValue)
    }
    if len(cipherText)%blockSize != 0 {
        return nil, errors.New("cipherText is not a multiple of the block size")
    }
    blockModel := cipher.NewCBCDecrypter(block, ivValue)
    plainText  := make([]byte, len(cipherText))
    blockModel.CryptBlocks(plainText, cipherText)
    plainText, e := PKCS5UnPadding(plainText, blockSize)
	if e != nil {
		return nil, e
	}

    return plainText, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
    padding := blockSize - len(src)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}

func PKCS5UnPadding(src []byte, blockSize int) ([]byte, error) {
    length    := len(src)
    if blockSize <= 0 {
        return nil, errors.New("invalid blocklen")
    }

	if length%blockSize != 0 || length == 0 {
		return nil, errors.New("invalid data len")
	}

    unpadding := int(src[length - 1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, errors.New("invalid padding")
	}

    padding := src[length - unpadding:]
	for i := 0; i < unpadding; i++ {
		if padding[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}

    return src[:(length - unpadding)], nil
}
