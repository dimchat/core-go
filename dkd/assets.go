/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

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
type BaseMoneyContent struct {
	//MoneyContent
	*BaseContent
}

func NewBaseMoneyContent(dict StringKeyMap, msgType MessageType, currency string, amount float64) *BaseMoneyContent {
	if dict != nil {
		// init money content with map
		return &BaseMoneyContent{
			BaseContent: NewBaseContent(dict, ""),
		}
	}
	// new money content
	if msgType == "" {
		msgType = ContentType.MONEY
	}
	content := &BaseMoneyContent{
		BaseContent: NewBaseContent(nil, msgType),
	}
	content.Set("currency", currency)
	content.Set("amount", amount)
	return content
}

// Override
func (content *BaseMoneyContent) Currency() string {
	return content.GetString("currency", "")
}

// Override
func (content *BaseMoneyContent) Amount() float64 {
	return content.GetFloat64("amount", 0)
}

// Override
func (content *BaseMoneyContent) SetAmount(amount float64) {
	content.Set("amount", amount)
}

/**
 *  Transfer Money
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x41),
 *      "sn"   : 123,
 *
 *      "currency" : "RMB",    // USD, USDT, ...
 *      "amount"   : 100.00,
 *      "remitter" : "{FROM}", // sender ID
 *      "remittee" : "{TO}"    // receiver ID
 *  }
 *  </pre></blockquote>
 */
type TransferMoneyContent struct {
	//TransferContent
	*BaseMoneyContent
}

func NewTransferMoneyContent(dict StringKeyMap, currency string, amount float64) *TransferMoneyContent {
	return &TransferMoneyContent{
		BaseMoneyContent: NewBaseMoneyContent(dict, ContentType.TRANSFER, currency, amount),
	}
}

// Override
func (content *TransferMoneyContent) Remitter() ID {
	sender := content.Get("remitter")
	return ParseID(sender)
}

// Override
func (content *TransferMoneyContent) SetRemitter(sender ID) {
	content.SetStringer("remitter", sender)
}

// Override
func (content *TransferMoneyContent) Remittee() ID {
	receiver := content.Get("remittee")
	return ParseID(receiver)
}

// Override
func (content *TransferMoneyContent) SetRemittee(receiver ID) {
	content.SetStringer("remittee", receiver)
}
