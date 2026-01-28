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
 *  History Command
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x89),
 *      "sn"   : 123,
 *
 *      "command" : "...", // command name
 *      "time"    : 0,     // command timestamp
 *      "extra"   : info   // command parameters
 *  }
 *  </pre></blockquote>
 */
type HistoryCommand interface {
	Command
}

const (
	//-------- group command names begin --------
	// founder/owner
	FOUND    = "found"
	ABDICATE = "abdicate"
	// member
	INVITE = "invite"
	EXPEL  = "expel"
	JOIN   = "join"
	QUIT   = "quit"
	//QUERY  = "query"
	RESET = "reset"
	// administrator/assistant
	HIRE   = "hire"
	FIRE   = "fire"
	RESIGN = "resign"
	//-------- group command names end --------
)

/**
 *  Group History
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x89),
 *      "sn"   : 123,
 *
 *      "command" : "reset",   // "invite", "quit", ...
 *      "time"    : 123.456,   // command timestamp
 *
 *      "group"   : "{GROUP_ID}",
 *      "members" : ["{MEMBER_ID}",]
 *  }
 *  </pre></blockquote>
 */
type GroupCommand interface {
	HistoryCommand

	Members() []ID
	SetMembers(members []ID)
}

//-------- Group Commands

type InviteCommand interface {
	GroupCommand
}

// Deprecated (use 'reset' instead)
type ExpelCommand interface {
	GroupCommand
}

/**
 *  Group history command: {
 *      "type" : i2s(0x89),
 *      "sn"   : 123,
 *
 *      "command" : "join",
 *      "time"    : 123.456,
 *
 *      "group"   : "{GROUP_ID}",
 *      "text"    : "May I?",
 *  }
 */
type JoinCommand interface {
	GroupCommand

	Ask() string
}

/**
 *  Group history command: {
 *      "type" : i2s(0x89),
 *      "sn"   : 123,
 *
 *      "command" : "quit",
 *      "time"    : 123.456,
 *
 *      "group"   : "{GROUP_ID}",
 *      "text"    : "Good bye!",
 *  }
 */
type QuitCommand interface {
	GroupCommand

	Bye() string
}

type ResetCommand interface {
	GroupCommand
}
