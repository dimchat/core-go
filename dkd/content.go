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
	content := new(SecretContent).InitWithMessage(msg)
	ObjectRetain(content)
	return content
}

func (content *SecretContent) Init(dict map[string]interface{}) *SecretContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content.setSecretMessage(nil)
	}
	return content
}

func (content *SecretContent) InitWithMessage(msg ReliableMessage) *SecretContent {
	if content.BaseContent.InitWithType(FORWARD) != nil {
		content.Set("forward", msg.GetMap(false))
		content.setSecretMessage(msg)
	}
	return content
}

func (content *SecretContent) Release() int {
	cnt := content.Dictionary.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		content.setSecretMessage(nil)
	}
	return cnt
}

func (content *SecretContent) setSecretMessage(secret ReliableMessage) {
	if secret != content._secret {
		ObjectRetain(secret)
		ObjectRelease(content._secret)
		content._secret = secret
	}
}

//-------- IForwardContent

func (content *SecretContent) ForwardMessage() ReliableMessage {
	if content._secret == nil {
		forward := content.Get("forward")
		content.setSecretMessage(ReliableMessageParse(forward))
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
	content := new(BaseTextContent).InitWithText(text)
	ObjectRetain(content)
	return content
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
	text := content.Get("text")
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
	IPageContent

	_icon []byte
}

func NewPageContent(url string, title string, desc string, icon []byte) PageContent {
	content := new(WebPageContent).InitWithURL(url, title, desc, icon)
	ObjectRetain(content)
	return content
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
	url := content.Get("URL")
	return url.(string)
}
func (content *WebPageContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *WebPageContent) Title() string {
	title := content.Get("title")
	if title == nil {
		return ""
	}
	return title.(string)
}
func (content *WebPageContent) SetTitle(title string) {
	content.Set("title", title)
}

func (content *WebPageContent) Description() string {
	desc := content.Get("desc")
	if desc == nil {
		return ""
	}
	return desc.(string)
}
func (content *WebPageContent) SetDescription(desc string) {
	content.Set("desc", desc)
}

func (content *WebPageContent) Icon() []byte {
	if content._icon == nil {
		b64 := content.Get("icon")
		if b64 != nil {
			content._icon = Base64Decode(b64.(string))
		}
	}
	return content._icon
}
func (content *WebPageContent) SetIcon(icon []byte) {
	if icon == nil {
		content.Remove("icon")
	} else {
		b64 := Base64Encode(icon)
		content.Set("icon", b64)
	}
	content._icon = icon
}
