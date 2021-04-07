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
	IForwardContent

	_secret ReliableMessage
}

func NewForwardContent(msg ReliableMessage) ForwardContent {
	return new(SecretContent).InitWithMessage(msg)
}

func (content *SecretContent) Init(dict map[string]interface{}) *SecretContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._secret = nil
	}
	return content
}

func (content *SecretContent) InitWithMessage(secret ReliableMessage) *SecretContent {
	if content.BaseContent.InitWithType(FORWARD) != nil {
		content.Set("forward", secret.GetMap(false))
		content._secret = secret
	}
	return content
}

//-------- IForwardContent

func (content *SecretContent) ForwardMessage() ReliableMessage {
	if content._secret == nil {
		forward := content.Get("forward")
		content._secret = ReliableMessageParse(forward)
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
	ITextContent
}

func NewTextContent(text string) TextContent {
	return new(BaseTextContent).InitWithText(text)
}

func (content *BaseTextContent) Init(dict map[string]interface{}) *BaseTextContent {
	if content.BaseContent.Init(dict) != nil {
	}
	return content
}

func (content *BaseTextContent) InitWithText(text string) *BaseTextContent {
	if content.InitWithType(TEXT) != nil {
		content.SetText(text)
	}
	return content
}

//-------- ITextContent

func (content *BaseTextContent) Text() string {
	text, ok := content.Get("text").(string)
	if ok {
		return text
	} else {
		return ""
	}
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
	IPageContent

	_icon []byte
}

func NewPageContent(url string, title string, desc string, icon []byte) PageContent {
	return new(WebPageContent).InitWithURL(url, title, desc, icon)
}

func (content *WebPageContent) Init(dict map[string]interface{}) *WebPageContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._icon = nil
	}
	return content
}

func (content *WebPageContent) InitWithURL(url string, title string, desc string, icon []byte) *WebPageContent {
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
	url, ok := content.Get("URL").(string)
	if ok {
		return url
	} else {
		return ""
	}
}
func (content *WebPageContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *WebPageContent) Title() string {
	text, ok := content.Get("title").(string)
	if ok {
		return text
	} else {
		return ""
	}
}
func (content *WebPageContent) SetTitle(title string) {
	content.Set("title", title)
}

func (content *WebPageContent) Description() string {
	text, ok := content.Get("desc").(string)
	if ok {
		return text
	} else {
		return ""
	}
}
func (content *WebPageContent) SetDescription(desc string) {
	content.Set("desc", desc)
}

func (content *WebPageContent) Icon() []byte {
	if content._icon == nil {
		b64, ok := content.Get("icon").(string)
		if ok {
			content._icon = Base64Decode(b64)
		}
	}
	return content._icon
}
func (content *WebPageContent) SetIcon(icon []byte) {
	if ValueIsNil(icon) {
		content.Remove("icon")
	} else {
		b64 := Base64Encode(icon)
		content.Set("icon", b64)
	}
	content._icon = icon
}
