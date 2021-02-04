/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
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
package core

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Message Packer
 *  ~~~~~~~~~~~~~~
 */
type Packer interface {

	/**
	 *  Get group ID which should be exposed to public network
	 *
	 * @param content - message content
	 * @return exposed group ID
	 */
	GetOvertGroup(content Content) ID

	//
	//  InstantMessage -> SecureMessage -> ReliableMessage -> Data
	//

	/**
	 *  Encrypt message content
	 *
	 * @param iMsg - plain message
	 * @return encrypted message
	 */
	EncryptMessage(iMsg InstantMessage) SecureMessage

	/**
	 *  Sign content data
	 *
	 * @param sMsg - encrypted message
	 * @return network message
	 */
	SignMessage(sMsg SecureMessage) ReliableMessage

	/**
	 *  Serialize network message
	 *
	 * @param rMsg - network message
	 * @return data package
	 */
	SerializeMessage(rMsg ReliableMessage) []byte

	//
	//  Data -> ReliableMessage -> SecureMessage -> InstantMessage
	//

	/**
	 *  Deserialize network message
	 *
	 * @param data - data package
	 * @return network message
	 */
	DeserializeMessage(data []byte) ReliableMessage

	/**
	 *  Verify encrypted content data
	 *
	 * @param rMsg - network message
	 * @return encrypted message
	 */
	VerifyMessage(rMsg ReliableMessage) SecureMessage

	/**
	 *  Decrypt message content
	 *
	 * @param sMsg - encrypted message
	 * @return plain message
	 */
	DecryptMessage(sMsg SecureMessage) InstantMessage
}

//
//   Implementations
//

type TransceiverPacker struct {
	Packer

	_transceiver *Transceiver
}

func (packer *TransceiverPacker) Init(transceiver *Transceiver) *TransceiverPacker {
	packer._transceiver = transceiver
	return packer
}

func (packer *TransceiverPacker) Transceiver() *Transceiver {
	return packer._transceiver
}

func (packer *TransceiverPacker) GetOvertGroup(content Content) ID {
	group := content.Group()
	if group == nil {
		return nil
	}
	if group.IsBroadcast() {
		// broadcast message is always overt
		return group
	}
	msgType := content.Type()
	if msgType == COMMAND || msgType == HISTORY {
		// group command should be sent to each member directly, so
		// don't expose group ID
		return nil
	}
	return group
}

func (packer *TransceiverPacker) EncryptMessage(iMsg InstantMessage) SecureMessage {
	// check message delegate
	if iMsg.Delegate() == nil {
		iMsg.SetDelegate(packer.Transceiver())
	}
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
	group := packer.Transceiver().GetOvertGroup(iMsg.Content())
	var password SymmetricKey
	if group == nil {
		// personal message or (group) command
		password = packer.Transceiver().GetCipherKey(sender, receiver, true)
	} else {
		password = packer.Transceiver().GetCipherKey(sender, group, true)
	}

	// 2. encrypt 'content' to 'data' for receiver/group members
	var sMsg SecureMessage
	if receiver.IsGroup() {
		// group message
		grp := packer.Transceiver().GetGroup(receiver)
		if grp == nil {
			// group not ready
			// TODO: suspend this message for waiting group's meta
			return nil
		}
		members := grp.GetMembers()
		if members == nil || len(members) == 0 {
			// group members not found
			// TODO: suspend this message for waiting group's membership
			return nil
		}
		sMsg = iMsg.Encrypt(password, members)
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
	sMsg.Envelope().SetType(iMsg.Content().Type())

	// OK
	return sMsg
}

func (packer *TransceiverPacker) SignMessage(sMsg SecureMessage) ReliableMessage {
	// check message delegate
	if sMsg.Delegate() == nil {
		sMsg.SetDelegate(packer.Transceiver())
	}
	// sign 'data' by sender
	return sMsg.Sign()
}

func (packer *TransceiverPacker) SerializeMessage(rMsg ReliableMessage) []byte {
	dict := rMsg.GetMap(false)
	return JSONEncode(dict)
}

func (packer *TransceiverPacker) DeserializeMessage(data []byte) ReliableMessage {
	dict := JSONDecode(data)
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
	return ReliableMessageParse(dict)
}

func (packer *TransceiverPacker) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	// check message delegate
	if rMsg.Delegate() == nil {
		rMsg.SetDelegate(packer.Transceiver())
	}
	//
	//  TODO: check [Meta Protocol]
	//        make sure the sender's meta exists
	//        (do in by application)
	//

	// verify 'data' with 'signature'
	return rMsg.Verify()
}

func (packer *TransceiverPacker) DecryptMessage(sMsg SecureMessage) InstantMessage {
	// check message delegate
	if sMsg.Delegate() == nil {
		sMsg.SetDelegate(packer.Transceiver())
	}
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
