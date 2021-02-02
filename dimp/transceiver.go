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
)

type Transceiver struct {
	InstantMessageDelegate
	ReliableMessageDelegate

	_entityDelegate EntityDelegate
	_keyDelegate CipherKeyDelegate
	_processor Processor
	_packer Packer
}

func (transceiver *Transceiver) Init() *Transceiver {
	transceiver._entityDelegate = nil
	transceiver._keyDelegate = nil
	transceiver._processor = nil
	transceiver._packer = nil
	return transceiver
}

/**
 *  Delegate for User/Group
 *
 * @param barrack - entity delegate
 */
func (transceiver *Transceiver) SetEntityDelegate(barrack EntityDelegate) {
	transceiver._entityDelegate = barrack
}
func (transceiver Transceiver) EntityDelegate() EntityDelegate {
	return transceiver._entityDelegate
}

/**
 *  Delegate for Cipher Key
 *
 * @param keyCache - key store
 */
func (transceiver *Transceiver) SetCipherKeyDelegate(keyCache CipherKeyDelegate) {
	transceiver._keyDelegate = keyCache
}
func (transceiver Transceiver) CipherKeyDelegate() CipherKeyDelegate {
	return transceiver._keyDelegate
}

/**
 *  Delegate for Processing Message
 *
 * @param processor - message processor
 */
func (transceiver *Transceiver) SetProcessor(processor Processor) {
	transceiver._processor = processor
}
func (transceiver *Transceiver) Processor() Processor {
	return transceiver._processor
}

/**
 *  Delegate for Packing Message
 *
 * @param packer - message packer
 */
func (transceiver *Transceiver) SetPacker(packer Packer) {
	transceiver._packer = packer
}
func (transceiver *Transceiver) Packer() Packer {
	return transceiver._packer
}

//
//  Interfaces for User/Group
//
func (transceiver *Transceiver) SelectLocalUser(receiver ID) *User {
	return transceiver.EntityDelegate().SelectLocalUser(receiver)
}

func (transceiver *Transceiver) GetUser(identifier ID) *User {
	return transceiver.EntityDelegate().GetUser(identifier)
}

func (transceiver *Transceiver) GetGroup(identifier ID) *Group {
	return transceiver.EntityDelegate().GetGroup(identifier)
}

//
//  Interfaces for Cipher Key
//
func (transceiver *Transceiver) GetCipherKey(sender, receiver ID, generate bool) SymmetricKey {
	return transceiver.CipherKeyDelegate().GetCipherKey(sender, receiver, generate)
}

func (transceiver *Transceiver) CacheCipherKey(sender, receiver ID, key SymmetricKey) {
	transceiver.CipherKeyDelegate().CacheCipherKey(sender, receiver, key)
}

//
//  Interfaces for Packing Message
//
func (transceiver *Transceiver) GetOvertGroup(content Content) ID {
	return transceiver.Packer().GetOvertGroup(content)
}

func (transceiver *Transceiver) EncryptMessage(iMsg InstantMessage) SecureMessage {
	return transceiver.Packer().EncryptMessage(iMsg)
}

func (transceiver *Transceiver) SignMessage(sMsg SecureMessage) ReliableMessage {
	return transceiver.Packer().SignMessage(sMsg)
}

func (transceiver *Transceiver) SerializeMessage(rMsg ReliableMessage) []byte {
	return transceiver.Packer().SerializeMessage(rMsg)
}

func (transceiver *Transceiver) DeserializeMessage(data []byte) ReliableMessage {
	return transceiver.Packer().DeserializeMessage(data)
}

func (transceiver *Transceiver) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	return transceiver.Packer().VerifyMessage(rMsg)
}

func (transceiver *Transceiver) DecryptMessage(sMsg SecureMessage) InstantMessage {
	return transceiver.Packer().DecryptMessage(sMsg)
}

//
//  Interfaces for Processing Message
//
func (transceiver *Transceiver) ProcessData(data []byte) []byte {
	return transceiver.Processor().ProcessData(data)
}

