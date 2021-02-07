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
 *  History command: {
 *      type : 0x89,
 *      sn   : 123,
 *
 *      command : "...", // command name
 *      time    : 0,     // command timestamp
 *      extra   : info   // command parameters
 *  }
 */
type GroupCommand struct {
	HistoryCommand
}

func (cmd *GroupCommand) Init(dict map[string]interface{}) *GroupCommand {
	if cmd.HistoryCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *GroupCommand) InitWithCommand(command string, group ID) *GroupCommand {
	if cmd.HistoryCommand.InitWithCommand(command) != nil {
		// group ID
		cmd.SetGroup(group)
	}
	return cmd
}

func (cmd *GroupCommand) InitWithMember(command string, group ID, member ID) *GroupCommand {
	if cmd.InitWithCommand(command, group) != nil {
		// member ID
		cmd.SetMember(member)
	}
	return cmd
}

func (cmd *GroupCommand) InitWithMembers(command string, group ID, members []ID) *GroupCommand {
	if cmd.InitWithCommand(command, group) != nil {
		// member ID list
		cmd.SetMembers(members)
	}
	return cmd
}

//-------- setter/getter --------

/*
 *  Member ID
 */
func (cmd *GroupCommand) Member() ID {
	member := cmd.Get("member")
	if member == nil {
		return nil
	}
	return IDParse(member)
}

func (cmd *GroupCommand) SetMember(member ID)  {
	cmd.Set("member", member.String())
}

/*
 *  Member ID list
 */
func (cmd *GroupCommand) Members() []ID {
	members := cmd.Get("members")
	if members == nil {
		return nil
	}
	switch members.(type) {
	case []interface{}:
		return IDConvert(members.([]interface{}))
	default:
		panic(members)
		return nil
	}
}

func (cmd *GroupCommand) SetMembers(members []ID)  {
	if members == nil {
		cmd.Set("members", nil)
	} else {
		cmd.Set("members", IDRevert(members))
	}
}

//-------- Group Commands

// "invite"
type InviteCommand struct {
	GroupCommand
}

func (cmd *InviteCommand) Init(dict map[string]interface{}) *InviteCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *InviteCommand) InitWithMember(group ID, member ID) *InviteCommand {
	if cmd.GroupCommand.InitWithMember(INVITE, group, member) != nil {
	}
	return cmd
}

func (cmd *InviteCommand) InitWithMembers(group ID, members []ID) *InviteCommand {
	if cmd.GroupCommand.InitWithMembers(INVITE, group, members) != nil {
	}
	return cmd
}

// "expel"
type ExpelCommand struct {
	GroupCommand
}

func (cmd *ExpelCommand) Init(dict map[string]interface{}) *ExpelCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *ExpelCommand) InitWithMember(group ID, member ID) *ExpelCommand {
	if cmd.GroupCommand.InitWithMember(EXPEL, group, member) != nil {
	}
	return cmd
}

func (cmd *ExpelCommand) InitWithMembers(group ID, members []ID) *ExpelCommand {
	if cmd.GroupCommand.InitWithMembers(EXPEL, group, members) != nil {
	}
	return cmd
}

// "join"
type JoinCommand struct {
	GroupCommand
}

func (cmd *JoinCommand) Init(dict map[string]interface{}) *JoinCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *JoinCommand) InitWithGroup(group ID) *JoinCommand {
	if cmd.GroupCommand.InitWithCommand(JOIN, group) != nil {
	}
	return cmd
}

// "quit"
type QuitCommand struct {
	GroupCommand
}

func (cmd *QuitCommand) Init(dict map[string]interface{}) *QuitCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *QuitCommand) InitWithGroup(group ID) *QuitCommand {
	if cmd.GroupCommand.InitWithCommand(QUIT, group) != nil {
	}
	return cmd
}

// "reset"
type ResetCommand struct {
	GroupCommand
}

func (cmd *ResetCommand) Init(dict map[string]interface{}) *ResetCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *ResetCommand) InitWithMembers(group ID, members []ID) *ResetCommand {
	if cmd.GroupCommand.InitWithMembers(RESET, group, members) != nil {
	}
	return cmd
}

// "query"
/**
 *  NOTICE:
 *      This command is just for querying group info,
 *      should not be saved in group history
 */
type QueryCommand struct {
	GroupCommand
}

func (cmd *QueryCommand) Init(dict map[string]interface{}) *QueryCommand {
	if cmd.GroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *QueryCommand) InitWithGroup(group ID) *QueryCommand {
	if cmd.GroupCommand.InitWithCommand(QUERY, group) != nil {
	}
	return cmd
}
