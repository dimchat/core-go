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

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// MetaCommand defines the interface for meta information exchange commands
//
// Extends the base Command interface for querying/updating entity meta information
//
//	Data structure: {
//	    "type"    : i2s(0x88),
//	    "sn"      : 123,
//
//	    "command" : "meta",  // Fixed command name: "meta"
//	    "did"     : "{ID}",  // Target entity ID (contact's ID)
//	    "meta"    : {...}    // Meta info (nil = query meta for the specified ID)
//	}
type MetaCommand interface {
	Command

	// ID returns the target entity ID for the meta command (contact's ID)
	ID() ID

	// Meta returns the meta information associated with the command
	//
	// Returns nil if the command is a meta query (rather than an update)
	Meta() Meta
}

// DocumentCommand defines the interface for document/profile exchange commands
//
// Extends MetaCommand for querying/updating entity documents (profiles/TAI)
//
//	Data structure: {
//	    "type"      : i2s(0x88),
//	    "sn"        : 123,
//
//	    "command"   : "documents",  // Fixed command name: "documents"
//	    "did"       : "{ID}",       // Target entity ID
//	    "meta"      : {...},        // Optional meta info (only for new friend handshake)
//	    "documents" : [...],        // List of entity documents (nil = query documents)
//	    "last_time" : 12345         // Timestamp of last document (for incremental query)
//	}
type DocumentCommand interface {
	MetaCommand

	// Documents returns the list of entity documents/profiles in the command
	//
	// Returns nil if the command is a document query (rather than an update)
	Documents() []Document

	// LastTime returns the timestamp of the last document from the sender
	//
	// Used for incremental document queries (only get documents newer than this time)
	LastTime() Time
}
