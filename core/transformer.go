/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
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
package core

import (
	"github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Message Transformer Implementations
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
type Transformer struct {
	dimp.Transformer

	_transceiver dimp.Transceiver
}

func (transformer *Transformer) Init(transceiver dimp.Transceiver) *Transformer {
	transformer._transceiver = transceiver
	return transformer
}

func (transformer *Transformer) Transceiver() dimp.Transceiver {
	return transformer._transceiver
}

//-------- InstantMessageDelegate

func isBroadcast(msg Message) bool {
	receiver := msg.Group()
	if receiver == nil {
		receiver = msg.Receiver()
	}
	return receiver.IsBroadcast()
}

func (transformer *Transformer) SerializeContent(content Content, _ SymmetricKey, _ InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be do in subclass
	dict := content.GetMap(false)
	return JSONEncode(dict)
}

func (transformer *Transformer) EncryptContent(data []byte, password SymmetricKey, _ InstantMessage) []byte {
	return password.Encrypt(data)
}

func (transformer *Transformer) EncodeData(data []byte, iMsg InstantMessage) string {
	if isBroadcast(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so no need to encode to Base64 here
		return UTF8Decode(data)
	}
	return Base64Encode(data)
}

func (transformer *Transformer) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if isBroadcast(iMsg) {
		// broadcast message has no key
		return nil
	}
	dict := password.GetMap(false)
	return JSONEncode(dict)
}

func (transformer *Transformer) EncryptKey(data []byte, receiver ID, _ InstantMessage) []byte {
	// TODO: make sure the receiver's public key exists
	contact := transformer.Transceiver().GetUser(receiver)
	// encrypt with receiver's public key
	return contact.Encrypt(data)
}

func (transformer *Transformer) EncodeKey(key []byte, _ InstantMessage) string {
	return Base64Encode(key)
}

//-------- SecureMessageDelegate

func (transformer *Transformer) DecodeKey(key string, _ SecureMessage) []byte {
	return Base64Decode(key)
}

func (transformer *Transformer) DecryptKey(key []byte, _ ID, _ ID, sMsg SecureMessage) []byte {
	// NOTICE: the receiver will be group ID in a group message here
	user := transformer.Transceiver().GetUser(sMsg.Receiver())
	// decrypt key data with the receiver/group member's private key
	return user.Decrypt(key)
}

func (transformer *Transformer) DeserializeKey(key []byte, sender ID, receiver ID, _ SecureMessage) SymmetricKey {
	// NOTICE: the receiver will be group ID in a group message here
	if key == nil {
		// get key from cache
		return transformer.Transceiver().GetCipherKey(sender, receiver, false)
	} else {
		dict := JSONDecode(key)
		// TODO: translate short keys
		//       'A' -> 'algorithm'
		//       'D' -> 'data'
		//       'V' -> 'iv'
		//       'M' -> 'mode'
		//       'P' -> 'padding'
		return SymmetricKeyParse(dict)
	}
}

func (transformer *Transformer) DecodeData(data string, sMsg SecureMessage) []byte {
	if isBroadcast(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so return the string data directly
		return UTF8Encode(data)
	}
	return Base64Decode(data)
}

func (transformer *Transformer) DecryptContent(data []byte, password SymmetricKey, _ SecureMessage) []byte {
	return password.Decrypt(data)
}

func (transformer *Transformer) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	dict := JSONDecode(data)
	// TODO: translate short keys
	//       'T' -> 'type'
	//       'N' -> 'sn'
	//       'G' -> 'group'
	content := ContentParse(dict)

	if !isBroadcast(sMsg) {
		sender := sMsg.Sender()
		group := transformer.Transceiver().GetOvertGroup(content)
		if group == nil {
			// personal message or (group) command
			// cache key with direction (sender -> receiver)
			receiver := sMsg.Receiver()
			transformer.Transceiver().CacheCipherKey(sender, receiver, password)
		} else {
			// group message (excludes group command)
			// cache the key with direction (sender -> group)
			transformer.Transceiver().CacheCipherKey(sender, group, password)
		}
	}

	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be do in subclass
	return content
}

func (transformer *Transformer) SignData(data []byte, sender ID, _ SecureMessage) []byte {
	user := transformer.Transceiver().GetUser(sender)
	return user.Sign(data)
}

func (transformer *Transformer) EncodeSignature(signature []byte, _ SecureMessage) string {
	return Base64Encode(signature)
}

//-------- ReliableMessageDelegate

func (transformer *Transformer) DecodeSignature(signature string, _ ReliableMessage) []byte {
	return Base64Decode(signature)
}

func (transformer *Transformer) VerifyDataSignature(data []byte, signature []byte, sender ID, _ ReliableMessage) bool {
	contact := transformer.Transceiver().GetUser(sender)
	return contact.Verify(data, signature)
}
