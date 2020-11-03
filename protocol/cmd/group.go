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
	. "github.com/dimchat/mkm-go/mkm"
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

func (cmd *GroupCommand) Init(dictionary map[string]interface{}) *GroupCommand {
	if cmd.HistoryCommand.Init(dictionary) != nil {
		// init
	}
	return cmd
}

func (cmd *GroupCommand) InitWithCommand(command string, group *ID) *GroupCommand {
	if cmd.HistoryCommand.InitWithCommand(command) != nil {
		// group ID
		cmd.SetGroup(group)
	}
	return cmd
}

func (cmd *GroupCommand) InitWithMember(command string, group *ID, member *ID) *GroupCommand {
	if cmd.InitWithCommand(command, group) != nil {
		// member ID
		cmd.SetMember(member)
	}
	return cmd
}

func (cmd *GroupCommand) InitWithMembers(command string, group *ID, members []string) *GroupCommand {
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
func (cmd *GroupCommand) GetMember() *ID {
	member := cmd.Get("member")
	if member == nil {
		return nil
	}
	handler := cmd.GetDelegate()
	return (*handler).GetID(member)
}

func (cmd *GroupCommand) SetMember(member *ID)  {
	cmd.Set("member", member.String.String())
}

/*
 *  Member ID list
 */
func (cmd *GroupCommand) GetMembers() []string {
	members := cmd.Get("members")
	if members == nil {
		return nil
	}
	return members.([]string)
}

func (cmd *GroupCommand) SetMembers(members []string)  {
	cmd.Set("members", members)
}

//-------- Group Commands

// "invite"
type InviteCommand struct {
	GroupCommand
}

func (cmd *InviteCommand) Init(dictionary map[string]interface{}) *InviteCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *InviteCommand) InitWithMember(group *ID, member *ID) *InviteCommand {
	cmd.GroupCommand.InitWithMember(INVITE, group, member)
	return cmd
}

func (cmd *InviteCommand) InitWithMembers(group *ID, members []string) *InviteCommand {
	cmd.GroupCommand.InitWithMembers(INVITE, group, members)
	return cmd
}

// "expel"
type ExpelCommand struct {
	GroupCommand
}

func (cmd *ExpelCommand) Init(dictionary map[string]interface{}) *ExpelCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *ExpelCommand) InitWithMember(group *ID, member *ID) *ExpelCommand {
	cmd.GroupCommand.InitWithMember(EXPEL, group, member)
	return cmd
}

func (cmd *ExpelCommand) InitWithMembers(group *ID, members []string) *ExpelCommand {
	cmd.GroupCommand.InitWithMembers(EXPEL, group, members)
	return cmd
}

// "join"
type JoinCommand struct {
	GroupCommand
}

func (cmd *JoinCommand) Init(dictionary map[string]interface{}) *JoinCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *JoinCommand) InitWithGroup(group *ID) *JoinCommand {
	cmd.GroupCommand.InitWithCommand(JOIN, group)
	return cmd
}

// "quit"
type QuitCommand struct {
	GroupCommand
}

func (cmd *QuitCommand) Init(dictionary map[string]interface{}) *QuitCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *QuitCommand) InitWithGroup(group *ID) *QuitCommand {
	cmd.GroupCommand.InitWithCommand(QUIT, group)
	return cmd
}

// "reset"
type ResetCommand struct {
	GroupCommand
}

func (cmd *ResetCommand) Init(dictionary map[string]interface{}) *ResetCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *ResetCommand) InitWithMembers(group *ID, members []string) *ResetCommand {
	cmd.GroupCommand.InitWithMembers(RESET, group, members)
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

func (cmd *QueryCommand) Init(dictionary map[string]interface{}) *QueryCommand {
	cmd.GroupCommand.Init(dictionary)
	return cmd
}

func (cmd *QueryCommand) InitWithGroup(group *ID) *QueryCommand {
	cmd.GroupCommand.InitWithCommand(QUERY, group)
	return cmd
}
