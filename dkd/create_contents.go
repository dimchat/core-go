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
	return NewBaseContent(dict, "")
}

/**
 *  Text Content
 */

func NewTextContent(text string) TextContent {
	return NewBaseTextContent(nil, text)
}

func NewTextContentWithMap(dict StringKeyMap) Content {
	return NewBaseTextContent(dict, "")
}

/**
 *  Name Card Content
 */

func NewNameCard(did ID, name string, avatar TransportableFile) NameCard {
	return NewNameCardContent(nil, did, name, avatar)
}

func NewNameCardWithMap(dict StringKeyMap) Content {
	return NewNameCardContent(dict, nil, "", nil)
}

/**
 *  Web Page Content
 */

func NewPageContentWithURL(url URL, title string, icon TransportableFile, desc string) PageContent {
	return NewWebPageContent(nil, title, icon, desc, url, "")
}
func NewPageContentWithHTML(html string, title string, icon TransportableFile, desc string) PageContent {
	return NewWebPageContent(nil, title, icon, desc, nil, html)
}

func NewPageContentWithMap(dict StringKeyMap) Content {
	return NewWebPageContent(dict, "", nil, "", nil, "")
}

/**
 *  Top-Secret Content
 */

func NewForwardMessage(msg ReliableMessage) ForwardContent {
	messages := []ReliableMessage{msg}
	return NewSecretContent(nil, messages)
}

func NewForwardMessages(messages []ReliableMessage) ForwardContent {
	return NewSecretContent(nil, messages)
}

func NewForwardContentWithMap(dict StringKeyMap) Content {
	return NewSecretContent(dict, nil)
}

/**
 *  Combine Forward Content
 */

func NewCombineMessages(title string, messages []InstantMessage) CombineContent {
	return NewCombineForwardContent(nil, title, messages)
}

func NewCombineContentWithMap(dict StringKeyMap) Content {
	return NewCombineForwardContent(dict, "", nil)
}

/**
 *  Array Content
 */

func NewArrayContent(contents []Content) ArrayContent {
	return NewListContent(nil, contents)
}

func NewArrayContentWithMap(dict StringKeyMap) Content {
	return NewListContent(dict, nil)
}

/**
 *  File Contents
 */

func NewFileContent(data TransportableData, filename string, url URL, password DecryptKey) FileContent {
	return NewBaseFileContent(nil, "", data, filename, url, password)
}
func NewFileContentWithData(data TransportableData, filename string) FileContent {
	return NewBaseFileContent(nil, "", data, filename, nil, nil)
}
func NewFileContentWithURL(url URL, password DecryptKey) FileContent {
	return NewBaseFileContent(nil, "", nil, "", url, password)
}

func NewFileContentWithMap(dict StringKeyMap) Content {
	return NewBaseFileContent(dict, "", nil, "", nil, nil)
}

/**
 *  Image Content
 */

func NewImageContent(data TransportableData, filename string, url URL, password DecryptKey) ImageContent {
	return NewImageFileContent(nil, data, filename, url, password, nil)
}
func NewImageContentWithData(data TransportableData, filename string) ImageContent {
	return NewImageFileContent(nil, data, filename, nil, nil, nil)
}
func NewImageContentWithURL(url URL, password DecryptKey) ImageContent {
	return NewImageFileContent(nil, nil, "", url, password, nil)
}

func NewImageContentWithMap(dict StringKeyMap) Content {
	return NewImageFileContent(dict, nil, "", nil, nil, nil)
}

/**
 *  Audio Content
 */

func NewAudioContent(data TransportableData, filename string, url URL, password DecryptKey) AudioContent {
	return NewAudioFileContent(nil, data, filename, url, password)
}
func NewAudioContentWithData(data TransportableData, filename string) AudioContent {
	return NewAudioFileContent(nil, data, filename, nil, nil)
}
func NewAudioContentWithURL(url URL, password DecryptKey) AudioContent {
	return NewAudioFileContent(nil, nil, "", url, password)
}

func NewAudioContentWithMap(dict StringKeyMap) Content {
	return NewAudioFileContent(dict, nil, "", nil, nil)
}

/**
 *  Video Content
 */

func NewVideoContent(data TransportableData, filename string, url URL, password DecryptKey) VideoContent {
	return NewVideoFileContent(nil, data, filename, url, password, nil)
}
func NewVideoContentWithData(data TransportableData, filename string) VideoContent {
	return NewVideoFileContent(nil, data, filename, nil, nil, nil)
}
func NewVideoContentWithURL(url URL, password DecryptKey) VideoContent {
	return NewVideoFileContent(nil, nil, "", url, password, nil)
}

func NewVideoContentWithMap(dict StringKeyMap) Content {
	return NewVideoFileContent(dict, nil, "", nil, nil, nil)
}

/**
 *  Money Contents
 */

func NewMoneyContent(currency string, amount float64) MoneyContent {
	return NewBaseMoneyContent(nil, "", currency, amount)
}

func NewMoneyContentWithMap(dict StringKeyMap) Content {
	return NewBaseMoneyContent(dict, "", "", 0)
}

func NewTransferContent(currency string, amount float64) TransferContent {
	return NewTransferMoneyContent(nil, currency, amount)
}

func NewTransferContentWithMap(dict StringKeyMap) Content {
	return NewTransferMoneyContent(dict, "", 0)
}

/**
 *  Quote Content
 */

func NewQuoteContent(text string, head Envelope, body Content) QuoteContent {
	origin := PurifyForQuote(head, body)
	return NewBaseQuoteContent(nil, text, origin)
}

func NewQuoteContentWithMap(dict StringKeyMap) Content {
	return NewBaseQuoteContent(dict, "", nil)
}

/**
 *  Application Customized Content
 */

func NewCustomizedContent(app, mod, act string) CustomizedContent {
	return NewAppCustomizedContent(nil, "", app, mod, act)
}

func NewCustomizedContentWithMap(dict StringKeyMap) Content {
	return NewAppCustomizedContent(dict, "", "", "", "")
}
