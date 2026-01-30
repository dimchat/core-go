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

import . "github.com/dimchat/mkm-go/format"

/**
 *  UTF-8 encoding
 */
type PlainData struct {
	BaseData
}

func (ted *PlainData) InitWithString(str string) {
	ted.BaseData.InitWithString(str)
}

func (ted *PlainData) InitWithBytes(bin []byte) {
	ted.BaseData.InitWithBytes(bin)
}

//
//  TransportableData
//

// Override
func (ted *PlainData) Encoding() string {
	return ""
}

// Override
func (ted *PlainData) Bytes() []byte {
	bin := ted._binary
	if bin == nil {
		bin = UTF8Encode(ted._string)
		ted._binary = bin
	}
	return bin
}

// Override
func (ted *PlainData) String() string {
	base64 := ted._string
	if base64 == "" {
		base64 = UTF8Decode(ted._binary)
		ted._string = base64
	}
	return base64
}

// Override
func (ted *PlainData) Size() int {
	return size(ted)
}

//
//  TransportableResource
//

// Override
func (ted *PlainData) Serialize() interface{} {
	return serialize(ted)
}

//
//  IObject
//

func (ted *PlainData) Equal(other interface{}) bool {
	return equals(ted, other)
}

//
//  Factory methods
//

func CreatePlainDataWithBytes(bin []byte) TransportableData {
	return &PlainData{
		BaseData: BaseData{
			_string: "",
			_binary: bin,
		},
	}
}

func CreatePlainDataWithString(str string) TransportableData {
	return &PlainData{
		BaseData: BaseData{
			_string: str,
			_binary: nil,
		},
	}
}

func ZeroPlainData() TransportableData {
	return &PlainData{
		BaseData: BaseData{
			_binary: []byte{},
			_string: "",
		},
	}
}
