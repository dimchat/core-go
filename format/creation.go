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
	. "github.com/dimchat/core-go/rfc"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

//
//  Plain Data
//

func NewPlainDataWithBytes(bytes []byte) TransportableData {
	return NewPlainData("", bytes)
}

func NewPlainDataWithString(text string) TransportableData {
	return NewPlainData(text, nil)
}

func ZeroPlainData() TransportableData {
	return NewPlainData("", []byte{})
}

//
//  Base-64
//

func NewBase64DataWithBytes(bytes []byte) TransportableData {
	return NewBase64Data("", bytes)
}

func NewBase64DataWithString(encoded string) TransportableData {
	return NewBase64Data(encoded, nil)
}

//
//  Data URI:
//
//      "data:image/jpg;base64,{BASE64_ENCODE}"
//      "data:audio/mp4;base64,{BASE64_ENCODE}"
//

func NewImageData(jpeg []byte) TransportableData {
	return NewEmbedDataWithType(MIMEType.IMAGE_JPG, jpeg)
}

func NewAudioData(mp4 []byte) TransportableData {
	return NewEmbedDataWithType(MIMEType.AUDIO_MP4, mp4)
}

func NewEmbedDataWithType(mimeType string, body []byte) TransportableData {
	head := NewDataHeader(mimeType, BASE_64, nil)
	return NewEmbedData("", body, nil, head)
}

func NewEmbedDataWithURI(uri DataURI) TransportableData {
	head := uri.Head()
	return NewEmbedData("", nil, uri, head)
}

//
//  PNF
//

func NewPortableNetworkFileWithMap(dict StringKeyMap) TransportableFile {
	return NewPortableNetworkFile(dict, nil, "", nil, nil)
}

func NewPortableNetworkFileWithData(data TransportableData, filename string,
	url URL, password DecryptKey,
) TransportableFile {
	return NewPortableNetworkFile(nil, data, filename, url, password)
}
