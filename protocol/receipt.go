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
package protocol

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Receipt Command
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
type ReceiptCommand interface {
	Command

	Text() string

	OriginalEnvelope() Envelope
	OriginalSerialNumber() SerialNumberType
	OriginalSignature() string
}

func PurifyForReceipt(head Envelope, body Content) StringKeyMap {
	if head == nil {
		return nil
	}
	info := head.CopyMap(false)
	if _, exists := info["data"]; exists {
		delete(info, "data")
		delete(info, "key")
		delete(info, "keys")
		delete(info, "meta")
		delete(info, "visa")
	}
	if body != nil {
		sn := body.SN()
		info["sn"] = sn
	}
	return info
}
