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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

/**
 *  Money message: {
 *      type : 0x40,
 *      sn   : 123,
 *
 *      currency : "RMB", // USD, USDT, ...
 *      amount   : 100.00
 *  }
 */
type BaseMoneyContent struct {
	BaseContent
}

func NewMoneyContent(msgType ContentType, currency string, amount float64) MoneyContent {
	content := new(BaseMoneyContent)
	content.InitWithType(msgType, currency, amount)
	return content
}

///* designated initializer */
//func (content *BaseMoneyContent) Init(dict map[string]interface{}) MoneyContent {
//	if content.BaseContent.Init(dict) != nil {
//	}
//	return content
//}

/* designated initializer */
func (content *BaseMoneyContent) InitWithType(msgType ContentType, currency string, amount float64) MoneyContent {
	if msgType == 0 {
		msgType = MONEY
	}
	if content.BaseContent.InitWithType(msgType) != nil {
		content.setCurrency(currency)
		content.SetAmount(amount)
	}
	return content
}

func (content *BaseMoneyContent) setCurrency(currency string) {
	content.Set("currency", currency)
}

//-------- IMoneyContent

func (content *BaseMoneyContent) Currency() string {
	text := content.Get("currency")
	if text == nil {
		return ""
	}
	return text.(string)
}

func (content *BaseMoneyContent) Amount() float64 {
	amount := content.Get("amount")
	if amount == nil {
		return 0.0
	}
	return amount.(float64)
}
func (content *BaseMoneyContent) SetAmount(amount float64) {
	content.Set("amount", amount)
}

/**
 *  Transfer money message: {
 *      type : 0x41,
 *      sn   : 123,
 *
 *      currency : "RMB", // USD, USDT, ...
 *      amount   : 100.00
 *  }
 */
type TransferMoneyContent struct {
	BaseMoneyContent
}

func NewTransferContent(currency string, amount float64) TransferContent {
	content := new(TransferMoneyContent)
	content.InitWithCurrency(currency, amount)
	return content
}

//func (content *TransferMoneyContent) Init(dict map[string]interface{}) TransferContent {
//	if content.BaseMoneyContent.Init(dict) != nil {
//	}
//	return content
//}

func (content *TransferMoneyContent) InitWithCurrency(currency string, amount float64) TransferContent {
	if content.BaseMoneyContent.InitWithType(TRANSFER, currency, amount) != nil {
	}
	return content
}

//-------- ITransferContent

func (content *TransferMoneyContent) Comment() string {
	text := content.Get("text")
	if text == nil {
		return ""
	}
	return text.(string)
}