func (transceiver *Transceiver) ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage {
	return transceiver.Processor().ProcessReliableMessage(rMsg)
}

func (transceiver *Transceiver) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage {
	return transceiver.Processor().ProcessSecureMessage(sMsg, rMsg)
}

func (transceiver *Transceiver) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	return transceiver.Processor().ProcessInstantMessage(iMsg, rMsg)
}

func (transceiver *Transceiver) ProcessContent(content Content, rMsg ReliableMessage) Content {
	return transceiver.Processor().ProcessContent(content, rMsg)
}

//-------- InstantMessageDelegate

func isBroadcast(msg Message) bool {
	receiver := msg.Group()
	if receiver == nil {
		receiver = msg.Receiver()
	}
	return receiver.IsBroadcast()
}

func (transceiver *Transceiver) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be do in subclass
	dict := content.GetMap(false)
	return JSONEncode(dict)
}

func (transceiver *Transceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	return password.Encrypt(data)
}

func (transceiver *Transceiver) EncodeData(data []byte, iMsg InstantMessage) string {
	if isBroadcast(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so no need to encode to Base64 here
		return UTF8Decode(data)
	}
	return Base64Encode(data)
}

func (transceiver *Transceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if isBroadcast(iMsg) {
		// broadcast message has no key
		return nil
	}
	dict := password.GetMap(false)
	return JSONEncode(dict)
}

func (transceiver *Transceiver) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	// TODO: make sure the receiver's public key exists
	contact := transceiver.GetUser(receiver)
	// encrypt with receiver's public key
	return contact.Encrypt(data)
}

func (transceiver *Transceiver) EncodeKey(key []byte, iMsg InstantMessage) string {
	return Base64Encode(key)
}

//-------- SecureMessageDelegate

func (transceiver *Transceiver) DecodeKey(key string, sMsg SecureMessage) []byte {
	return Base64Decode(key)
}

func (transceiver *Transceiver) DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte {
	// NOTICE: the receiver will be group ID in a group message here
	user := transceiver.GetUser(sMsg.Receiver())
	// decrypt key data with the receiver/group member's private key
	return user.Decrypt(key)
}

func (transceiver *Transceiver) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	// NOTICE: the receiver will be group ID in a group message here
	if key == nil {
		// get key from cache
		return transceiver.GetCipherKey(sender, receiver, false)
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

func (transceiver *Transceiver) DecodeData(data string, sMsg SecureMessage) []byte {
	if isBroadcast(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so return the string data directly
		return UTF8Encode(data)
	}
	return Base64Decode(data)
}

func (transceiver *Transceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	return password.Decrypt(data)
}

func (transceiver *Transceiver) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	dict := JSONDecode(data)
	// TODO: translate short keys
	//       'T' -> 'type'
	//       'N' -> 'sn'
	//       'G' -> 'group'
	content := ContentParse(dict)

	if !isBroadcast(sMsg) {
		sender := sMsg.Sender()
		group := transceiver.GetOvertGroup(content)
		if group == nil {
			// personal message or (group) command
			// cache key with direction (sender -> receiver)
			receiver := sMsg.Receiver()
			transceiver.CacheCipherKey(sender, receiver, password)
		} else {
			// group message (excludes group command)
			// cache the key with direction (sender -> group)
			transceiver.CacheCipherKey(sender, group, password)
		}
	}

	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be do in subclass
	return content
}

func (transceiver *Transceiver) SignData(data []byte, sender ID, sMsg SecureMessage) []byte {
	user := transceiver.GetUser(sender)
	return user.Sign(data)
}

func (transceiver *Transceiver) EncodeSignature(signature []byte, sMsg SecureMessage) string {
	return Base64Encode(signature)
}

//-------- ReliableMessageDelegate

func (transceiver *Transceiver) DecodeSignature(signature string, rMsg ReliableMessage) []byte {
	return Base64Decode(signature)
}

func (transceiver *Transceiver) VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool {
	contact := transceiver.GetUser(sender)
	return contact.Verify(data, signature)
}
