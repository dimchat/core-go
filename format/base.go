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
	"bytes"
	"fmt"

	. "github.com/dimchat/mkm-go/format"
)

//goland:noinspection GoSnakeCaseUsage
const (
	BASE_64 = "base64"
	BASE_58 = "base58"
	HEX     = "hex"
)

type EncodeData interface {
	TransportableData

	EncodedString() string
	DecodedBytes() []byte
}

type BaseData struct {
	//EncodeData

	encoded string // encoded string
	bytes   []byte // decoded bytes
}

func NewBaseData(encoded string, bytes []byte) *BaseData {
	return &BaseData{
		encoded: encoded,
		bytes:   bytes,
	}
}

// protected
func (ted *BaseData) EncodedString() string {
	return ted.encoded
}

// protected
func (ted *BaseData) DecodedBytes() []byte {
	return ted.bytes
}

// Override
func (ted *BaseData) IsEmpty() bool {
	// 1. check inner bytes
	bin := ted.bytes
	if len(bin) > 0 {
		return false
	}
	// 2. check inner string
	str := ted.encoded
	return str == ""
}

//
//  TransportableData
//

func size(ted TransportableData) int {
	//str := ted.String()
	//return len(str)
	bin := ted.Bytes()
	return len(bin)
}

//
//  TransportableResource
//

func serialize(ted TransportableData) any {
	return ted.String()
}

//
//  IObject
//

func equals(ted EncodeData, other any) bool {
	if other == nil {
		return ted.IsEmpty()
	} else if other == ted {
		// same object
		return true
	}
	// check value types
	switch v := other.(type) {
	case EncodeData:
		return dataEquals(ted, v)
	case TransportableData:
		return tedEquals(ted, v)
	case fmt.Stringer:
		// compare with encoded string
		return ted.String() == v.String()
	case string:
		// compare with encoded string
		return ted.String() == v
	default:
		//panic(fmt.Sprintf("unknown data: %v", other))
	}
	return false
}

func dataEquals(self, other EncodeData) bool {
	if other == nil || other.IsEmpty() {
		return self.IsEmpty()
	}
	// compare with inner string
	thisString := self.EncodedString()
	thatString := other.EncodedString()
	if thisString != "" && thatString != "" {
		return thisString == thatString
	}
	// compare with inner bytes
	thisBytes := self.DecodedBytes()
	thatBytes := other.DecodedBytes()
	if thisBytes != nil && thatBytes != nil {
		return bytes.Equal(thisBytes, thatBytes)
	}
	// compare with decoded bytes
	thisBytes = self.Bytes()
	thatBytes = other.Bytes()
	return bytes.Equal(thisBytes, thatBytes)
}

func tedEquals(self EncodeData, other TransportableData) bool {
	if other == nil || other.IsEmpty() {
		return self.IsEmpty()
	}
	// compare with encoded string
	thisString := self.EncodedString()
	if thisString != "" {
		thatString := other.String()
		return thisString == thatString
	}
	// compare with decoded bytes
	thisBytes := self.DecodedBytes()
	thatBytes := other.Bytes()
	return bytes.Equal(thisBytes, thatBytes)
}
