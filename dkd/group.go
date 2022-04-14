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

	_member ID
	_members []ID
}

/* designated initializer */
func (cmd *BaseGroupCommand) Init(dict map[string]interface{}) GroupCommand {
	if cmd.BaseHistoryCommand.Init(dict) != nil {
		// lazy load
		cmd._member = nil
		cmd._members = nil
	}
	return cmd
}

/* designated initializer */
func (cmd *BaseGroupCommand) InitWithGroupCommand(command string, group ID, member ID, members []ID) GroupCommand {
	if cmd.BaseHistoryCommand.InitWithCommand(command) != nil {
		cmd.SetGroup(group)
		cmd.SetMember(member)
		cmd.SetMembers(members)
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithCommand(command string, group ID) GroupCommand {
	if cmd.InitWithGroupCommand(command, group, nil, nil) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMember(command string, group ID, member ID) GroupCommand {
	if cmd.InitWithGroupCommand(command, group, member, nil) != nil {
	}
	return cmd
}

func (cmd *BaseGroupCommand) InitWithMembers(command string, group ID, members []ID) GroupCommand {
	if cmd.InitWithGroupCommand(command, group, nil, members) != nil {
	}
	return cmd
}

//-------- IGroupCommand

/*
 *  Member ID
 */
func (cmd *BaseGroupCommand) Member() ID {
	if cmd._member == nil {
		member := cmd.Get("member")
		cmd._member = IDParse(member)
	}
	return cmd._member
}

func (cmd *BaseGroupCommand) SetMember(member ID) {
	if ValueIsNil(member) {
		cmd.Remove("member")
	} else {
		cmd.Set("member", member.String())
	}
	cmd._member = member
}

/*
 *  Member ID list
 */
func (cmd *BaseGroupCommand) Members() []ID {
	if cmd._members == nil {
		members := cmd.Get("members")
		if members != nil {
			cmd._members = IDConvert(members)
		}
	}
	return cmd._members
}

func (cmd *BaseGroupCommand) SetMembers(members []ID) {
	if ValueIsNil(members) {
		cmd.Remove("members")
	} else {
		cmd.Set("members", IDRevert(members))
	}
	cmd._members = members
}
