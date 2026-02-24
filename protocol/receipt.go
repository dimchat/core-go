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

// ReceiptCommand defines the interface for message receipt/acknowledgment commands
//
// Extends the Command interface for confirming message delivery/receipt
// Used to notify senders that their messages have been received/processed
//
//	Data structure: {
//	    "type"    : i2s(0x88),
//	    "sn"      : 456,
//
//	    "command" : "receipt",  // Fixed command name: "receipt"
//	    "text"    : "...",      // Optional receipt message/comment
//	    "origin"  : {           // Envelope of the original message being acknowledged
//	        "sender"    : "...",
//	        "receiver"  : "...",
//	        "time"      : 0,
//
//	        "sn"        : 123,
//	        "signature" : "..."
//	    }
//	}
type ReceiptCommand interface {
	Command

	// Text returns the optional receipt message/comment (acknowledgment note)
	//
	// Typical values: "Message received", "Read", "Processed", etc.
	Text() string

	// OriginalEnvelope returns the envelope of the original message being acknowledged
	OriginalEnvelope() Envelope

	// OriginalSerialNumber returns the serial number (SN) of the original message
	//
	// Extracts the "sn" field from the "origin" envelope for quick reference
	OriginalSerialNumber() SerialNumberType

	// OriginalSignature returns the signature of the original message
	//
	// Used to verify the authenticity of the message being acknowledged
	OriginalSignature() string
}

func PurifyForReceipt(head Envelope, body Content) StringKeyMap {
	helper := GetQuoteHelper()
	return helper.PurifyForReceipt(head, body)
}
