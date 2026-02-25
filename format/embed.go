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
	//TransportableData
	*BaseData

	dataURI  DataURI
	dataHead DataHeader
}

func NewEmbedData(encoded string, bytes []byte, uri DataURI, head DataHeader) *EmbedData {
	return &EmbedData{
		BaseData: NewBaseData(encoded, bytes),
		dataURI:  uri,
		dataHead: head,
	}
}

// protected
func (ted *EmbedData) DataCoder() DataCoder {
	encoding := ted.Encoding()
	return GetDataCoder(encoding)
}

// protected
func (ted *EmbedData) encodeDataURI() DataURI {
	uri := ted.dataURI
	if uri == nil {
		coder := ted.DataCoder()
		data := ted.bytes
		if coder == nil || len(data) == 0 {
			// cannot encode data
			return nil
		}
		base64 := coder.Encode(data)
		uri = NewDataURI(ted.dataHead, base64)
		ted.dataURI = uri
	}
	return uri
}

//
//  TransportableData
//

// Override
func (ted *EmbedData) Encoding() string {
	return ted.dataHead.Encoding() // "base64"
}

// Override
func (ted *EmbedData) Bytes() []byte {
	bin := ted.bytes
	if bin == nil {
		uri := ted.dataURI
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
		ted.bytes = bin
	}
	return bin
}

// Override
func (ted *EmbedData) String() string {
	base64 := ted.encoded
	if base64 == "" {
		uri := ted.encodeDataURI()
		if uri != nil {
			base64 = uri.String()
			ted.encoded = base64
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
func (ted *EmbedData) Serialize() any {
	return serialize(ted)
}

//
//  IObject
//

func (ted *EmbedData) Equal(other any) bool {
	return equals(ted, other)
}
