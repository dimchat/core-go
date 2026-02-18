/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

func NewContentWithMap(dict StringKeyMap) Content {
	content := &BaseContent{}
	return content.InitWithMap(dict)
}

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
func NewTextContent(text string) TextContent {
	content := &BaseTextContent{}
	return content.Init(text)
}

func NewTextContentWithMap(dict StringKeyMap) TextContent {
	content := &BaseTextContent{}
	content.InitWithMap(dict)
	return content
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
func NewNameCard(did ID, name string, avatar TransportableFile) NameCard {
	content := &NameCardContent{}
	return content.Init(did, name, avatar)
}

func NewNameCardWithMap(dict StringKeyMap) NameCard {
	content := &NameCardContent{}
	return content.InitWithMap(dict)
}

/**
 *  Web Page
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
func NewPageContentWithURL(url URL, title string, icon TransportableFile, desc string) PageContent {
	content := &WebPageContent{}
	return content.Init(title, icon, desc, url, "")
}
func NewPageContentWithHTML(html string, title string, icon TransportableFile, desc string) PageContent {
	content := &WebPageContent{}
	return content.Init(title, icon, desc, nil, html)
}

func NewPageContentWithMap(dict StringKeyMap) PageContent {
	content := &WebPageContent{}
	return content.InitWithMap(dict)
}

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
func NewForwardMessage(msg ReliableMessage) ForwardContent {
	content := &SecretContent{}
	return content.InitWithMessage(msg)
}

func NewForwardMessages(messages []ReliableMessage) ForwardContent {
	content := &SecretContent{}
	return content.InitWithMessages(messages)
}

func NewForwardContentWithMap(dict StringKeyMap) ForwardContent {
	content := &SecretContent{}
	return content.InitWithMap(dict)
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
func NewCombineMessages(title string, messages []InstantMessage) CombineContent {
	content := &CombineForwardContent{}
	return content.Init(title, messages)
}

func NewCombineContentWithMap(dict StringKeyMap) CombineContent {
	content := &CombineForwardContent{}
	return content.InitWithMap(dict)
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
func NewArrayContent(contents []Content) ArrayContent {
	content := &ListContent{}
	return content.Init(contents)
}

func NewArrayContentWithMap(dict StringKeyMap) ArrayContent {
	content := &ListContent{}
	return content.InitWithMap(dict)
}

/**
 *  File Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x10),
 *      "sn"   : 123,
 *
 *      "data"     : "...",        // base64_encode(fileContent)
 *      "filename" : "photo.png",
 *
 *      "URL"      : "http://...", // download from CDN
 *      // before fileContent uploaded to a public CDN,
 *      // it should be encrypted by a symmetric key
 *      "key"      : {             // symmetric key to decrypt file data
 *          "algorithm" : "AES",   // "DES", ...
 *          "data"      : "{BASE64_ENCODE}",
 *          ...
 *      }
 *  }
 *  </pre></blockquote>
 */
func NewFileContent(data TransportableData, filename string, url URL, key DecryptKey) FileContent {
	content := &BaseFileContent{}
	return content.InitWithType(ContentType.FILE, data, filename, url, key)
}
func NewFileContentWithData(data TransportableData, filename string) FileContent {
	return NewFileContent(data, filename, nil, nil)
}
func NewFileContentWithURL(url URL, key DecryptKey) FileContent {
	return NewFileContent(nil, "", url, key)
}

func NewFileContentWithMap(dict StringKeyMap) FileContent {
	content := &BaseFileContent{}
	return content.InitWithMap(dict)
}

/**
 *  Image Content
 */
func NewImageContent(data TransportableData, filename string, url URL, key DecryptKey) ImageContent {
	content := &ImageFileContent{}
	return content.Init(data, filename, url, key)
}
func NewImageContentWithData(data TransportableData, filename string) ImageContent {
	return NewImageContent(data, filename, nil, nil)
}
func NewImageContentWithURL(url URL, key DecryptKey) ImageContent {
	return NewImageContent(nil, "", url, key)
}

func NewImageContentWithMap(dict StringKeyMap) ImageContent {
	content := &ImageFileContent{}
	return content.InitWithMap(dict)
}

/**
 *  Audio Content
 */
func NewAudioContent(data TransportableData, filename string, url URL, key DecryptKey) AudioContent {
	content := &AudioFileContent{}
	return content.Init(data, filename, url, key)
}
func NewAudioContentWithData(data TransportableData, filename string) AudioContent {
	return NewAudioContent(data, filename, nil, nil)
}
func NewAudioContentWithURL(url URL, key DecryptKey) AudioContent {
	return NewAudioContent(nil, "", url, key)
}

func NewAudioContentWithMap(dict StringKeyMap) AudioContent {
	content := &AudioFileContent{}
	content.InitWithMap(dict)
	return content
}

/**
 *  Video Content
 */
func NewVideoContent(data TransportableData, filename string, url URL, key DecryptKey) VideoContent {
	content := &VideoFileContent{}
	return content.Init(data, filename, url, key)
}
func NewVideoContentWithData(data TransportableData, filename string) VideoContent {
	return NewVideoContent(data, filename, nil, nil)
}
func NewVideoContentWithURL(url URL, key DecryptKey) VideoContent {
	return NewVideoContent(nil, "", url, key)
}

func NewVideoContentWithMap(dict StringKeyMap) VideoContent {
	content := &VideoFileContent{}
	return content.InitWithMap(dict)
}

/**
 *  Money Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x40),
 *      "sn"   : 123,
 *
 *      "currency" : "RMB", // USD, USDT, ...
 *      "amount"   : 100.00
 *  }
 *  </pre></blockquote>
 */
func NewMoneyContent(currency string, amount float64) MoneyContent {
	content := &BaseMoneyContent{}
	return content.Init(currency, amount)
}

func NewMoneyContentWithMap(dict StringKeyMap) MoneyContent {
	content := &BaseMoneyContent{}
	content.InitWithMap(dict)
	return content
}

func NewTransferContent(currency string, amount float64) TransferContent {
	content := &TransferMoneyContent{}
	return content.Init(currency, amount)
}

func NewTransferContentWithMap(dict StringKeyMap) TransferContent {
	content := &TransferMoneyContent{}
	content.InitWithMap(dict)
	return content
}

/**
 *  Quote Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x37),
 *      "sn"   : 456,
 *
 *      "text"   : "...",  // text message
 *      "origin" : {       // original message envelope
 *          "sender"   : "...",
 *          "receiver" : "...",
 *
 *          "type"     : i2s(0x01),
 *          "sn"       : 123,
 *      }
 *  }
 *  </pre></blockquote>
 */
func NewQuoteContent(text string, head Envelope, body Content) QuoteContent {
	origin := PurifyForQuote(head, body)
	content := &BaseQuoteContent{}
	return content.Init(text, origin)
}

func NewQuoteContentWithMap(dict StringKeyMap) QuoteContent {
	content := &BaseQuoteContent{}
	return content.InitWithMap(dict)
}

/**
 *  Application Customized message
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0xCC),
 *      "sn"   : 123,
 *
 *      "app"   : "{APP_ID}",  // application (e.g.: "chat.dim.sechat")
 *      "mod"   : "{MODULE}",  // module name (e.g.: "drift_bottle")
 *      "act"   : "{ACTION}",  // action name (3.g.: "throw")
 *      "extra" : info         // action parameters
 *  }
 *  </pre></blockquote>
 */
func NewCustomizedContent(app, mod, act string) CustomizedContent {
	content := &AppCustomizedContent{}
	return content.Init(app, mod, act)
}

func NewCustomizedContentWithMap(dict StringKeyMap) CustomizedContent {
	content := &AppCustomizedContent{}
	content.InitWithMap(dict)
	return content
}
