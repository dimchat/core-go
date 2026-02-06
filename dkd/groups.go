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
 *  History Command Content
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
type BaseHistoryCommand struct {
	//HistoryCommand
	BaseCommand
}

func (content *BaseHistoryCommand) Init(cmd string) HistoryCommand {
	if content.BaseCommand.InitWithType(ContentType.HISTORY, cmd) != nil {
	}
	return content
}

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
type BaseGroupCommand struct {
	//GroupCommand
	BaseHistoryCommand
}

func (content *BaseGroupCommand) Init(cmd string, group ID, members []ID) GroupCommand {
	if content.BaseHistoryCommand.Init(cmd) != nil {
		content.SetGroup(group)
		if members != nil {
			content.SetMembers(members)
		}
	}
	return content
}

// Override
func (content *BaseGroupCommand) Members() []ID {
	members := content.Get("members")
	if members != nil {
		return IDConvert(members)
	}
	// get from 'member'
	single := ParseID(content.Get("member"))
	if single != nil {
		return []ID{
			single,
		}
	}
	// failed to get group members
	return nil
}

func (content *BaseGroupCommand) SetMembers(members []ID) {
	if members == nil {
		content.Remove("members")
	} else {
		content.Set("members", IDRevert(members))
	}
	content.Remove("member")
}

type InviteGroupCommand struct {
	//InviteCommand
	BaseGroupCommand
}

func (content *InviteGroupCommand) Init(group ID, members []ID) InviteCommand {
	if content.BaseGroupCommand.Init(INVITE, group, members) != nil {
	}
	return content
}

/**
 *  Deprecated, use 'reset' instead
 */
type ExpelGroupCommand struct {
	//ExpelCommand
	BaseGroupCommand
}

func (content *ExpelGroupCommand) Init(group ID, members []ID) ExpelCommand {
	if content.BaseGroupCommand.Init(EXPEL, group, members) != nil {
	}
	return content
}

type JoinGroupCommand struct {
	//JoinCommand
	BaseGroupCommand
}

func (content *JoinGroupCommand) Init(group ID) JoinCommand {
	if content.BaseGroupCommand.Init(JOIN, group, nil) != nil {
	}
	return content
}

// Override
func (content *JoinGroupCommand) Ask() string {
	return content.GetString("text", "")
}

type QuitGroupCommand struct {
	//QuitCommand
	BaseGroupCommand
}

func (content *QuitGroupCommand) Init(group ID) QuitCommand {
	if content.BaseGroupCommand.Init(QUIT, group, nil) != nil {
	}
	return content
}

// Override
func (content *QuitGroupCommand) Bye() string {
	return content.GetString("text", "")
}

type ResetGroupCommand struct {
	//ResetCommand
	BaseGroupCommand
}

func (content *ResetGroupCommand) Init(group ID, members []ID) ResetCommand {
	if content.BaseGroupCommand.Init(RESET, group, members) != nil {
	}
	return content
}
