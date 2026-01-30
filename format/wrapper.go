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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  PNF Wrapper
 */
type TransportableFileWrapper interface {

	// serialize data
	Map() StringKeyMap

	/**
	 *  file data
	 */
	Data() TransportableData
	SetData(data TransportableData)

	/**
	 *  file name
	 */
	Filename() string
	SetFilename(filename string)

	/**
	 *  download URL
	 */
	URL() URL
	SetURL(url URL)

	/**
	 *  decrypt key
	 */
	Password() DecryptKey
	SetPassword(password DecryptKey)
}

/**
 *  Wrapper Factory
 */
type TransportableFileWrapperFactory interface {
	CreateTransportableFileWrapper(content StringKeyMap) TransportableFileWrapper
}

var sharedTransportableFileWrapperFactory TransportableFileWrapperFactory = &PortableNetworkFileWrapperFactory{}

func SetTransportableFileWrapperFactory(factory TransportableFileWrapperFactory) {
	sharedTransportableFileWrapperFactory = factory
}

func GetTransportableFileWrapperFactory() TransportableFileWrapperFactory {
	return sharedTransportableFileWrapperFactory
}

type PortableNetworkFileWrapperFactory struct {
	//TransportableFileWrapperFactory
}

func (factory *PortableNetworkFileWrapperFactory) CreateTransportableFileWrapper(content StringKeyMap) TransportableFileWrapper {
	wrapper := &PortableNetworkFileWrapper{}
	wrapper.Init(content)
	return wrapper
}

/**
 *  File Content MixIn
 *
 *  <blockquote><pre>
 *  {
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
type PortableNetworkFileWrapper struct {
	//TransportableFileWrapper

	_dictionary StringKeyMap

	// file content (not encrypted)
	_attachment TransportableData

	// download from CDN
	_remoteURL URL
	// key to decrypt data downloaded from CDN
	_password DecryptKey
}

func (wrapper *PortableNetworkFileWrapper) Init(dictionary StringKeyMap) TransportableFileWrapper {
	if ValueIsNil(dictionary) {
		// create empty map
		dictionary = NewMap()
	}
	wrapper._dictionary = dictionary
	// lazy load
	wrapper._attachment = nil
	wrapper._remoteURL = nil
	wrapper._password = nil
	return wrapper
}

func (wrapper *PortableNetworkFileWrapper) Get(key string) interface{} {
	value, exists := wrapper._dictionary[key]
	if !exists {
		return nil
	}
	return value
}

func (wrapper *PortableNetworkFileWrapper) Set(key string, value interface{}) {
	if ValueIsNil(value) {
		delete(wrapper._dictionary, key)
	} else {
		wrapper._dictionary[key] = value
	}
}

func (wrapper *PortableNetworkFileWrapper) Remove(key string) {
	delete(wrapper._dictionary, key)
}

func (wrapper *PortableNetworkFileWrapper) GetString(key string) string {
	value := wrapper.Get(key)
	return ConvertString(value, "")
}

//-------- TransportableFileWrapper

// Override
func (wrapper *PortableNetworkFileWrapper) Map() StringKeyMap {
	// serialize 'data'
	ted := wrapper._attachment
	if ted != nil && wrapper.Get("data") == nil {
		wrapper.Set("data", ted.Serialize())
	}
	// serialize 'key'
	pwd := wrapper._password
	if pwd != nil && wrapper.Get("key") == nil {
		wrapper.Set("key", pwd.Map())
	}
	// OK
	return wrapper._dictionary
}

// Override
func (wrapper *PortableNetworkFileWrapper) Data() TransportableData {
	ted := wrapper._attachment
	if ted == nil {
		base64 := wrapper.Get("data")
		ted = ParseTransportableData(base64)
		wrapper._attachment = ted
	}
	return ted
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetData(ted TransportableData) {
	wrapper.Remove("data")
	//if ted != nil {
	//	wrapper.Set("data", ted.Serialize())
	//}
	wrapper._attachment = ted
}

// Override
func (wrapper *PortableNetworkFileWrapper) Filename() string {
	return wrapper.GetString("filename")
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetFilename(filename string) {
	if filename == "" {
		wrapper.Remove("filename")
	} else {
		wrapper.Set("filename", filename)
	}
}

// Override
func (wrapper *PortableNetworkFileWrapper) URL() URL {
	remote := wrapper._remoteURL
	if remote == nil {
		locator := wrapper.GetString("URL")
		if locator != "" {
			remote = ParseURL(locator)
			wrapper._remoteURL = remote
		}
	}
	return remote
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetURL(remote URL) {
	if remote == nil {
		wrapper.Remove("URL")
	} else {
		wrapper.Set("URL", remote.String())
	}
	wrapper._remoteURL = remote
}

// Override
func (wrapper *PortableNetworkFileWrapper) Password() DecryptKey {
	pwd := wrapper._password
	if pwd == nil {
		info := wrapper.Get("key")
		pwd = ParseSymmetricKey(info)
		wrapper._password = pwd
	}
	return pwd
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetPassword(pwd DecryptKey) {
	wrapper.Remove("key")
	//if pwd != nil {
	//	wrapper.Set("key", pwd.Map())
	//}
	wrapper._password = pwd
}
