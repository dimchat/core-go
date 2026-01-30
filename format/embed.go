/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
package format

import (
	"fmt"

	. "github.com/dimchat/core-go/rfc"
	. "github.com/dimchat/mkm-go/format"
)

type EmbedData struct {
	BaseData

	_dataURI  DataURI
	_dataHead DataHeader
}

func (ted *EmbedData) InitWithURI(uri DataURI) {
	ted.BaseData.InitWithString(uri.String())
	ted._dataURI = uri
	ted._dataHead = uri.Head()
}

func (ted *EmbedData) InitWithBytes(body []byte, head DataHeader) {
	ted.BaseData.InitWithBytes(body)
	ted._dataURI = nil
	ted._dataHead = head
}

// protected
func (ted *EmbedData) DataCoder() DataCoder {
	encoding := ted.Encoding()
	return GetDataCoder(encoding)
}

// protected
func (ted *EmbedData) encodeDataURI() DataURI {
	uri := ted._dataURI
	if uri == nil {
		coder := ted.DataCoder()
		data := ted._binary
		if coder == nil || data == nil || len(data) == 0 {
			// cannot encode data
			return nil
		}
		base64 := coder.Encode(data)
		uri = NewDataURI(ted._dataHead, base64)
		ted._dataURI = uri
	}
	return uri
}

//
//  TransportableData
//

// Override
func (ted *EmbedData) Encoding() string {
	return ted._dataHead.Encoding() // "base64"
}

// Override
func (ted *EmbedData) Bytes() []byte {
	bin := ted._binary
	if bin == nil {
		uri := ted._dataURI
		if uri == nil || uri.IsEmpty() {
			panic(fmt.Sprintf("data URI error: %v", uri))
			return nil
		}
		coder := ted.DataCoder()
		base64 := uri.Body()
		if coder == nil || base64 == "" {
			//panic(fmt.Sprintf("cannot decode data: %s", ted.Encoding()))
			return nil
		}
		bin = coder.Decode(base64)
		ted._binary = bin
	}
	return bin
}

// Override
func (ted *EmbedData) String() string {
	base64 := ted._string
	if base64 == "" {
		uri := ted.encodeDataURI()
		if uri != nil {
			base64 = uri.String()
			ted._string = base64
		}
	}
	return base64
}

// Override
func (ted *EmbedData) Size() int {
	return size(ted)
}

//
//  TransportableResource
//

// Override
func (ted *EmbedData) Serialize() interface{} {
	return serialize(ted)
}

//
//  IObject
//

func (ted *EmbedData) Equal(other interface{}) bool {
	return equals(ted, other)
}

//
//  Factories:
//
//      "data:image/jpg;base64,{BASE64_ENCODE}"
//      "data:audio/mp4;base64,{BASE64_ENCODE}"
//

func CreateImageData(jpeg []byte) TransportableData {
	return CreateEmbedData(MIMEType.IMAGE_JPG, jpeg)
}

func CreateAudioData(mp4 []byte) TransportableData {
	return CreateEmbedData(MIMEType.AUDIO_MP4, mp4)
}

func CreateEmbedData(mimeType string, body []byte) TransportableData {
	head := NewDataHeader(mimeType, BASE_64, nil)
	return NewEmbedData(head, body)
}

func NewEmbedData(head DataHeader, body []byte) TransportableData {
	return &EmbedData{
		BaseData: BaseData{
			_string: "",
			_binary: body,
		},
		_dataURI:  nil,
		_dataHead: head,
	}
}
