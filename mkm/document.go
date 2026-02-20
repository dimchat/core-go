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
	//Document
	*Dictionary

	data      string            // JSONEncode(properties)
	signature TransportableData // LocalUser(id).Sign(data)

	properties StringKeyMap
	status     int8 // 1 for valid, -1 for invalid
}

func NewBaseDocument(dict StringKeyMap, docType DocumentType, data string, signature TransportableData) *BaseDocument {
	var properties StringKeyMap
	var status int8
	if dict != nil {
		// document info from network, waiting for verify
		properties = nil // lazy load
		status = 0
	} else if data == "" || signature == nil {
		// new empty document
		dict = NewMap()
		// document type
		dict["type"] = docType
		// initialize properties with created time
		properties = NewMap()
		properties["type"] = docType // deprecated
		properties["created_time"] = TimeToFloat64(TimeNow())
		status = 0
	} else {
		// document with data and signature loaded from local storage
		dict = NewMap()
		// document type
		dict["type"] = docType
		// document data (JsON)
		dict["data"] = data
		// document signature (Base64)
		dict["signature"] = signature.Serialize()
		properties = nil // lazy load
		// all documents must be verified before saving into local storage
		status = 1
	}
	return &BaseDocument{
		Dictionary: NewDictionary(dict),
		data:       data,
		signature:  signature,
		properties: properties,
		status:     status,
	}
}

/**
 *  Get serialized properties
 *
 * @return JSON string
 */
func (doc *BaseDocument) getData() string {
	json := doc.data
	if json == "" {
		json = doc.GetString("data", "")
		doc.data = json
	}
	return json
}

/**
 *  Get signature for serialized properties
 *
 * @return signature data
 */
func (doc *BaseDocument) getSignature() TransportableData {
	ted := doc.signature
	if ted == nil {
		base64 := doc.Get("signature")
		ted = ParseTransportableData(base64)
		doc.signature = ted
	}
	return ted
}

//-------- TAI

// Override
func (doc *BaseDocument) IsValid() bool {
	return doc.status > 0
}

// Override
func (doc *BaseDocument) Verify(publicKey VerifyKey) bool {
	//if doc._status > 0 {
	//	// already verify OK
	//	return true
	//}
	data := doc.getData()
	signature := doc.getSignature()
	if data == "" {
		// NOTICE: if data is empty, signature should be empty at the same time
		//         this happen while entity document not found
		if signature == nil || signature.IsEmpty() {
			doc.status = 0
		} else {
			// data signature error
			doc.status = -1
		}
	} else if signature == nil || signature.IsEmpty() {
		// signature error
		doc.status = -1
	} else if publicKey.Verify(UTF8Encode(data), signature.Bytes()) {
		// signature matched
		doc.status = 1
	} else {
		// public key not matched,
		// no need to affect the status here
		return false
	}
	// NOTICE: if status is 0, it doesn't mean the document is invalid,
	//         try another key
	return doc.status == 1
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
	ted := NewBase64DataWithBytes(signature)
	// 3. update 'data' & 'signature' fields
	doc.Set("data", data)                 // JSON string
	doc.Set("signature", ted.Serialize()) // BASE-64
	doc.data = data
	doc.signature = ted
	// 4. update status
	doc.status = 1
	return signature
}

// Override
func (doc *BaseDocument) Properties() StringKeyMap {
	if doc.status < 0 {
		// invalid
		return nil
	}
	info := doc.properties
	if info == nil {
		data := doc.getData()
		if data == "" {
			// create new properties
			info = NewMap()
		} else {
			// get properties from data
			info = JSONDecodeMap(data)
		}
		doc.properties = info
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
	doc.status = 0
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
	doc.data = ""
	doc.signature = nil
}

//-------- IDocument

// Override
func (doc *BaseDocument) Time() Time {
	timestamp := doc.GetProperty("time")
	return ConvertTime(timestamp, nil)
}
