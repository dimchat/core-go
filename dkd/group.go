/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
package dkd

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
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
type BaseGroupCommand struct {
	BaseHistoryCommand
	IGroupCommand
}

func (cmd *BaseGroupCommand) Init(dict map[string]interface{}) *BaseGroupCommand {
	if cmd.BaseHistoryCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithCommand(command string, group ID) *BaseGroupCommand {
	if cmd.BaseHistoryCommand.InitWithCommand(command) != nil {
		// group ID
		cmd.SetGroup(group)
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMember(command string, group ID, member ID) *BaseGroupCommand {
	if cmd.InitWithCommand(command, group) != nil {
		// member ID
		cmd.SetMember(member)
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMembers(command string, group ID, members []ID) *BaseGroupCommand {
	if cmd.InitWithCommand(command, group) != nil {
		// member ID list
		cmd.SetMembers(members)
	}
	return cmd
}

//-------- IGroupCommand

/*
 *  Member ID
 */
func (cmd *BaseGroupCommand) Member() ID {
	member := cmd.Get("member")
	if member == nil {
		return nil
	}
	return IDParse(member)
}

func (cmd *BaseGroupCommand) SetMember(member ID) {
	cmd.Set("member", member.String())
}

/*
 *  Member ID list
 */
func (cmd *BaseGroupCommand) Members() []ID {
	members := cmd.Get("members")
	if members == nil {
		return nil
	} else {
		return IDConvert(members)
	}
}

func (cmd *BaseGroupCommand) SetMembers(members []ID) {
	if members == nil {
		cmd.Set("members", nil)
	} else {
		cmd.Set("members", IDRevert(members))
	}
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
type InviteGroupCommand struct {
	BaseGroupCommand
	IInviteCommand
}

func NewInviteCommand(group ID, members []ID) InviteCommand {
	return new(InviteGroupCommand).InitWithMembers(group, members)
}

func (cmd *InviteGroupCommand) Init(dict map[string]interface{}) *InviteGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *InviteGroupCommand) InitWithMember(group ID, member ID) *InviteGroupCommand {
	if cmd.BaseGroupCommand.InitWithMember(INVITE, group, member) != nil {
	}
	return cmd
}

func (cmd *InviteGroupCommand) InitWithMembers(group ID, members []ID) *InviteGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(INVITE, group, members) != nil {
	}
	return cmd
}

//-------- IInviteCommand

func (cmd *InviteGroupCommand) InviteMembers() []ID {
	member := cmd.Member()
	if member != nil {
		return []ID{member}
	} else {
		return cmd.Members()
	}
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
type ExpelGroupCommand struct {
	BaseGroupCommand
	IExpelCommand
}

func NewExpelCommand(group ID, members []ID) ExpelCommand {
	return new(ExpelGroupCommand).InitWithMembers(group, members)
}

func (cmd *ExpelGroupCommand) Init(dict map[string]interface{}) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *ExpelGroupCommand) InitWithMember(group ID, member ID) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.InitWithMember(EXPEL, group, member) != nil {
	}
	return cmd
}

func (cmd *ExpelGroupCommand) InitWithMembers(group ID, members []ID) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(EXPEL, group, members) != nil {
	}
	return cmd
}

//-------- IExpelCommand

func (cmd *ExpelGroupCommand) ExpelMembers() []ID {
	member := cmd.Member()
	if member != nil {
		return []ID{member}
	} else {
		return cmd.Members()
	}
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
type JoinGroupCommand struct {
	BaseGroupCommand
	IJoinCommand
}

func NewJoinCommand(group ID) JoinCommand {
	return new(JoinGroupCommand).InitWithGroup(group)
}

func (cmd *JoinGroupCommand) Init(dict map[string]interface{}) *JoinGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *JoinGroupCommand) InitWithGroup(group ID) *JoinGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(JOIN, group) != nil {
	}
	return cmd
}

//-------- IJoinCommand

func (cmd *JoinGroupCommand) Ask() string {
	text := cmd.Get("text")
	if text == nil {
		return ""
	} else {
		return text.(string)
	}
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
type QuitGroupCommand struct {
	BaseGroupCommand
	IQuitCommand
}

func NewQuitCommand(group ID) QuitCommand {
	return new(QuitGroupCommand).InitWithGroup(group)
}

func (cmd *QuitGroupCommand) Init(dict map[string]interface{}) *QuitGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *QuitGroupCommand) InitWithGroup(group ID) *QuitGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(QUIT, group) != nil {
	}
	return cmd
}

//-------- IQuitCommand

func (cmd *QuitGroupCommand) Bye() string {
	text := cmd.Get("text")
	if text == nil {
		return ""
	} else {
		return text.(string)
	}
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
type ResetGroupCommand struct {
	BaseGroupCommand
	IResetCommand
}

func NewResetCommand(group ID, members []ID) ResetCommand {
	return new(ResetGroupCommand).InitWithMembers(group, members)
}

func (cmd *ResetGroupCommand) Init(dict map[string]interface{}) *ResetGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *ResetGroupCommand) InitWithMembers(group ID, members []ID) *ResetGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(RESET, group, members) != nil {
	}
	return cmd
}

//-------- IResetCommand

func (cmd *ResetGroupCommand) AllMembers() []ID {
	return cmd.Members()
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
type QueryGroupCommand struct {
	BaseGroupCommand
	IQueryCommand
}

func NewQueryCommand(group ID) QueryCommand {
	return new(QueryGroupCommand).InitWithGroup(group)
}

func (cmd *QueryGroupCommand) Init(dict map[string]interface{}) *QueryGroupCommand {
	if cmd.BaseGroupCommand.Init(dict) != nil {
	}
	return cmd
}

func (cmd *QueryGroupCommand) InitWithGroup(group ID) *QueryGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(QUERY, group) != nil {
	}
	return cmd
}

//-------- IQueryCommand

func (cmd *QueryGroupCommand) Text() string {
	text := cmd.Get("text")
	if text == nil {
		return ""
	} else {
		return text.(string)
	}
}
