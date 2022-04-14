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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Top-Secret message: {
 *      type : 0xFF,
 *      sn   : 456,
 *
 *      forward : {...}  // reliable (secure + certified) message
 *  }
 */
type SecretContent struct {
	BaseContent

	_secret ReliableMessage
}

func NewForwardContent(msg ReliableMessage) ForwardContent {
	content := new(SecretContent)
	content.InitWithMessage(msg)
	return content
}

func (content *SecretContent) Init(dict map[string]interface{}) ForwardContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._secret = nil
	}
	return content
}

func (content *SecretContent) InitWithMessage(secret ReliableMessage) ForwardContent {
	if content.BaseContent.InitWithType(FORWARD) != nil {
		content.Set("forward", secret.Map())
		content._secret = secret
	}
	return content
}

//-------- IForwardContent

func (content *SecretContent) ForwardMessage() ReliableMessage {
	if content._secret == nil {
		msg := content.Get("forward")
		content._secret = ReliableMessageParse(msg)
	}
	return content._secret
}

/**
 *  Text message: {
 *      type : 0x01,
 *      sn   : 123,
 *
 *      text : "..."
 *  }
 */

type BaseTextContent struct {
	BaseContent
}

func NewTextContent(text string) TextContent {
	content := new(BaseTextContent)
	content.InitWithText(text)
	return content
}

//func (content *BaseTextContent) Init(dict map[string]interface{}) TextContent {
//	if content.BaseContent.Init(dict) != nil {
//	}
//	return content
//}

func (content *BaseTextContent) InitWithText(text string) TextContent {
	if content.InitWithType(TEXT) != nil {
		content.SetText(text)
	}
	return content
}

//-------- ITextContent

func (content *BaseTextContent) Text() string {
	text := content.Get("text")
	if text == nil {
		return ""
	}
	return text.(string)
}

func (content *BaseTextContent) SetText(text string) {
	content.Set("text", text)
}

/**
 *  Web Page message: {
 *      type : 0x20,
 *      sn   : 123,
 *
 *      URL   : "https://github.com/moky/dimp", // Page URL
 *      icon  : "...",                          // base64_encode(icon)
 *      title : "...",
 *      desc  : "..."
 *  }
 */
type WebPageContent struct {
	BaseContent

	_icon []byte
}

func NewPageContent(url string, title string, desc string, icon []byte) PageContent {
	content := new(WebPageContent)
	content.InitWithURL(url, title, desc, icon)
	return content
}

func (content *WebPageContent) Init(dict map[string]interface{}) PageContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._icon = nil
	}
	return content
}

func (content *WebPageContent) InitWithURL(url string, title string, desc string, icon []byte) PageContent {
	if content.BaseContent.InitWithType(PAGE) != nil {
		content.SetURL(url)
		content.SetTitle(title)
		content.SetDescription(desc)
		content.SetIcon(icon)
	}
	return content
}

//-------- IPageContent

func (content *WebPageContent) URL() string {
	text := content.Get("URL")
	if text == nil {
		return ""
	}
	return text.(string)
}
func (content *WebPageContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *WebPageContent) Title() string {
	text := content.Get("title")
	if text == nil {
		return ""
	}
	return text.(string)
}
func (content *WebPageContent) SetTitle(title string) {
	content.Set("title", title)
}

func (content *WebPageContent) Description() string {
	text := content.Get("desc")
	if text == nil {
		return ""
	}
	return text.(string)
}
func (content *WebPageContent) SetDescription(desc string) {
	content.Set("desc", desc)
}

func (content *WebPageContent) Icon() []byte {
	if content._icon == nil {
		base64 := content.Get("icon")
		if base64 != nil {
			content._icon = Base64Decode(base64.(string))
		}
	}
	return content._icon
}
func (content *WebPageContent) SetIcon(icon []byte) {
	if ValueIsNil(icon) {
		content.Remove("icon")
	} else {
		base64 := Base64Encode(icon)
		content.Set("icon", base64)
	}
	content._icon = icon
}
