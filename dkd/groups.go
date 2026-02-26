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
 *  History Command
 */

type BaseHistoryCommand struct {
	//HistoryCommand
	*BaseCommand
}

func NewBaseHistoryCommand(dict StringKeyMap, cmd string) *BaseHistoryCommand {
	return &BaseHistoryCommand{
		BaseCommand: NewBaseCommand(dict, ContentType.HISTORY, cmd),
	}
}

// Override
func (content *BaseHistoryCommand) Event() string {
	event := content.Get("event")
	if event == nil {
		event = content.Get("command")
	}
	return ConvertString(event, "")
}

/**
 *  Group History
 */

type BaseGroupCommand struct {
	//GroupCommand
	*BaseHistoryCommand
}

func NewBaseGroupCommand(dict StringKeyMap, cmd string, group ID, members []ID) *BaseGroupCommand {
	content := &BaseGroupCommand{
		BaseHistoryCommand: NewBaseHistoryCommand(dict, cmd),
	}
	if group != nil {
		content.SetGroup(group)
	}
	if members != nil {
		content.SetMembers(members)
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
	*BaseGroupCommand
}

func NewInviteGroupCommand(dict StringKeyMap, group ID, members []ID) *InviteGroupCommand {
	return &InviteGroupCommand{
		BaseGroupCommand: NewBaseGroupCommand(dict, INVITE, group, members),
	}
}

// Deprecated, use 'reset' instead
type ExpelGroupCommand struct {
	//ExpelCommand
	*BaseGroupCommand
}

func NewExpelGroupCommand(dict StringKeyMap, group ID, members []ID) *ExpelGroupCommand {
	return &ExpelGroupCommand{
		BaseGroupCommand: NewBaseGroupCommand(dict, EXPEL, group, members),
	}
}

type JoinGroupCommand struct {
	//JoinCommand
	*BaseGroupCommand
}

func NewJoinGroupCommand(dict StringKeyMap, group ID) *JoinGroupCommand {
	return &JoinGroupCommand{
		BaseGroupCommand: NewBaseGroupCommand(dict, JOIN, group, nil),
	}
}

// Override
func (content *JoinGroupCommand) Ask() string {
	return content.GetString("text", "")
}

type QuitGroupCommand struct {
	//QuitCommand
	*BaseGroupCommand
}

func NewQuitGroupCommand(dict StringKeyMap, group ID) *QuitGroupCommand {
	return &QuitGroupCommand{
		BaseGroupCommand: NewBaseGroupCommand(dict, QUIT, group, nil),
	}
}

// Override
func (content *QuitGroupCommand) Bye() string {
	return content.GetString("text", "")
}

type ResetGroupCommand struct {
	//ResetCommand
	*BaseGroupCommand
}

func NewResetGroupCommand(dict StringKeyMap, group ID, members []ID) *ResetGroupCommand {
	return &ResetGroupCommand{
		BaseGroupCommand: NewBaseGroupCommand(dict, RESET, group, members),
	}
}
