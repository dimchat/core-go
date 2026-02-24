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

import . "github.com/dimchat/mkm-go/protocol"

// History command name constants for account operations
// These values are used as the "command" field in HistoryCommand messages
const (
	REGISTER = "register" // Command for account registration (create new account)
	SUICIDE  = "suicide"  // Command for account deletion (permanently remove account)
)

// HistoryCommand defines the base interface for historical operation command messages
//
// Extends the Command interface for recording timestamped system operations (account/group changes)
//
//	Standard data structure: {
//	    "type"    : i2s(0x89),
//	    "sn"      : 123,
//
//	    "command" : "...",  // History command name (e.g., "register", "found", "invite")
//	    "time"    : 0,      // Timestamp when the command was executed (float)
//	    "extra"   : info    // Optional command-specific parameters
//	}
type HistoryCommand interface {
	Command
}

// Group command name constants for group membership/role operations
// These values are used as the "command" field in GroupCommand messages
const (
	// Founder/Owner operations
	FOUND    = "found"    // Command to create/found a new group (initiate group)
	ABDICATE = "abdicate" // Command for group owner to transfer ownership (step down)

	// Group member operations
	INVITE = "invite" // Command to invite users to join the group
	EXPEL  = "expel"  // Command to remove members from the group (deprecated: use RESET)
	JOIN   = "join"   // Command for user to request joining the group
	QUIT   = "quit"   // Command for member to leave the group voluntarily
	//QUERY  = "query"    // Reserved for group membership query (Deprecated)
	RESET = "reset" // Command to reset group membership (replaces EXPEL)

	// Administrator operations
	HIRE   = "hire"   // Command to appoint a group administrator
	FIRE   = "fire"   // Command to remove a group administrator
	RESIGN = "resign" // Command for administrator to step down voluntarily
)

// GroupCommand defines the interface for group membership/role change history commands
//
// Extends HistoryCommand for recording group operations (invite, expel, hire, etc.)
//
//	Standard data structure: {
//	    "type"    : i2s(0x89),
//	    "sn"      : 123,
//
//	    "command" : "reset",  // Group command name (e.g., "invite", "quit", "hire")
//	    "time"    : 123.456,  // Timestamp when the group command was executed
//
//	    "group"   : "{GROUP_ID}",
//	    "members" : ["{MEMBER_ID}",]
//	}
type GroupCommand interface {
	HistoryCommand

	// Members returns the list of member IDs affected by the group command
	Members() []ID
	SetMembers(members []ID)
}

//-------- Group Commands

// InviteCommand defines the interface for group member invitation commands
//
// Extends GroupCommand for the "invite" group operation
type InviteCommand interface {
	GroupCommand
}

// ExpelCommand defines the interface for group member expulsion commands
//
// Extends GroupCommand for the "expel" group operation
//
// Deprecated: Use ResetCommand ("reset") instead for membership removal
type ExpelCommand interface {
	GroupCommand
}

// JoinCommand defines the interface for group join request commands
//
// Extends GroupCommand for the "join" group operation (user-initiated join request)
//
//	Data structure: {
//	    "type"    : i2s(0x89),
//	    "sn"      : 123,
//
//	    "command" : "join",        // Fixed command name: "join"
//	    "time"    : 123.456,       // Timestamp of the join request
//
//	    "group"   : "{GROUP_ID}",  // Target group ID
//	    "text"    : "May I?",      // Optional join request message/comment
//	}
type JoinCommand interface {
	GroupCommand

	// Ask returns the join request message/comment (from "text" field)
	//
	// Typical value: "May I?"
	Ask() string
}

// QuitCommand defines the interface for group leave commands
//
// Extends GroupCommand for the "quit" group operation (member-initiated leave)
//
//	Data structure: {
//	    "type"    : i2s(0x89),
//	    "sn"      : 123,
//
//	    "command" : "quit",        // Fixed command name: "quit"
//	    "time"    : 123.456,       // Timestamp of the leave action
//
//	    "group"   : "{GROUP_ID}",  // Target group ID
//	    "text"    : "Goodbye!",    // Optional leave message/comment
//	}
type QuitCommand interface {
	GroupCommand

	// Bye returns the leave message/comment (from "text" field)
	//
	// Typical value: "Goodbye!"
	Bye() string
}

// ResetCommand defines the interface for group membership reset commands
//
// Extends GroupCommand for the "reset" group operation (replaces deprecated EXPEL)
//
// Used to update/reset group membership (remove members, update roles, etc.)
type ResetCommand interface {
	GroupCommand
}
