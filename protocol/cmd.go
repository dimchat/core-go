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
)

const (
	//-------- history command names begin --------
	// account
	REGISTER = "register"
	SUICIDE  = "suicide"
	//-------- history command names end --------
)

/**
 *  History command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "...", // command name
 *      time    : 0,     // command timestamp
 *      // extra info
 *  }
 */
type HistoryCommand interface {
	Command
	IHistoryCommand
}
type IHistoryCommand interface {

}

/**
 *  Meta command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "meta", // command name
 *      ID      : "{ID}", // contact's ID
 *      meta    : {...}   // when meta is empty, means query meta for ID
 *  }
 */
type MetaCommand interface {
	Command
	IMetaCommand
}
type IMetaCommand interface {

	ID() ID
	Meta() Meta
}

/**
 *  Document Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command   : "document",  // command name
 *      ID        : "{ID}",      // entity ID
 *      meta      : {...},       // only for handshaking with new friend
 *      profile   : {...},       // when profile is empty, means query for ID
 *      signature : "..."        // old profile's signature for querying
 *  }
 */
type DocumentCommand interface {
	MetaCommand
	IDocumentCommand
}
type IDocumentCommand interface {

	Document() Document
	Signature() string
}
