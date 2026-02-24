/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
package dkd

import (
	. "github.com/dimchat/core-go/format"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

// EncryptedMessage (SecureMessage) represents an end-to-end encrypted message
//
// # Extends BaseMessage with symmetrically encrypted content and asymmetrically encrypted keys
//
// Core security design: Content is encrypted with a symmetric key, which is then
// encrypted with the receiver's public key (stored in the "keys" field)
//
//	Data structure: {
//	    // Envelope metadata
//	    "sender"   : "moki@xxx",
//	    "receiver" : "hulk@yyy",
//	    "time"     : 123,
//
//	    // Encrypted content and keys
//	    "data"     : "...",    // base64_encode( symmetric_encrypt(content))
//	    "keys"     : {
//	        "{ID}"   : "...",  // base64_encode(asymmetric_encrypt(pwd))
//	        "digest" : "..."   // hash(pwd.data)
//	    }
//	}
type EncryptedMessage struct {
	//SecureMessage
	*BaseMessage

	// data stores the Base64-encoded symmetrically encrypted message content
	//
	// Generated via: base64_encode(symmetric_encrypt(plaintext_content))
	data TransportableData
}

func NewEncryptedMessage(dict StringKeyMap, data TransportableData) *EncryptedMessage {
	return &EncryptedMessage{
		BaseMessage: NewBaseMessage(dict, nil),
		data:        data,
	}
}

//-------- ISecureMessage

// Override
func (msg *EncryptedMessage) Data() TransportableData {
	ted := msg.data
	if ted == nil {
		text := msg.Get("data")
		if text == nil {
			//panic(fmt.Sprintf("message data not found: %v", msg.Map()))
		} else if !IsBroadcastMessage(msg) {
			// message content had been encrypted by a symmetric key,
			// so the data should be encoded here (with algorithm 'base64' as default).
			ted = ParseTransportableData(text)
		} else if str, ok := text.(string); ok {
			ted = NewPlainDataWithString(str)
		} else {
			//panic(fmt.Sprintf("content data error: %v", text))
		}
		msg.data = ted
	}
	return ted
}

// Override
func (msg *EncryptedMessage) EncryptedKeys() StringKeyMap {
	keys := msg.Get("keys")
	if dict, ok := keys.(StringKeyMap); ok {
		return dict
	}
	// TODO: get from 'key'
	return nil
}

// Override
func (msg *EncryptedMessage) Map() StringKeyMap {
	// serialize 'data'
	ted := msg.data
	if ted != nil && !msg.Contains("data") {
		msg.Set("data", ted.Serialize())
	}
	// OK
	return msg.BaseMessage.Map()
}
