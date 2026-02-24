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
package protocol

import . "github.com/dimchat/dkd-go/protocol"

// ForwardContent defines the interface for top-secret forwarded message content
//
// Extends the base Content interface for transmitting secure/certified messages
// Used to forward one or more reliable (secure + certified) messages
//
//	Data structure: {
//	    "type"    : i2s(0xFF),
//	    "sn"      : 456,
//
//	    "forward" : {...},  // Single reliable (secure + certified) message (optional)
//	    "secrets" : [...]   // List of reliable (secure + certified) messages
//	}
type ForwardContent interface {
	Content

	// SecretMessages returns the list of reliable (secure + certified) messages
	//
	// Maps to the "secrets" field in the data structure (includes "forward" if present)
	SecretMessages() []ReliableMessage
}

// CombineContent defines the interface for combined/merged message content
//
// Extends the base Content interface for grouping multiple chat messages (chat history)
//
//	Data structure: {
//	    "type"     : i2s(0xCF),
//	    "sn"       : 123,
//
//	    "title"    : "...",  // Chat title/description for the combined messages
//	    "messages" : [...]   // List of InstantMessage instances (chat history)
//	}
type CombineContent interface {
	Content

	// Title returns the chat title/description for the combined messages
	Title() string

	// Messages returns the list of InstantMessage instances (chat history)
	Messages() []InstantMessage
}

// ArrayContent defines the interface for array-based message content
//
// Extends the base Content interface for wrapping a list of heterogeneous Content instances
//
//	Data structure: {
//	   "type"     : i2s(0xCA),
//	   "sn"       : 123,
//
//	   "contents" : [...]  // content array
//	}
type ArrayContent interface {
	Content

	// Contents returns the list of heterogeneous Content instances
	Contents() []Content
}
