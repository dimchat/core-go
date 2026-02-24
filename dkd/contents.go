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
 */

type BaseTextContent struct {
	//TextContent
	*BaseContent
}

func NewBaseTextContent(dict StringKeyMap, text string) *BaseTextContent {
	if dict != nil {
		// init text content with map
		return &BaseTextContent{
			BaseContent: NewBaseContent(dict, ""),
		}
	}
	// new text content
	content := &BaseTextContent{
		BaseContent: NewBaseContent(nil, ContentType.TEXT),
	}
	content.Set("text", text)
	return content
}

// Override
func (content *BaseTextContent) Text() string {
	return content.GetString("text", "")
}

/**
 *  Web Page Content
 */

type WebPageContent struct {
	//PageContent
	*BaseContent

	// small image
	icon TransportableFile

	// web URL
	url URL
}

func NewWebPageContent(dict StringKeyMap,
	title string, icon TransportableFile, desc string,
	url URL, html string,
) *WebPageContent {
	if dict != nil {
		// init page content with map
		return &WebPageContent{
			BaseContent: NewBaseContent(dict, ""),
			// lazy load
			icon: nil,
			url:  nil,
		}
	}
	// new page content
	content := &WebPageContent{
		BaseContent: NewBaseContent(nil, ContentType.PAGE),
		icon:        icon,
		url:         url,
	}
	content.Set("title", title)
	//content.Set("icon", icon.Serialize())
	content.Set("desc", desc)
	if url != nil {
		content.Set("URL", url.String())
	}
	if html != "" {
		content.Set("HTML", html)
	}
	return content
}

// Override
func (content *WebPageContent) Map() StringKeyMap {
	// serialize 'icon'
	img := content.icon
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
	img := content.icon
	if img == nil {
		url := content.Get("icon")
		img = ParseTransportableFile(url)
		content.icon = img
	}
	return img
}

// Override
func (content *WebPageContent) SetIcon(img TransportableFile) {
	content.Remove("icon")
	//content.SetMapper("icon", img)
	content.icon = img
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
	url := content.url
	if url == nil {
		text := content.GetString("URL", "")
		if text != "" {
			url = ParseURL(text)
			content.url = url
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
	return content.GetString("HTML", "")
}

// Override
func (content *WebPageContent) SetHTML(html string) {
	content.Set("HTML", html)
}

/**
 *  Name Card Content
 */

type NameCardContent struct {
	//NameCard
	*BaseContent

	image TransportableFile
}

func NewNameCardContent(dict StringKeyMap, did ID, name string, avatar TransportableFile) *NameCardContent {
	if dict != nil {
		// init name card with map
		return &NameCardContent{
			BaseContent: NewBaseContent(dict, ""),
			// lazy load
			image: nil,
		}
	}
	// new name card
	content := &NameCardContent{
		BaseContent: NewBaseContent(dict, ContentType.NAME_CARD),
		image:       avatar,
	}
	// ID
	content.Set("did", did.String())
	// name
	content.Set("name", name)
	//// avatar
	//content.Set("avatar", avatar.Serialize()) // lazy serialize
	return content
}

// Override
func (content *NameCardContent) Map() StringKeyMap {
	// serialize 'avatar'
	img := content.image
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
	img := content.image
	if img == nil {
		url := content.Get("avatar")
		img = ParseTransportableFile(url)
		content.image = img
	}
	return img
}
