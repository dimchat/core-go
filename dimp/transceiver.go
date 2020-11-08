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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
)

type Transceiver struct {
	InstantMessageDelegate
	ReliableMessageDelegate

	_entityDelegate EntityDelegate
	_cipherKeyDelegate CipherKeyDelegate
}

func (transceiver *Transceiver) Init() *Transceiver {
	transceiver._entityDelegate = nil
	transceiver._cipherKeyDelegate = nil
	return transceiver
}

func (transceiver Transceiver) EntityDelegate() EntityDelegate {
	return transceiver._entityDelegate
}

func (transceiver *Transceiver) SetEntityDelegate(delegate EntityDelegate) {
	transceiver._entityDelegate = delegate
}

func (transceiver Transceiver) CipherKeyDelegate() CipherKeyDelegate {
	return transceiver._cipherKeyDelegate
}

func (transceiver *Transceiver) SetCipherKeyDelegate(delegate CipherKeyDelegate) {
	transceiver._cipherKeyDelegate = delegate
}

//--------

func checkMessageDelegate(msg Message, delegate MessageDelegate)  {
	holder, ok := msg.(MessageDelegateHolder)
	if ok && holder.Delegate() == nil {
		holder.SetDelegate(delegate)
	}
}

func (transceiver Transceiver) isBroadcast(msg Message) bool {
	checkMessageDelegate(msg, transceiver)

	receiver := msg.Group()
	if receiver == nil {
		receiver = msg.Receiver()
	}
	return AddressIsBroadcast(receiver.Address())
}

func (transceiver Transceiver) getSymmetricKey(from, to ID) SymmetricKey {
	keyCache := transceiver.CipherKeyDelegate()
	// get old key from cache
	key := keyCache.GetCipherKey(from, to)
	if key == nil {
		// create new key and cache it
		key = keyCache.GenerateCipherKey(AES)
		keyCache.CacheCipherKey(from, to, key)
	}
	return key
}

func getOvertGroup(content Content) ID {
	group := content.Group()
	if group == nil {
		return nil
	}
	if AddressIsBroadcast(group.Address()) {
		// broadcast message is always overt
		return group
	}
	ct := content.Type()
	if ct == COMMAND || ct == HISTORY {
		// group command should be sent to each member directly, so
		// don't expose group ID
		return nil
	}
	return group
}

//-------- Transform

func (transceiver Transceiver) EncryptMessage(iMsg InstantMessage) SecureMessage {
	checkMessageDelegate(iMsg, transceiver)

	sender := iMsg.Sender()
	receiver := iMsg.Receiver()
	// if 'group' exists and the 'receiver' is a group ID,
	// they must be equal

	// NOTICE: while sending group message, don't split it before encrypting.
	//         this means you could set group ID into message content, but
	//         keep the "receiver" to be the group ID;
	//         after encrypted (and signed), you could split the message
	//         with group members before sending out, or just send it directly
	//         to the group assistant to let it split messages for you!
	//    BUT,
	//         if you don't want to share the symmetric key with other members,
	//         you could split it (set group ID into message content and
	//         set contact ID to the "receiver") before encrypting, this usually
	//         for sending group command to assistant robot, which should not
	//         share the symmetric key (group msg key) with other members.

	// 1. get symmetric key
	group := getOvertGroup(iMsg.Content())
	var password SymmetricKey
	if group == nil {
		// personal message or (group) command
		password = transceiver.getSymmetricKey(sender, receiver)
	} else {
		password = transceiver.getSymmetricKey(sender, group)
	}

	// 2. encrypt 'content' to 'data' for receiver/group members
	var sMsg SecureMessage
	if AddressIsGroup(receiver.Address()) {
		// group message
		delegate := transceiver.EntityDelegate()
		grp := delegate.GetGroup(receiver)
		if grp == nil {
			panic("failed to get group: " + receiver.String())
		}
		sMsg = iMsg.Encrypt(password, grp.GetMembers())
	} else {
		// personal message (or split group message)
		sMsg = iMsg.Encrypt(password, nil)
	}
	if sMsg == nil {
		// public key for encryption not found
		// TODO: suspend this message for waiting receiver's meta
		return nil
	}
	
	// overt group ID
	if group != nil && !receiver.Equal(group) {
		// NOTICE: this help the receiver knows the group ID
		//         when the group message separated to multi-messages,
		//         if don't want the others know you are the group members,
		//         remove it.
		sMsg.Envelope().SetGroup(group)
	}

	// NOTICE: copy content type to envelope
	//         this help the intermediate nodes to recognize message type
	sMsg.Envelope().SetType(iMsg.Type())

	// OK
	return sMsg
}

func (transceiver Transceiver) SignMessage(sMsg SecureMessage) ReliableMessage {
	checkMessageDelegate(sMsg, transceiver)

	// sign 'data' by sender
	return sMsg.Sign()
}

func (transceiver Transceiver) SerializeMessage(rMsg ReliableMessage) []byte {
	dict := rMsg.GetMap(false)
	return JSONBytesFromMap(dict)
}

func (transceiver Transceiver) DeserializeMessage(data []byte) ReliableMessage {
	dict := JSONMapFromBytes(data)
	// TODO: translate short keys
	//       'S' -> 'sender'
	//       'R' -> 'receiver'
	//       'W' -> 'time'
	//       'T' -> 'type'
	//       'G' -> 'group'
	//       ------------------
	//       'D' -> 'data'
	//       'V' -> 'signature'
	//       'K' -> 'key'
	//       ------------------
	//       'M' -> 'meta'
	return CreateReliableMessage(dict)
}

