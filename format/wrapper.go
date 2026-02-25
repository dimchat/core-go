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

// TransportableFileWrapper defines the interface for wrapping PNF (Portable Network File) data
//
// Provides a standardized way to access and modify core PNF properties (data, filename, URL, decrypt key)
// Implements serialization to StringKeyMap for data interchange
type TransportableFileWrapper interface {

	// Map serializes the PNF wrapper data to a StringKeyMap
	//
	// Returns a key-value representation of all PNF properties for network transmission/storage
	Map() StringKeyMap

	/**
	 *  file data
	 */

	// Data returns the raw file content (TransportableData format, unencrypted)
	//
	// Typically Base64-encoded binary data of the file
	Data() TransportableData
	SetData(data TransportableData)

	/**
	 *  file name
	 */

	// Filename returns the name of the file (including extension)
	//
	// Example: "photo.png", "document.pdf", "audio.mp3"
	Filename() string
	SetFilename(filename string)

	/**
	 *  download URL
	 */

	// URL returns the CDN download URL for the file (alternative to direct data)
	//
	// Used when file content is stored remotely (instead of embedding in Data())
	URL() URL
	SetURL(url URL)

	/**
	 *  decrypt key
	 */

	// Password returns the symmetric decrypt key for CDN-hosted encrypted files
	//
	// Maps to the "key" field in the PNF data structure (required for encrypted remote files)
	Password() DecryptKey
	SetPassword(password DecryptKey)
}

func CreateTransportableFileWrapper(content StringKeyMap,
	data TransportableData, filename string, url URL, password DecryptKey,
) TransportableFileWrapper {
	factory := sharedTransportableFileWrapperFactory
	return factory.CreateTransportableFileWrapper(content, data, filename, url, password)
}

/**
 *  Wrapper Factory
 */

type TransportableFileWrapperFactory interface {

	// Create PNF Wrapper
	CreateTransportableFileWrapper(content StringKeyMap,
		data TransportableData, filename string, url URL, password DecryptKey,
	) TransportableFileWrapper
}

var sharedTransportableFileWrapperFactory TransportableFileWrapperFactory = &pnfWrapperFactory{}

func SetTransportableFileWrapperFactory(factory TransportableFileWrapperFactory) {
	sharedTransportableFileWrapperFactory = factory
}

func GetTransportableFileWrapperFactory() TransportableFileWrapperFactory {
	return sharedTransportableFileWrapperFactory
}

type pnfWrapperFactory struct {
	//TransportableFileWrapperFactory
}

func (pnfWrapperFactory) CreateTransportableFileWrapper(content StringKeyMap,
	data TransportableData, filename string, url URL, key DecryptKey,
) TransportableFileWrapper {
	return NewPortableNetworkFileWrapper(content, data, filename, url, key)
}

// PortableNetworkFileWrapper is a concrete implementation of TransportableFileWrapper (Mixin)
//
// Provides storage and access to PNF (Portable Network File) data with CDN support and encryption
//
//	Standard PNF data structure: {
//	    "data"     : "...",        // Base64-encoded raw file content (unencrypted)
//	    "filename" : "photo.png",  // File name with extension
//
//	    "URL"      : "http://...", // Optional CDN download URL (alternative to "data" field)
//	    "key"      : {             // Optional symmetric decrypt key (for CDN-hosted encrypted files)
//	        "algorithm" : "AES",   // Encryption algorithm (e.g., DES, AES)
//	        "data"      : "{BASE64_ENCODE}",
//	        ...
//	    }
//	 }
//
// Security Note: Files uploaded to public CDNs MUST be encrypted with a symmetric key
// (the "key" field provides decryption capability for authorized recipients)
type PortableNetworkFileWrapper struct {
	//TransportableFileWrapper

	// dictionary stores the serialized StringKeyMap representation of the PNF data
	//
	// Used for network transmission/storage and property serialization
	dictionary StringKeyMap

	// attachment stores the raw, unencrypted file content (TransportableData format)
	//
	// This field holds the direct file data (Base64-encoded) when not using CDN URL
	attachment TransportableData

	// remoteURL stores the CDN download URL for remote file access
	//
	// Alternative to embedding file data in the "data" field
	remoteURL URL

	// password stores the symmetric decrypt key for CDN-hosted encrypted files
	//
	// Required to decrypt files downloaded from public CDNs (maps to "key" field in data structure)
	password DecryptKey
}

func NewPortableNetworkFileWrapper(dict StringKeyMap,
	data TransportableData, filename string, url URL, password DecryptKey,
) *PortableNetworkFileWrapper {
	if filename != "" {
		dict["filename"] = filename
	}
	if url != nil {
		dict["url"] = url.String()
	}
	return &PortableNetworkFileWrapper{
		dictionary: dict,
		attachment: data,
		remoteURL:  url,
		password:   password,
	}
}

func (wrapper *PortableNetworkFileWrapper) Get(key string) any {
	value, exists := wrapper.dictionary[key]
	if !exists {
		return nil
	}
	return value
}

func (wrapper *PortableNetworkFileWrapper) Set(key string, value any) {
	if ValueIsNil(value) {
		delete(wrapper.dictionary, key)
	} else {
		wrapper.dictionary[key] = value
	}
}

func (wrapper *PortableNetworkFileWrapper) Remove(key string) {
	delete(wrapper.dictionary, key)
}

func (wrapper *PortableNetworkFileWrapper) GetString(key string) string {
	value := wrapper.Get(key)
	return ConvertString(value, "")
}

//-------- TransportableFileWrapper

// Override
func (wrapper *PortableNetworkFileWrapper) Map() StringKeyMap {
	// serialize 'data'
	ted := wrapper.attachment
	if ted != nil && wrapper.Get("data") == nil {
		wrapper.Set("data", ted.Serialize())
	}
	// serialize 'key'
	pwd := wrapper.password
	if pwd != nil && wrapper.Get("key") == nil {
		wrapper.Set("key", pwd.Map())
	}
	// OK
	return wrapper.dictionary
}

// Override
func (wrapper *PortableNetworkFileWrapper) Data() TransportableData {
	ted := wrapper.attachment
	if ted == nil {
		base64 := wrapper.Get("data")
		ted = ParseTransportableData(base64)
		wrapper.attachment = ted
	}
	return ted
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetData(ted TransportableData) {
	wrapper.Remove("data")
	//if ted != nil {
	//	wrapper.Set("data", ted.Serialize())
	//}
	wrapper.attachment = ted
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
	remote := wrapper.remoteURL
	if remote == nil {
		locator := wrapper.GetString("URL")
		if locator != "" {
			remote = ParseURL(locator)
			wrapper.remoteURL = remote
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
	wrapper.remoteURL = remote
}

// Override
func (wrapper *PortableNetworkFileWrapper) Password() DecryptKey {
	pwd := wrapper.password
	if pwd == nil {
		info := wrapper.Get("key")
		pwd = ParseSymmetricKey(info)
		wrapper.password = pwd
	}
	return pwd
}

// Override
func (wrapper *PortableNetworkFileWrapper) SetPassword(pwd DecryptKey) {
	wrapper.Remove("key")
	//if pwd != nil {
	//	wrapper.Set("key", pwd.Map())
	//}
	wrapper.password = pwd
}
