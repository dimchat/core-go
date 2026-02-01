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
	"fmt"

	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Receipt Command Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x88),
 *      "sn"   : 456,
 *
 *      "command" : "receipt",
 *      "text"    : "...",  // text message
 *      "origin"  : {       // original message envelope
 *          "sender"    : "...",
 *          "receiver"  : "...",
 *          "time"      : 0,
 *
 *          "sn"        : 123,
 *          "signature" : "..."
 *      }
 *  }
 *  </pre></blockquote>
 */
type BaseReceiptCommand struct {
	//ReceiptCommand
	BaseCommand

	/**
	 *  original message envelope
	 */
	_envelope Envelope
}

func (content *BaseReceiptCommand) InitWithMap(dict StringKeyMap) ReceiptCommand {
	if content.BaseCommand.InitWithMap(dict) != nil {
		// lazy load
		content._envelope = nil
	}
	return content
}

func (content *BaseReceiptCommand) Init(text string, origin StringKeyMap) ReceiptCommand {
	if content.BaseCommand.Init(RECEIPT) != nil {
		// text message
		content.Set("text", text)
		// original envelope of message responding to,
		// includes 'sn' and 'signature'
		if origin != nil {
			content.Set("origin", origin)
		}
		// lazy load
		content._envelope = nil
	}
	return content
}

// Override
func (content *BaseReceiptCommand) Text() string {
	return content.GetString("text", "")
}

// protected
func (content *BaseReceiptCommand) Origin() StringKeyMap {
	origin := content.Get("origin")
	if origin == nil {
		return nil
	} else if dict, ok := origin.(StringKeyMap); ok {
		return dict
	}
	panic(fmt.Sprintf("Invalid origin: %v", origin))
	return nil
}

// Override
func (content *BaseReceiptCommand) OriginalEnvelope() Envelope {
	env := content._envelope
	if env == nil {
		origin := content.Origin()
		env = ParseEnvelope(origin)
		content._envelope = env
	}
	return env
}

// Override
func (content *BaseReceiptCommand) OriginalSerialNumber() SerialNumberType {
	origin := content.Origin()
	if origin == nil {
		return 0
	}
	sn := origin["sn"]
	return ConvertUInt64(sn, 0)
}

// Override
func (content *BaseReceiptCommand) OriginalSignature() string {
	origin := content.Origin()
	if origin == nil {
		return ""
	}
	signature := origin["signature"]
	return ConvertString(signature, "")
}