func (transceiver Transceiver) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	checkMessageDelegate(rMsg, transceiver)

	//
	//  TODO: check [Meta Protocol]
	//        make sure the sender's meta exists
	//        (do in by application)
	//

	// verify 'data' with 'signature'
	return rMsg.Verify()
}

func (transceiver Transceiver) DecryptMessage(sMsg SecureMessage) InstantMessage {
	checkMessageDelegate(sMsg, transceiver)

	//
	//  NOTICE: make sure the receiver is YOU!
	//          which means the receiver's private key exists;
	//          if the receiver is a group ID, split it first
	//

	// decrypt 'data' to 'content'
	return sMsg.Decrypt()

	// TODO: check top-secret message
	//       (do it by application)
}

//-------- MessageDelegate

func (transceiver Transceiver) GetID(identifier interface{}) ID {
	delegate := transceiver.EntityDelegate()
	return delegate.GetID(identifier)
}

//-------- InstantMessageDelegate

func (transceiver Transceiver) GetContent(dictionary interface{}) Content {
	dict := dictionary.(map[string]interface{})
	content := new(BaseContent).Init(dict)
	content.SetDelegate(transceiver)
	return content
}

func (transceiver Transceiver) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         before serialize content, this job should be do in subclass
	dict := content.GetMap(false)
	return JSONBytesFromMap(dict)
}

func (transceiver Transceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	return password.Encrypt(data)
}

func (transceiver Transceiver) EncodeData(data []byte, iMsg InstantMessage) string {
	if transceiver.isBroadcast(iMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so no need to encode to Base64 here
		return UTF8StringFromBytes(data)
	}
	return Base64Encode(data)
}

func (transceiver Transceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	if transceiver.isBroadcast(iMsg) {
		// broadcast message has no key
		return nil
	}
	dict := password.GetMap(false)
	return JSONBytesFromMap(dict)
}

func (transceiver Transceiver) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	// TODO: make sure the receiver's public key exists
	delegate := transceiver.EntityDelegate()
	contact := delegate.GetUser(receiver)
	// encrypt with receiver's public key
	return contact.Encrypt(data)
}

func (transceiver Transceiver) EncodeKey(key []byte, iMsg InstantMessage) string {
	return Base64Encode(key)
}

//-------- SecureMessageDelegate

func (transceiver Transceiver) DecodeKey(key string, sMsg SecureMessage) []byte {
	return Base64Decode(key)
}

func (transceiver Transceiver) DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte {
	// NOTICE: the receiver will be group ID in a group message here
	delegate := transceiver.EntityDelegate()
	user := delegate.GetUser(sMsg.Receiver())
	// decrypt key data with the receiver/group member's private key
	return user.Decrypt(key)
}

func (transceiver Transceiver) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	// NOTICE: the receiver will be group ID in a group message here
	keyCache := transceiver.CipherKeyDelegate()
	if key == nil {
		// get key from cache
		return keyCache.GetCipherKey(sender, receiver)
	} else {
		dict := JSONMapFromBytes(key)
		// TODO: translate short keys
		//       'A' -> 'algorithm'
		//       'D' -> 'data'
		//       'V' -> 'iv'
		//       'M' -> 'mode'
		//       'P' -> 'padding'
		return keyCache.ParseCipherKey(dict)
	}
}

func (transceiver Transceiver) DecodeData(data string, sMsg SecureMessage) []byte {
	if transceiver.isBroadcast(sMsg) {
		// broadcast message content will not be encrypted (just encoded to JsON),
		// so return the string data directly
		return UTF8BytesFromString(data)
	}
	return Base64Decode(data)
}

func (transceiver Transceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	return password.Decrypt(data)
}

func (transceiver Transceiver) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	dict := JSONMapFromBytes(data)
	// TODO: translate short keys
	//       'T' -> 'type'
	//       'N' -> 'sn'
	//       'G' -> 'group'
	content := transceiver.GetContent(dict)

	if transceiver.isBroadcast(sMsg) {
		keyCache := transceiver.CipherKeyDelegate()
		sender := sMsg.Sender()
		group := getOvertGroup(content)
		if group == nil {
			// personal message or (group) command
			// cache key with direction (sender -> receiver)
			receiver := sMsg.Receiver()
			keyCache.CacheCipherKey(sender, receiver, password)
		} else {
			// group message (excludes group command)
			// cache the key with direction (sender -> group)
			keyCache.CacheCipherKey(sender, group, password)
		}
	}

	// NOTICE: check attachment for File/Image/Audio/Video message content
	//         after deserialize content, this job should be do in subclass
	return content
}

func (transceiver Transceiver) SignData(data []byte, sender ID, sMsg SecureMessage) []byte {
	delegate := transceiver.EntityDelegate()
	user := delegate.GetUser(sender)
	return user.Sign(data)
}

func (transceiver Transceiver) EncodeSignature(signature []byte, sMsg SecureMessage) string {
	return Base64Encode(signature)
}

//-------- ReliableMessageDelegate

func (transceiver Transceiver) DecodeSignature(signature string, rMsg ReliableMessage) []byte {
	return Base64Decode(signature)
}

func (transceiver Transceiver) VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool {
	delegate := transceiver.EntityDelegate()
	contact := delegate.GetUser(sender)
	return contact.Verify(data, signature)
}
