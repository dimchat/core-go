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

/**
 *  Instant Message
 *
 *  <blockquote><pre>
 *  data format: {
 *      //-- envelope
 *      "sender"   : "moki@xxx",
 *      "receiver" : "hulk@yyy",
 *      "time"     : 123,
 *
 *      //-- content
 *      "content"  : {...}
 *  }
 *  </pre></blockquote>
 */
type PlainMessage struct {
	//InstantMessage
	BaseMessage

	_content Content
}

func (msg *PlainMessage) InitWithMap(dict StringKeyMap) InstantMessage {
	if msg.BaseMessage.InitWithMap(dict) != nil {
		// lazy load
		msg._content = nil
	}
	return msg
}

func (msg *PlainMessage) Init(head Envelope, body Content) InstantMessage {
	if msg.BaseMessage.InitWithEnvelope(head) != nil {
		msg.SetContent(body)
		//msg._content = body
	}
	return msg
}

//-------- IMessage

// Override
func (msg *PlainMessage) Time() Time {
	head := msg.Envelope()
	body := msg.Content()
	if body == nil {
		//panic("body is nil")
		if head == nil {
			//panic("head is nil")
			return nil
		}
		return head.Time()
	}
	// get body.time or head.time
	when := body.Time()
	if TimeIsNil(when) && head != nil {
		when = head.Time()
	}
	return when
}

// Override
func (msg *PlainMessage) Group() ID {
	body := msg.Content()
	if body == nil {
		//panic("body is nil")
		return nil
	}
	return body.Group()
}

// Override
func (msg *PlainMessage) Type() MessageType {
	body := msg.Content()
	if body == nil {
		//panic("body is nil")
		return ""
	}
	return body.Type()
}

//-------- IInstantMessage

// Override
func (msg *PlainMessage) Content() Content {
	body := msg._content
	if body == nil {
		info := msg.Get("content")
		body = ParseContent(info)
		msg._content = body
	}
	return body
}

// protected
func (msg *PlainMessage) SetContent(content Content) {
	msg.Remove("content")
	//msg.SetMapper("content", content)
	msg._content = content
}

// Override
func (msg *PlainMessage) Map() StringKeyMap {
	// serialize 'content'
	body := msg._content
	if body != nil && !msg.Contains("content") {
		msg.Set("content", body.Map())
	}
	// OK
	return msg.BaseMessage.Map()
}

//
//  Factories
//

func NewInstantMessage(head Envelope, body Content) InstantMessage {
	msg := &PlainMessage{}
	return msg.Init(head, body)
}

func NewInstantMessageWithMap(dict StringKeyMap) InstantMessage {
	msg := &PlainMessage{}
	return msg.InitWithMap(dict)
}
