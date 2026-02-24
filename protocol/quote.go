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

// QuoteContent defines the interface for quoted/reply message content
//
// Extends the base Content interface for replying to a specific message (with context)
//
//	Data structure: {
//	    "type"   : i2s(0x37),
//	    "sn"     : 456,
//
//	    "text"   : "...",  // Reply text/message
//	    "origin" : {       // Envelope of the original quoted message
//	        "sender"   : "...",
//	        "receiver" : "...",
//	        "type"     : i2s(0x01),  // Message type of the original content
//	        "sn"       : 123,        // Serial number of the original message
//	    }
//	}
type QuoteContent interface {
	Content

	// Text returns the reply text/message for the quote
	Text() string

	// OriginalEnvelope returns the envelope of the original quoted message
	//
	// Contains sender/receiver/type/sn of the message being replied to
	OriginalEnvelope() Envelope
	OriginalSerialNumber() SerialNumberType
}

func PurifyForQuote(head Envelope, body Content) StringKeyMap {
	helper := GetQuoteHelper()
	return helper.PurifyForQuote(head, body)
}
