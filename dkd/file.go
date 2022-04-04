/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package dkd

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  File message: {
 *      type : 0x10,
 *      sn   : 123,
 *
 *      URL      : "http://", // upload to CDN
 *      data     : "...",     // if (!URL) base64_encode(fileContent)
 *      filename : "..."
 *  }
 */
type BaseFileContent struct {
	BaseContent

	_data []byte       // file data (plaintext)
	_key SymmetricKey  // symmetric key to decrypt the encrypted data from URL
}

func NewFileContent(msgType ContentType, filename string, data []byte) FileContent {
	content := new(BaseFileContent)
	content.InitWithType(msgType, filename, data)
	return content
}

/* designated initializer */
func (content *BaseFileContent) Init(dict map[string]interface{}) *BaseFileContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._data = nil
		content._key = nil
	}
	return content
}

/* designated initializer */
func (content *BaseFileContent) InitWithType(msgType ContentType, filename string, data []byte) *BaseFileContent {
	if msgType == 0 {
		msgType = FILE
	}
	if content.BaseContent.InitWithType(msgType) != nil {
		content.SetFilename(filename)
		content.SetData(data)
		content._key = nil
	}
	return content
}

//-------- IFileContent

func (content *BaseFileContent) URL() string {
	text := content.Get("URL")
	if text == nil {
		return ""
	}
	return text.(string)
}
func (content *BaseFileContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *BaseFileContent) Data() []byte {
	if content._data == nil {
		b64 := content.Get("data")
		if b64 != nil {
			content._data = Base64Decode(b64.(string))
		}
	}
	return content._data
}
func (content *BaseFileContent) SetData(data []byte) {
	if ValueIsNil(data) {
		content.Remove("data")
	} else {
		b64 := Base64Encode(data)
		content.Set("data", b64)
	}
	content._data = data
}

func (content *BaseFileContent) Filename() string {
	text := content.Get("filename")
	if text == nil {
		return ""
	}
	return text.(string)
}
func (content *BaseFileContent) SetFilename(filename string) {
	content.Set("filename", filename)
}

func (content *BaseFileContent) Password() SymmetricKey {
	if content._key == nil {
		dict := content.Get("password")
		content._key = SymmetricKeyParse(dict)
	}
	return content._key
}

func (content *BaseFileContent) SetPassword(password SymmetricKey) {
	if ValueIsNil(password) {
		content.Remove("password")
	} else {
		content.Set("password", password.GetMap(false))
	}
	content._key = password
}
