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
	//-------- group command names begin --------
	// founder/owner
	FOUND    = "found"
	ABDICATE = "abdicate"
	// member
	INVITE   = "invite"
	EXPEL    = "expel"
	JOIN     = "join"
	QUIT     = "quit"
	QUERY    = "query"
	RESET    = "reset"
	// administrator/assistant
	HIRE     = "hire"
	FIRE     = "fire"
	RESIGN   = "resign"
	//-------- group command names end --------
)

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "{NAME}",      // join, quit, ...
 *      group   : "{GROUP_ID}",
 *      // extra info: member or members
 *  }
 */
type GroupCommand interface {
	HistoryCommand
	IGroupCommand
}
type IGroupCommand interface {

	Member() ID
	SetMember(member ID)

	Members() []ID
	SetMembers(members []ID)
}

//-------- Group Commands

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "invite",
 *      group   : "{GROUP_ID}",
 *      members : [],            // member ID list
 *  }
 */
type InviteCommand interface {
	GroupCommand
	IInviteCommand
}
type IInviteCommand interface {

	InviteMembers() []ID
}

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "expel",
 *      group   : "{GROUP_ID}",
 *      members : [],            // member ID list
 *  }
 */
type ExpelCommand interface {
	GroupCommand
	IExpelCommand
}
type IExpelCommand interface {

	ExpelMembers() []ID
}

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "join",
 *      group   : "{GROUP_ID}",
 *      text    : "May I?",
 *  }
 */
type JoinCommand interface {
	GroupCommand
	IJoinCommand
}
type IJoinCommand interface {

	Ask() string
}

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "quit",
 *      group   : "{GROUP_ID}",
 *      text    : "Good bye!",
 *  }
 */
type QuitCommand interface {
	GroupCommand
	IQuitCommand
}
type IQuitCommand interface {

	Bye() string
}

/**
 *  Group history command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "reset",
 *      group   : "{GROUP_ID}",
 *      members : [],            // member ID list
 *  }
 */
type ResetCommand interface {
	GroupCommand
	IResetCommand
}
type IResetCommand interface {

	AllMembers() []ID
}

/**
 *  Group command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "query",
 *      group   : "{GROUP_ID}",
 *      text    : "May I?",
 *  }
 */
type QueryCommand interface {
	GroupCommand
	IQueryCommand
}
//  NOTICE:
//      This command is just for querying group info,
//      should not be saved in group history
type IQueryCommand interface {

	Text() string
}
