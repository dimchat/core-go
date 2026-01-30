/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
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
package mkm

import (
	. "github.com/dimchat/core-go/format"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Base Document
 *  ~~~~~~~~~~~~~
 */
type BaseDocument struct {
	Dictionary

	_data      string            // JSONEncode(properties)
	_signature TransportableData // LocalUser(id).Sign(data)

	_properties StringKeyMap
	_status     int8 // 1 for valid, -1 for invalid
}

/**
 *  Create Entity Document
 *
 *  @param dict - info
 */
func (doc *BaseDocument) Init(dict StringKeyMap) Document {
	if doc.Dictionary.Init(dict) != nil {
		// lazy load
		doc._data = ""
		doc._signature = nil
		doc._properties = nil
		doc._status = 0
	}
	return doc
}

/**
 *  Create with data and signature loaded from local storage
 *
 * @param docType   - document type
 * @param data      - document data in JsON format
 * @param signature - signature of document data in Base64 format
 */
func (doc *BaseDocument) InitWithType(docType DocumentType, data string, signature TransportableData) Document {
	if data == "" || signature == nil || signature.IsEmpty() {
		panic("data or signature is nil")
		return nil
	}
	if doc.Dictionary.Init(NewMap()) != nil {
		// document type
		doc.Set("type", docType)

		// document data (JsON)
		doc.Set("data", data)
		doc._data = data

		// document signature (Base64)
		doc.Set("signature", signature.Serialize())
		doc._signature = signature

		doc._properties = nil // lazy

		// all documents must be verified before saving into local storage
		doc._status = 1
	}
	return doc
}

/**
 *  Create a new empty document
 *
 * @param type - document type
 */
func (doc *BaseDocument) InitEmptyDocument(docType DocumentType) Document {
	if doc.Dictionary.Init(NewMap()) != nil {
		// document type
		doc.Set("type", docType)

		// document data & signature
		doc._data = ""
		doc._signature = nil

		// initialize properties with created time
		info := NewMap()
		info["type"] = docType // deprecated
		info["created_time"] = TimeToFloat64(TimeNow())
		doc._properties = info

		doc._status = 0
	}
	return doc
}

/**
 *  Get serialized properties
 *
 * @return JsON string
 */
func (doc *BaseDocument) data() string {
	json := doc._data
	if json == "" {
		json := doc.GetString("data", "")
		doc._data = json
	}
	return json
}

/**
 *  Get signature for serialized properties
 *
 * @return signature data
 */
func (doc *BaseDocument) signature() TransportableData {
	ted := doc._signature
	if ted == nil {
		base64 := doc.Get("signature")
		ted = ParseTransportableData(base64)
		doc._signature = ted
	}
	return ted
}

//-------- TAI

// Override
func (doc *BaseDocument) IsValid() bool {
	return doc._status > 0
}

// Override
func (doc *BaseDocument) Verify(publicKey VerifyKey) bool {
	//if doc._status > 0 {
	//	// already verify OK
	//	return true
	//}
	data := doc.data()
	signature := doc.signature()
	if data == "" {
		// NOTICE: if data is empty, signature should be empty at the same time
		//         this happen while entity document not found
		if signature == nil || signature.IsEmpty() {
			doc._status = 0
		} else {
			// data signature error
			doc._status = -1
		}
	} else if signature == nil || signature.IsEmpty() {
		// signature error
		doc._status = -1
	} else if publicKey.Verify(UTF8Encode(data), signature.Bytes()) {
		// signature matched
		doc._status = 1
	} else {
		// public key not matched,
		// no need to affect the status here
		return false
	}
	// NOTICE: if status is 0, it doesn't mean the document is invalid,
	//         try another key
	return doc._status == 1
}

// Override
func (doc *BaseDocument) Sign(privateKey SignKey) []byte {
	//if doc._status > 0 {
	//	// already signed/verified
	//	return doc.signature()
	//}
	// 1. update sign time
	doc.SetProperty("time", TimeToFloat64(TimeNow()))
	// 2. encode & sign
	info := doc.Properties()
	if info == nil {
		return nil
	}
	data := JSONEncodeMap(info)
	signature := privateKey.Sign(UTF8Encode(data))
	ted := CreateBase64DataWithBytes(signature)
	// 3. update 'data' & 'signature' fields
	doc.Set("data", data)
	doc.Set("signature", ted.Serialize())
	doc._data = data
	doc._signature = ted
	// 4. update status
	doc._status = 1
	return signature
}

// Override
func (doc *BaseDocument) Properties() StringKeyMap {
	if doc._status < 0 {
		// invalid
		return nil
	}
	info := doc._properties
	if info == nil {
		data := doc.data()
		if data == "" {
			// create new properties
			info = NewMap()
		} else {
			// get properties from data
			info = JSONDecodeMap(data)
		}
		doc._properties = info
	}
	return info
}

// Override
func (doc *BaseDocument) GetProperty(name string) interface{} {
	dict := doc.Properties()
	if dict == nil {
		return nil
	}
	value, exists := dict[name]
	if !exists {
		return nil
	}
	return value
}

// Override
func (doc *BaseDocument) SetProperty(name string, value interface{}) {
	// 1. reset status
	doc._status = 0
	// 2. update property value with name
	properties := doc.Properties()
	if properties == nil {
		// failed to get properties
	} else if ValueIsNil(value) {
		delete(properties, name)
	} else {
		properties[name] = value
	}
	// 3. clear data signature after properties changed
	doc.Remove("data")
	doc.Remove("signature")
	doc._data = ""
	doc._signature = nil
}

//-------- IDocument

// Override
func (doc *BaseDocument) Time() Time {
	timestamp := doc.GetProperty("time")
	return ConvertTime(timestamp, nil)
	//return TimeParse(timestamp)
}
