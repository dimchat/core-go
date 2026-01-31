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
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Envelope for message
 *  <p>
 *      This class is used to create a message envelope
 *      which contains 'sender', 'receiver' and 'time'
 *  </p>
 *
 *  <blockquote><pre>
 *  data format: {
 *      "sender"   : "moki@xxx",
 *      "receiver" : "hulk@yyy",
 *      "time"     : 123
 *  }
 *  </pre></blockquote>
 */
type MessageEnvelope struct {
	//Envelope
	Dictionary

	_sender   ID
	_receiver ID
	_time     Time
}

func (env *MessageEnvelope) Init(dict StringKeyMap) Envelope {
	if env.Dictionary.Init(dict) != nil {
		// lazy load
		env._sender = nil
		env._receiver = nil
		env._time = nil
	}
	return env
}

//-------- IEnvelope

// Override
func (env *MessageEnvelope) Sender() ID {
	sender := env._sender
	if sender == nil {
		did := env.Get("sender")
		sender = ParseID(did)
		env._sender = sender
	}
	return sender
}

// Override
func (env *MessageEnvelope) Receiver() ID {
	receiver := env._receiver
	if receiver == nil {
		did := env.Get("receiver")
		receiver = ParseID(did)
		if receiver == nil {
			receiver = ANYONE
		}
		env._receiver = receiver
	}
	return receiver
}

// Override
func (env *MessageEnvelope) Time() Time {
	when := env._time
	if when == nil {
		when = env.GetTime("time", nil)
		if when == nil {
			when = TimeNil()
		}
		env._time = when
	}
	return when
}

/*
 *  Group ID
 *  ~~~~~~~~
 *  when a group message was split/trimmed to a single message
 *  the 'receiver' will be changed to a member ID, and
 *  the group ID will be saved as 'group'.
 */
// Override
func (env *MessageEnvelope) Group() ID {
	group := env.Get("group")
	return ParseID(group)
}

// Override
func (env *MessageEnvelope) SetGroup(group ID) {
	env.SetStringer("group", group)
}

/*
 *  Message Type
 *  ~~~~~~~~~~~~
 *  because the message content will be encrypted, so
 *  the intermediate nodes(station) cannot recognize what kind of it.
 *  we pick out the content type and set it in envelope
 *  to let the station do its job.
 */
// Override
func (env *MessageEnvelope) Type() MessageType {
	return env.GetString("type", "")
}

// Override
func (env *MessageEnvelope) SetType(msgType MessageType) {
	env.Set("type", msgType)
}

//
//  Factory
//

func NewEnvelope(from, to ID, when Time) Envelope {
	if to == nil {
		to = ANYONE
	}
	if TimeIsNil(when) {
		when = TimeNow()
	}
	dict := NewMap()
	dict["sender"] = from.String()
	dict["receiver"] = to.String()
	dict["time"] = TimeToFloat64(when)
	// create with dictionary
	env := &MessageEnvelope{}
	if env.Init(dict) != nil {
		env._sender = from
		env._receiver = to
		env._time = when
	}
	return env
}
