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
package protocol

import (
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

/**
 *  Top-Secret message: {
 *      type : 0xFF,
 *      sn   : 456,
 *
 *      forward : {...}  // reliable (secure + certified) message
 *  }
 */
type ForwardContent struct {
	BaseContent

	_secret ReliableMessage
}

func (content *ForwardContent) Init(dictionary map[string]interface{}) *ForwardContent {
	if content.BaseContent.Init(dictionary) != nil {
		// init
	}
	return content
}

func (content *ForwardContent) InitWithMessage(msg ReliableMessage) *ForwardContent {
	if content.InitWithType(FORWARD) != nil {
		dict := msg.GetMap(false)
		content.Set("forward", dict)
	}
	return content
}

func (content *ForwardContent) GetMessage() ReliableMessage {
	if content._secret == nil {
		forward := content.Get("forward")
		if forward == nil {
			panic("forward message not found")
		} else {
			dict, ok := forward.(map[string]interface{})
			if ok {
				content._secret = ReliableMessageParse(dict)
			} else {
				panic("forward message error")
			}
		}
	}
	return content._secret
}
