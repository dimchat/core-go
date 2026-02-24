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
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

// NetworkMessage (ReliableMessage) represents a signed and encrypted message
//
// # Extends EncryptedMessage with a digital signature for authenticity verification
//
// The signature is generated using the sender's private key to prove message origin
//
//	Data structure: {
//	    // Envelope metadata
//	    "sender"   : "moki@xxx",
//	    "receiver" : "hulk@yyy",
//	    "time"     : 123,
//
//	    // Encrypted content and keys (from EncryptedMessage)
//	    "data"     : "...",    // base64_encode( symmetric_encrypt(content))
//	    "keys"     : {
//	        "{ID}"   : "...",  // base64_encode(asymmetric_encrypt(pwd))
//	        "digest" : "..."   // hash(pwd.data)
//	    },
//
//	    // Digital signature for authenticity
//	    "signature": "..."   // base64_encode(asymmetric_sign(data))
//	}
type NetworkMessage struct {
	//ReliableMessage
	*EncryptedMessage

	signature TransportableData
}

func NewNetworkMessage(dict StringKeyMap, data, signature TransportableData) *NetworkMessage {
	return &NetworkMessage{
		EncryptedMessage: NewEncryptedMessage(dict, data),
		signature:        signature,
	}
}

//-------- IReliableMessage

// Override
func (msg *NetworkMessage) Signature() TransportableData {
	ted := msg.signature
	if ted == nil {
		base64 := msg.Get("signature")
		ted = ParseTransportableData(base64)
		msg.signature = ted
	}
	return ted
}

// Override
func (msg *NetworkMessage) Map() StringKeyMap {
	// serialize 'signature'
	ted := msg.signature
	if ted != nil && !msg.Contains("signature") {
		msg.Set("signature", ted.Serialize())
	}
	// OK
	return msg.EncryptedMessage.Map()
}
