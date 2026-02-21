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
	. "github.com/dimchat/core-go/format"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Base File Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x10),
 *      "sn"   : 123,
 *
 *      "data"     : "...",        // base64_encode(fileContent)
 *      "filename" : "photo.png",
 *
 *      "URL"      : "http://...", // download from CDN
 *      // before fileContent uploaded to a public CDN,
 *      // it should be encrypted by a symmetric key
 *      "key"      : {             // symmetric key to decrypt file data
 *          "algorithm" : "AES",   // "DES", ...
 *          "data"      : "{BASE64_ENCODE}",
 *          ...
 *      }
 *  }
 *  </pre></blockquote>
 */
type BaseFileContent struct {
	//FileContent
	*BaseContent

	wrapper TransportableFileWrapper
}

func NewBaseFileContent(dict StringKeyMap,
	msgType MessageType,
	data TransportableData, filename string,
	url URL, password DecryptKey,
) *BaseFileContent {
	if dict != nil {
		// init file content with map
		return &BaseFileContent{
			BaseContent: NewBaseContent(dict, ""),
			wrapper:     CreateTransportableFileWrapper(dict, nil, "", nil, nil),
		}
	}
	// new file content
	if msgType == "" {
		msgType = ContentType.FILE
	}
	content := &BaseFileContent{
		BaseContent: NewBaseContent(nil, msgType),
	}
	content.wrapper = CreateTransportableFileWrapper(content.Map(), data, filename, url, password)
	return content
}

// Override
func (content *BaseFileContent) Map() StringKeyMap {
	// call wrapper to serialize 'data' & 'key'
	return content.wrapper.Map()
}

/**
 *  file data
 */

// Override
func (content *BaseFileContent) Data() TransportableData {
	return content.wrapper.Data()
}

// Override
func (content *BaseFileContent) SetData(data TransportableData) {
	content.wrapper.SetData(data)
}

/**
 *  file name
 */

// Override
func (content *BaseFileContent) Filename() string {
	return content.wrapper.Filename()
}

// Override
func (content *BaseFileContent) SetFilename(filename string) {
	content.wrapper.SetFilename(filename)
}

/**
 *  download URL
 */

// Override
func (content *BaseFileContent) URL() URL {
	return content.wrapper.URL()
}

// Override
func (content *BaseFileContent) SetURL(url URL) {
	content.wrapper.SetURL(url)
}

/**
 *  decrypt key
 */

// Override
func (content *BaseFileContent) Password() DecryptKey {
	return content.wrapper.Password()
}

// Override
func (content *BaseFileContent) SetPassword(password DecryptKey) {
	content.wrapper.SetPassword(password)
}
