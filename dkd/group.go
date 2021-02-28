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
	. "github.com/dimchat/mkm-go/types"
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

	_member ID
	_members []ID
}

func (cmd *BaseGroupCommand) Init(this GroupCommand, dict map[string]interface{}) *BaseGroupCommand {
	if cmd.BaseHistoryCommand.Init(this, dict) != nil {
		cmd.setMember(nil)
		cmd.setMembers(nil)
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithGroupCommand(this GroupCommand, command string, group ID, member ID, members []ID) *BaseGroupCommand {
	if cmd.InitWithCommand(this, command, group) != nil {
		cmd.SetGroup(group)
		cmd.SetMember(member)
		cmd.setMembers(members)
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithCommand(this GroupCommand, command string, group ID) *BaseGroupCommand {
	if cmd.InitWithGroupCommand(this, command, group, nil, nil) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMember(this GroupCommand, command string, group ID, member ID) *BaseGroupCommand {
	if cmd.InitWithGroupCommand(this, command, group, member, nil) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMembers(this GroupCommand, command string, group ID, members []ID) *BaseGroupCommand {
	if cmd.InitWithGroupCommand(this, command, group, nil, members) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) Release() int {
	cnt := cmd.BaseHistoryCommand.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		cmd.setMember(nil)
		cmd.setMembers(nil)
	}
	return cnt
}

func (cmd *BaseGroupCommand) setMember(member ID) {
	if member != cmd._member {
		ObjectRetain(member)
		ObjectRelease(cmd._member)
		cmd._member = member
	}
}

func (cmd *BaseGroupCommand) setMembers(members []ID) {
	if members != nil {
		for _, item := range members {
			item.Retain()
		}
	}
	if cmd._members != nil {
		for _, item := range cmd._members {
			item.Release()
		}
	}
	cmd._members = members
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
	cmd := new(InviteGroupCommand)
	return cmd.InitWithMembers(cmd, group, members)
}

func (cmd *InviteGroupCommand) Init(this InviteCommand, dict map[string]interface{}) *InviteGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *InviteGroupCommand) InitWithMember(this InviteCommand, group ID, member ID) *InviteGroupCommand {
	if cmd.BaseGroupCommand.InitWithMember(this, INVITE, group, member) != nil {
	}
	return cmd
}

func (cmd *InviteGroupCommand) InitWithMembers(this InviteCommand, group ID, members []ID) *InviteGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(this, INVITE, group, members) != nil {
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
	cmd := new(ExpelGroupCommand)
	return cmd.InitWithMembers(cmd, group, members)
}

func (cmd *ExpelGroupCommand) Init(this ExpelCommand, dict map[string]interface{}) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *ExpelGroupCommand) InitWithMember(this ExpelCommand, group ID, member ID) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.InitWithMember(this, EXPEL, group, member) != nil {
	}
	return cmd
}

func (cmd *ExpelGroupCommand) InitWithMembers(this ExpelCommand, group ID, members []ID) *ExpelGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(this, EXPEL, group, members) != nil {
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
	cmd := new(JoinGroupCommand)
	return cmd.InitWithGroup(cmd, group)
}

func (cmd *JoinGroupCommand) Init(this JoinCommand, dict map[string]interface{}) *JoinGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *JoinGroupCommand) InitWithGroup(this JoinCommand, group ID) *JoinGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(this, JOIN, group) != nil {
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
	cmd := new(QuitGroupCommand)
	return cmd.InitWithGroup(cmd, group)
}

func (cmd *QuitGroupCommand) Init(this QuitCommand, dict map[string]interface{}) *QuitGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *QuitGroupCommand) InitWithGroup(this QuitCommand, group ID) *QuitGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(this, QUIT, group) != nil {
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
	cmd := new(ResetGroupCommand)
	return cmd.InitWithMembers(cmd, group, members)
}

func (cmd *ResetGroupCommand) Init(this ResetCommand, dict map[string]interface{}) *ResetGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *ResetGroupCommand) InitWithMembers(this ResetCommand, group ID, members []ID) *ResetGroupCommand {
	if cmd.BaseGroupCommand.InitWithMembers(this, RESET, group, members) != nil {
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
	cmd := new(QueryGroupCommand)
	return cmd.InitWithGroup(cmd, group)
}

func (cmd *QueryGroupCommand) Init(this QueryCommand, dict map[string]interface{}) *QueryGroupCommand {
	if cmd.BaseGroupCommand.Init(this, dict) != nil {
	}
	return cmd
}

func (cmd *QueryGroupCommand) InitWithGroup(this QueryCommand, group ID) *QueryGroupCommand {
	if cmd.BaseGroupCommand.InitWithCommand(this, QUERY, group) != nil {
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
