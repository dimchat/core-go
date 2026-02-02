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

type PortableNetworkFile struct {
	//TransportableFile
	Dictionary

	_wrapper TransportableFileWrapper
}

func (pnf *PortableNetworkFile) createWrapper() TransportableFileWrapper {
	dict := pnf.Map()
	factory := GetTransportableFileWrapperFactory()
	return factory.CreateTransportableFileWrapper(dict)
}

func (pnf *PortableNetworkFile) InitWithMap(dict StringKeyMap) TransportableFile {
	if pnf.Dictionary.InitWithMap(dict) != nil {
		pnf._wrapper = pnf.createWrapper()
	}
	return pnf
}

func (pnf *PortableNetworkFile) Init(data TransportableData, filename string, url URL, key DecryptKey) TransportableFile {
	if pnf.Dictionary.Init() != nil {
		wrapper := pnf.createWrapper()
		pnf._wrapper = wrapper
		// file data
		if data != nil {
			wrapper.SetData(data)
		}
		// file name
		if filename != "" {
			wrapper.SetFilename(filename)
		}
		// download URL
		if url != nil {
			wrapper.SetURL(url)
		}
		// decrypt key
		if key != nil {
			wrapper.SetPassword(key)
		}
	}
	return pnf
}

func (pnf *PortableNetworkFile) getURIString() string {
	// serialize
	dict := pnf._wrapper.Map()
	// check 'URL'
	url := pnf.URL()
	if url != nil {
		count := len(dict)
		if count == 1 {
			// this PNF info contains 'URL' only,
			// so return the URI string here.
			return url.String()
		} else if count == 2 {
			// check filename
			_, exists := dict["filename"]
			if exists {
				// ignore 'filename' field
				return url.String()
			}
		}
		// this PNF info contains other params,
		// cannot serialize it as a string.
		return ""
	}
	// check 'data'
	text := pnf.GetString("data", "")
	uri := ParseDataURI(text)
	if uri != nil {
		count := len(dict)
		if count == 1 {
			// this PNF info contains 'data' only,
			// and it is a data URI,
			// so return the URI string here.
			return text
		} else if count == 2 {
			// check filename
			value, exists := dict["filename"]
			if exists {
				if value == nil {
					// nothing changed
					return text
				}
				filename := ConvertString(value, "")
				// add 'filename' to data URI
				uri = newDataURI(uri, filename)
				return uri.String()
			}
		}
		// this PNF info contains other params,
		// cannot serialize it as a string.
		return ""
	}
	// the file data was saved into local storage,
	// so there is just a 'filename' here,
	// cannot build URI string
	return ""
}

func newDataURI(uri DataURI, filename string) DataURI {
	head := uri.Head()
	extra := NewMap()
	// copy extra values
	keys := head.ExtraKeys()
	if keys != nil {
		// copy all entries
		for _, name := range keys {
			extra[name] = head.ExtraValue(name)
		}
	}
	if filename != "" {
		// update 'filename'
		extra["filename"] = filename
	} else {
		// erase 'filename'
		delete(extra, "filename")
	}
	// create data URI with new extra info
	head = NewDataHeader(head.MimeType(), head.Encoding(), extra)
	return NewDataURI(head, uri.Body())
}

// Override
func (pnf *PortableNetworkFile) String() string {
	uri := pnf.getURIString()
	if uri != "" {
		// this PNF can be simplified to a URI string
		return uri
	}
	// return JsON string
	dict := pnf._wrapper.Map()
	return JSONEncodeMap(dict)
}

// Override
func (pnf *PortableNetworkFile) Map() StringKeyMap {
	// call wrapper to serialize 'data' & 'key'
	return pnf._wrapper.Map()
}

// Override
func (pnf *PortableNetworkFile) Serialize() interface{} {
	uri := pnf.getURIString()
	if uri != "" {
		// this PNF can be simplified to a URI string
		return uri
	}
	// return inner map
	return pnf._wrapper.Map()
}

/**
 *  file data
 */

// Override
func (pnf *PortableNetworkFile) Data() TransportableData {
	return pnf._wrapper.Data()
}

// Override
func (pnf *PortableNetworkFile) SetData(data TransportableData) {
	pnf._wrapper.SetData(data)
}

/**
 *  file name
 */

// Override
func (pnf *PortableNetworkFile) Filename() string {
	return pnf._wrapper.Filename()
}

// Override
func (pnf *PortableNetworkFile) SetFilename(filename string) {
	pnf._wrapper.SetFilename(filename)
}

/**
 *  download URL
 */

// Override
func (pnf *PortableNetworkFile) URL() URL {
	return pnf._wrapper.URL()
}

// Override
func (pnf *PortableNetworkFile) SetURL(url URL) {
	pnf._wrapper.SetURL(url)
}

/**
 *  decrypt key
 */

// Override
func (pnf *PortableNetworkFile) Password() DecryptKey {
	return pnf._wrapper.Password()
}

// Override
func (pnf *PortableNetworkFile) SetPassword(password DecryptKey) {
	pnf._wrapper.SetPassword(password)
}
