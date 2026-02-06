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
package dkd

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Top-Secret Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0xFF),
 *      "sn"   : 456,
 *
 *      "forward" : {...}  // reliable (secure + certified) message
 *      "secrets" : [...]  // reliable (secure + certified) messages
 *  }
 *  </pre></blockquote>
 */
type SecretContent struct {
	//ForwardContent
	BaseContent

	_forward ReliableMessage
	_secrets []ReliableMessage
}

func (content *SecretContent) InitWithMap(dict StringKeyMap) ForwardContent {
	if content.BaseContent.InitWithMap(dict) != nil {
		// lazy load
		content._forward = nil
		content._secrets = nil
	}
	return content
}

func (content *SecretContent) InitWithMessage(msg ReliableMessage) ForwardContent {
	if content.BaseContent.InitWithType(ContentType.FORWARD) != nil {
		content._forward = msg
		content._secrets = nil
		//content.Set("forward", msg.Map())
	}
	return content
}

func (content *SecretContent) InitWithMessages(messages []ReliableMessage) ForwardContent {
	if content.BaseContent.InitWithType(ContentType.FORWARD) != nil {
		content._forward = nil
		content._secrets = messages
		//content.Set("secrets", ReliableMessageRevert(messages))
	}
	return content
}

// Override
func (content *SecretContent) Map() StringKeyMap {
	msg := content._forward
	messages := content._secrets
	if messages != nil {
		// serialize 'secrets'
		if !content.Contains("secrets") {
			content.Set("secrets", ReliableMessageRevert(messages))
		}
	} else if msg != nil {
		// serialize 'forward'
		if !content.Contains("forward") {
			content.Set("forward", msg.Map())
		}
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *SecretContent) ForwardMessage() ReliableMessage {
	msg := content._forward
	if msg == nil {
		info := content.Get("forward")
		msg = ParseReliableMessage(info)
		content._forward = msg
	}
	return msg
}

// Override
func (content *SecretContent) SecretMessages() []ReliableMessage {
	messages := content._secrets
	if messages == nil {
		info := content.Get("secrets")
		if info != nil {
			// get from secrets
			messages = ReliableMessageConvert(info)
		} else {
			// get from 'forward'
			msg := content.ForwardMessage()
			if msg != nil {
				messages = []ReliableMessage{msg}
			} else {
				messages = []ReliableMessage{}
			}
		}
		content._secrets = messages
	}
	return messages
}

/**
 *  Combine Forward Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0xCF),
 *      "sn"   : 123,
 *
 *      "title"    : "...",  // chat title
 *      "messages" : [...]   // chat history
 *  }
 *  </pre></blockquote>
 */
type CombineForwardContent struct {
	//CombineContent
	BaseContent

	_history []InstantMessage
}

func (content *CombineForwardContent) InitWithMap(dict StringKeyMap) CombineContent {
	if content.BaseContent.InitWithMap(dict) != nil {
		// lazy load
		content._history = nil
	}
	return content
}

func (content *CombineForwardContent) Init(title string, messages []InstantMessage) CombineContent {
	if content.BaseContent.InitWithType(ContentType.COMBINE_FORWARD) != nil {
		// chat name
		content.Set("title", title)
		// chat history
		content._history = messages // lazy serialize
	}
	return content
}

// Override
func (content *CombineForwardContent) Map() StringKeyMap {
	// serialize 'messages'
	messages := content._history
	if messages != nil && !content.Contains("messages") {
		content.Set("messages", InstantMessageRevert(messages))
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *CombineForwardContent) Title() string {
	return content.GetString("title", "")
}

// Override
func (content *CombineForwardContent) Messages() []InstantMessage {
	messages := content._history
	if messages == nil {
		array := content.Get("messages")
		if array != nil {
			messages = InstantMessageConvert(array)
		} else {
			messages = []InstantMessage{}
		}
		content._history = messages
	}
	return messages
}

/**
 *  Array Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0xCA),
 *      "sn"   : 123,
 *
 *      "contents" : [...]  // content array
 *  }
 *  </pre></blockquote>
 */
type ListContent struct {
	//ArrayContent
	BaseContent

	_list []Content
}

func (content *ListContent) InitWithMap(dict StringKeyMap) ArrayContent {
	if content.BaseContent.InitWithMap(dict) != nil {
		// lazy load
		content._list = nil
	}
	return content
}

func (content *ListContent) Init(contents []Content) ArrayContent {
	if content.BaseContent.InitWithType(ContentType.ARRAY) != nil {
		// content list
		content._list = contents // lazy serialize
	}
	return content
}

// Override
func (content *ListContent) Map() StringKeyMap {
	// serialize 'contents'
	contents := content._list
	if contents != nil && !content.Contains("contents") {
		content.Set("contents", ContentRevert(contents))
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *ListContent) Contents() []Content {
	contents := content._list
	if contents == nil {
		array := content.Get("contents")
		if array != nil {
			contents = ContentConvert(array)
		} else {
			contents = []Content{}
		}
		content._list = contents
	}
	return contents
}
