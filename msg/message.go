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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/*
 *  Message Transforming
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *     Instant Message <-> Secure Message <-> Reliable Message
 *     +-------------+     +------------+     +--------------+
 *     |  sender     |     |  sender    |     |  sender      |
 *     |  receiver   |     |  receiver  |     |  receiver    |
 *     |  time       |     |  time      |     |  time        |
 *     |             |     |            |     |              |
 *     |  content    |     |  data      |     |  data        |
 *     +-------------+     |  key/keys  |     |  key/keys    |
 *                         +------------+     |  signature   |
 *                                            +--------------+
 *     Algorithm:
 *         data      = password.encrypt(content)
 *         key       = receiver.public_key.encrypt(password)
 *         signature = sender.private_key.sign(data)
 */

/**
 *  Message with Envelope
 *  <p>
 *      Base classes for messages
 *  </p>
 *  <p>
 *      This class is used to create a message
 *      with the envelope fields, such as 'sender', 'receiver', and 'time'
 *  </p>
 *
 *  <blockquote><pre>
 *  data format: {
 *      //-- envelope
 *      "sender"   : "moki@xxx",
 *      "receiver" : "hulk@yyy",
 *      "time"     : 123,
 *
 *      //-- body
 *      ...
 *  }
 *  </pre></blockquote>
 */
type BaseMessage struct {
	//Message
	*Dictionary

	envelope Envelope
}

func NewBaseMessage(dict StringKeyMap, head Envelope) *BaseMessage {
	if dict == nil {
		dict = head.Map()
	}
	return &BaseMessage{
		Dictionary: NewDictionary(dict),
		envelope:   head,
	}
}

//-------- IMessage

// Override
func (msg *BaseMessage) Envelope() Envelope {
	head := msg.envelope
	if head == nil {
		head = ParseEnvelope(msg.Map())
		msg.envelope = head
	}
	return head
}

// Override
func (msg *BaseMessage) Sender() ID {
	head := msg.Envelope()
	if head == nil {
		//panic("head is nil")
		return nil
	}
	return head.Sender()
}

// Override
func (msg *BaseMessage) Receiver() ID {
	head := msg.Envelope()
	if head == nil {
		//panic("head is nil")
		return nil
	}
	return head.Receiver()
}

// Override
func (msg *BaseMessage) Time() Time {
	head := msg.Envelope()
	if head == nil {
		//panic("head is nil")
		return nil
	}
	return head.Time()
}

// Override
func (msg *BaseMessage) Group() ID {
	head := msg.Envelope()
	if head == nil {
		//panic("head is nil")
		return nil
	}
	return head.Group()
}

// Override
func (msg *BaseMessage) Type() MessageType {
	head := msg.Envelope()
	if head == nil {
		//panic("head is nil")
		return ""
	}
	return head.Type()
}

//--------

func IsBroadcastMessage(msg Message) bool {
	receiver := msg.Receiver()
	if receiver == nil || receiver.IsBroadcast() {
		return true
	}
	// check exposed group
	overtGroup := msg.Get("group")
	if overtGroup == nil {
		return false
	}
	group := ParseID(overtGroup)
	return group != nil && group.IsBroadcast()
}
