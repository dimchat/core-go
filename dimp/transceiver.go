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
package dimp

import (
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

type ITransceiver interface {
	InstantMessageDelegate
	//SecureMessageDelegate
	ReliableMessageDelegate

	// ref(barrack)
	EntityDelegate() EntityDelegate
	SetEntityDelegate(barrack EntityDelegate)
}

/**
 *  Message Transceiver
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Transceiver struct {

	_barrack EntityDelegate
}

func (transceiver *Transceiver) Init() *Transceiver {
	transceiver._barrack = nil
	return transceiver
}

func (transceiver *Transceiver) EntityDelegate() EntityDelegate {
	return transceiver._barrack
}
func (transceiver *Transceiver) SetEntityDelegate(barrack EntityDelegate) {
	transceiver._barrack = barrack
}

func (transceiver *Transceiver) IsBroadcast(msg Message) bool {
	receiver := msg.Group()
	if receiver == nil {
		receiver = msg.Receiver()
	}
	return receiver.IsBroadcast()
}

//-------- IInstantMessageDelegate

func (transceiver *Transceiver) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be do in subclass
	dict := content.GetMap(false)
	return JSONEncodeMap(dict)
}

func (transceiver *Transceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	return password.Encrypt(data)
}

func (transceiver *Transceiver) EncodeData(data []byte, iMsg InstantMessage) string {
	if transceiver.IsBroadcast(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so no need to encode to Base64 here
		return UTF8Decode(data)
	}
	return Base64Encode(data)
}

func (transceiver *Transceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if transceiver.IsBroadcast(iMsg) {
		// broadcast message has no key
		return nil
	}
	dict := password.GetMap(false)
	return JSONEncodeMap(dict)
}

func (transceiver *Transceiver) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	// TODO: make sure the receiver's public key exists
	barrack := transceiver.EntityDelegate()
	contact := barrack.GetUser(receiver)
	// encrypt with receiver's public key
	return contact.Encrypt(data)
}

func (transceiver *Transceiver) EncodeKey(data []byte, iMsg InstantMessage) string {
	return Base64Encode(data)
}

//-------- ISecureMessageDelegate

func (transceiver *Transceiver) DecodeKey(key interface{}, sMsg SecureMessage) []byte {
	if ValueIsNil(key) {
		return nil
	}
	return Base64Decode(key.(string))
}

func (transceiver *Transceiver) DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte {
	// NOTICE: the receiver will be group ID in a group message here
	barrack := transceiver.EntityDelegate()
	user := barrack.GetUser(sMsg.Receiver())
	// decrypt key data with the receiver/group member's private key
	return user.Decrypt(key)
}

func (transceiver *Transceiver) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	// NOTICE: the receiver will be group ID in a group message here
	dict := JSONDecodeMap(key)
	// TODO: translate short keys
	//       'A' -> 'algorithm'
	//       'D' -> 'data'
	//       'V' -> 'iv'
	//       'M' -> 'mode'
	//       'P' -> 'padding'
	return SymmetricKeyParse(dict)
}

func (transceiver *Transceiver) DecodeData(data interface{}, sMsg SecureMessage) []byte {
	text := data.(string)
	if transceiver.IsBroadcast(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so return the string data directly
		return UTF8Encode(text)
	}
	return Base64Decode(text)
}

func (transceiver *Transceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	return password.Decrypt(data)
}

func (transceiver *Transceiver) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	dict := JSONDecodeMap(data)
	// TODO: translate short keys
	//       'T' -> 'type'
	//       'N' -> 'sn'
	//       'G' -> 'group'
	return ContentParse(dict)
}

func (transceiver *Transceiver) SignData(data []byte, sender ID, sMsg SecureMessage) []byte {
	barrack := transceiver.EntityDelegate()
	user := barrack.GetUser(sender)
	return user.Sign(data)
}

func (transceiver *Transceiver) EncodeSignature(signature []byte, sMsg SecureMessage) string {
	return Base64Encode(signature)
}

//-------- IReliableMessageDelegate

func (transceiver *Transceiver) DecodeSignature(signature interface{}, rMsg ReliableMessage) []byte {
	if ValueIsNil(signature) {
		// should not happen
		return nil
	}
	return Base64Decode(signature.(string))
}

func (transceiver *Transceiver) VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool {
	barrack := transceiver.EntityDelegate()
	contact := barrack.GetUser(sender)
	return contact.Verify(data, signature)
}
