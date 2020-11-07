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
	. "github.com/dimchat/mkm-go/format"
)

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
type PageContent struct {
	BaseContent

	_icon []byte
}

func (content *PageContent) Init(dictionary map[string]interface{}) *PageContent {
	if content.BaseContent.Init(dictionary) != nil {
		// lazy load
		content._icon = nil
	}
	return content
}

func (content *PageContent) InitWithURL(url string, title string, desc string, icon []byte) *PageContent {
	if content.BaseContent.InitWithType(PAGE) != nil {
		content.SetURL(url)
		content.SetTitle(title)
		content.SetDescription(desc)
		content.SetIcon(icon)
	}
	return content
}

//-------- setter/getter --------

func (content *PageContent) GetURL() string {
	url := content.Get("URL")
	return url.(string)
}

func (content *PageContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *PageContent) GetTitle() string {
	title := content.Get("title")
	return title.(string)
}

func (content *PageContent) SetTitle(title string) {
	content.Set("title", title)
}

func (content *PageContent) GetDescription() string {
	desc := content.Get("desc")
	if desc == nil {
		return ""
	} else {
		return desc.(string)
	}
}

func (content *PageContent) SetDescription(desc string) {
	content.Set("desc", desc)
}

func (content *PageContent) GetIcon() []byte {
	if content._icon == nil {
		b64 := content.Get("icon")
		if b64 != nil {
			content._icon = Base64Decode(b64.(string))
		}
	}
	return content._icon
}

func (content *PageContent) SetIcon(icon []byte) {
	if icon == nil {
		content.Set("icon", nil)
	} else {
		b64 := Base64Encode(icon)
		content.Set("icon", b64)
	}
	content._icon = icon
}
