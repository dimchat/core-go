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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Text Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x01),
 *      "sn"   : 123,
 *
 *      "text" : "..."
 *  }
 *  </pre></blockquote>
 */
type BaseTextContent struct {
	//TextContent
	BaseContent
}

func (content *BaseTextContent) Init(text string) TextContent {
	if content.InitWithType(ContentType.TEXT) != nil {
		content.Set("text", text)
	}
	return content
}

// Override
func (content *BaseTextContent) Text() string {
	return content.GetString("text", "")
}

/**
 *  Web Page Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x20),
 *      "sn"   : 123,
 *
 *      "title" : "...",                // Web title
 *      "desc"  : "...",
 *      "icon"  : "data:image/x-icon;base64,...",
 *
 *      "URL"   : "https://github.com/moky/dimp",
 *
 *      "HTML"      : "...",            // Web content
 *      "mime_type" : "text/html",      // Content-Type
 *      "encoding"  : "utf8",
 *      "base"      : "about:blank"     // Base URL
 *  }
 *  </pre></blockquote>
 */
type WebPageContent struct {
	//PageContent
	BaseContent

	// small image
	_icon TransportableFile

	// web URL
	_url URL
}

func (content *WebPageContent) InitWithMap(dict StringKeyMap) PageContent {
	if content.BaseContent.InitWithMap(dict) != nil {
		// lazy load
		content._icon = nil
		content._url = nil
	}
	return content
}

//func (content *WebPageContent) Init(
//	title string,
//	icon TransportableFile,
//	desc string,
//	url URL,
//	html string,
//) PageContent {
//	if content.BaseContent.InitWithType(ContentType.PAGE) != nil {
//
//		content.SetTitle(title)
//		content.SetIcon(icon)
//
//		content.SetDescription(desc)
//
//		content.SetURL(url)
//		content.SetHTML(html)
//	}
//	return content
//}

// Override
func (content *WebPageContent) Map() StringKeyMap {
	// serialize 'icon'
	img := content._icon
	if img != nil && !content.Contains("icon") {
		content.Set("icon", img.Serialize())
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *WebPageContent) Title() string {
	return content.GetString("title", "")
}

// Override
func (content *WebPageContent) SetTitle(title string) {
	content.Set("title", title)
}

// Override
func (content *WebPageContent) Icon() TransportableFile {
	img := content._icon
	if img == nil {
		url := content.Get("icon")
		img = ParseTransportableFile(url)
		content._icon = img
	}
	return img
}

// Override
func (content *WebPageContent) SetIcon(img TransportableFile) {
	content.Remove("icon")
	//content.SetMapper("icon", img)
	content._icon = img
}

// Override
func (content *WebPageContent) Description() string {
	return content.GetString("desc", "")
}

// Override
func (content *WebPageContent) SetDescription(desc string) {
	content.Set("desc", desc)
}

// Override
func (content *WebPageContent) URL() URL {
	url := content._url
	if url == nil {
		text := content.GetString("URL", "")
		if text != "" {
			url = ParseURL(text)
			content._url = url
		}
	}
	return url
}

// Override
func (content *WebPageContent) SetURL(url URL) {
	content.Set("URL", url.String())
}

// Override
func (content *WebPageContent) HTML() string {
	return content.GetString("html", "")
}

// Override
func (content *WebPageContent) SetHTML(html string) {
	content.Set("html", html)
}

/**
 *  Name Card
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x33),
 *      "sn"   : 123,
 *
 *      "did"    : "{ID}",        // contact's ID
 *      "name"   : "{nickname}",  // contact's name
 *      "avatar" : "{URL}",       // avatar - PNF(URL)
 *      ...
 *  }
 *  </pre></blockquote>
 */
type NameCardContent struct {
	//NameCard
	BaseContent

	_image TransportableFile
}

func (content *NameCardContent) InitWithMap(dict StringKeyMap) NameCard {
	if content.BaseContent.InitWithMap(dict) != nil {
		// lazy load
		content._image = nil
	}
	return content
}

func (content *NameCardContent) Init(did ID, name string, avatar TransportableFile) NameCard {
	if content.BaseContent.InitWithType(ContentType.NAME_CARD) != nil {
		// ID
		content.Set("did", did.String())
		// name
		content.Set("name", name)
		// avatar
		content._image = avatar // lazy serialize
	}
	return content
}

// Override
func (content *NameCardContent) Map() StringKeyMap {
	// serialize 'avatar'
	img := content._image
	if img != nil && !content.Contains("avatar") {
		content.Set("avatar", img.Serialize())
	}
	// OK
	return content.BaseContent.Map()
}

// Override
func (content *NameCardContent) ID() ID {
	did := content.Get("did")
	return ParseID(did)
}

// Override
func (content *NameCardContent) Name() string {
	return content.GetString("name", "")
}

// Override
func (content *NameCardContent) Avatar() TransportableFile {
	img := content._image
	if img == nil {
		url := content.Get("avatar")
		img = ParseTransportableFile(url)
		content._image = img
	}
	return img
}

//
//  Factories
//

func NewTextContent(text string) TextContent {
	content := &BaseTextContent{}
	return content.Init(text)
}

func NewNameCard(did ID, name string, avatar TransportableFile) NameCard {
	content := &NameCardContent{}
	return content.Init(did, name, avatar)
}

func NewPageContentWithURL(url URL, title string, icon TransportableFile, desc string) PageContent {
	content := &WebPageContent{}
	content.InitWithMap(NewMap())
	content.SetURL(url)
	content.SetTitle(title)
	content.SetIcon(icon)
	content.SetDescription(desc)
	return content
}

func NewPageContentWithHTML(html string, title string, icon TransportableFile, desc string) PageContent {
	content := &WebPageContent{}
	content.InitWithMap(NewMap())
	content.SetHTML(html)
	content.SetTitle(title)
	content.SetIcon(icon)
	content.SetDescription(desc)
	return content
}
