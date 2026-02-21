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
	*BaseContent

	secrets []ReliableMessage
}

func NewSecretContent(dict StringKeyMap, messages []ReliableMessage) *SecretContent {
	return &SecretContent{
		BaseContent: NewBaseContent(dict, ContentType.FORWARD),
		secrets:     messages,
	}
}

// Override
func (content *SecretContent) Map() StringKeyMap {
	// serialize 'secret' messages
	messages := content.secrets
	if messages != nil && !content.Contains("secrets") {
		content.Set("secrets", ReliableMessageRevert(messages))
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *SecretContent) SecretMessages() []ReliableMessage {
	messages := content.secrets
	if messages == nil {
		secrets := content.Get("secrets")
		if secrets != nil {
			// get from secrets
			messages = ReliableMessageConvert(secrets)
		} else {
			// get from 'forward'
			forward := content.Get("forward")
			msg := ParseReliableMessage(forward)
			if msg != nil {
				messages = []ReliableMessage{msg}
			} else {
				messages = []ReliableMessage{}
			}
		}
		content.secrets = messages
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
	*BaseContent

	history []InstantMessage
}

func NewCombineForwardContent(dict StringKeyMap, title string, messages []InstantMessage) *CombineForwardContent {
	if dict != nil {
		// init combine content with map
		return &CombineForwardContent{
			BaseContent: NewBaseContent(dict, ""),
			// lazy load
			history: nil,
		}
	}
	// new combine content
	content := &CombineForwardContent{
		BaseContent: NewBaseContent(dict, ContentType.COMBINE_FORWARD),
		history:     messages,
	}
	// chat name
	content.Set("title", title)
	return content
}

// Override
func (content *CombineForwardContent) Map() StringKeyMap {
	// serialize 'messages'
	messages := content.history
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
	messages := content.history
	if messages == nil {
		array := content.Get("messages")
		if array != nil {
			messages = InstantMessageConvert(array)
		} else {
			messages = []InstantMessage{}
		}
		content.history = messages
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
	*BaseContent

	list []Content
}

func NewListContent(dict StringKeyMap, contents []Content) *ListContent {
	return &ListContent{
		BaseContent: NewBaseContent(dict, ContentType.ARRAY),
		list:        contents,
	}
}

// Override
func (content *ListContent) Map() StringKeyMap {
	// serialize 'contents'
	contents := content.list
	if contents != nil && !content.Contains("contents") {
		content.Set("contents", ContentRevert(contents))
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *ListContent) Contents() []Content {
	contents := content.list
	if contents == nil {
		array := content.Get("contents")
		if array != nil {
			contents = ContentConvert(array)
		} else {
			contents = []Content{}
		}
		content.list = contents
	}
	return contents
}
