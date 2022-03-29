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
package core

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/protocol"
)

type CommandParser func(map[string]interface{})Command

/**
 *  General Command Factory
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type GeneralCommandFactory struct {
	IContentFactory
	ICommandFactory

	_parse CommandParser
}

func NewGeneralCommandFactory(fn CommandParser) *GeneralCommandFactory {
	return new(GeneralCommandFactory).Init(fn)
}

func (factory *GeneralCommandFactory) Init(fn CommandParser) *GeneralCommandFactory {
	factory._parse = fn
	return factory
}

//-------- IContentFactory

func (factory *GeneralCommandFactory) ParseContent(content map[string]interface{}) Content {
	// get factory by command name
	command := CommandGetName(content)
	cmdFactory := CommandGetFactory(command)
	if cmdFactory == nil {
		// check for group command
		if ContentGetGroup(content) != nil {
			cmdFactory = CommandGetFactory("group")
		}
		if cmdFactory == nil {
			cmdFactory = factory
		}
	}
	return cmdFactory.ParseCommand(content)
}

//-------- ICommandFactory

func (factory *GeneralCommandFactory) ParseCommand(cmd map[string]interface{}) Command {
	return factory._parse(cmd)
}

/**
 *  History Command Factory
 */
type HistoryCommandFactory struct {
	GeneralCommandFactory
}

func NewHistoryCommandFactory(fn CommandParser) *HistoryCommandFactory {
	return new(HistoryCommandFactory).Init(fn)
}

func (factory *HistoryCommandFactory) Init(fn CommandParser) *HistoryCommandFactory {
	if factory.GeneralCommandFactory.Init(fn) != nil {
	}
	return factory
}

/**
 *  Group Command Factory
 */
type GroupCommandFactory struct {
	GeneralCommandFactory
}

func NewGroupCommandFactory(fn CommandParser) *GroupCommandFactory {
	return new(GroupCommandFactory).Init(fn)
}

func (factory *GroupCommandFactory) Init(fn CommandParser) *GroupCommandFactory {
	if factory.GeneralCommandFactory.Init(fn) != nil {
	}
	return factory
}

func (factory *GroupCommandFactory) ParseContent(content map[string]interface{}) Content {
	command := CommandGetName(content)
	// get factory by command name
	cmdFactory := CommandGetFactory(command)
	if cmdFactory == nil {
		cmdFactory = factory
	}
	return cmdFactory.ParseCommand(content)
}

/**
 *  Register core command parsers
 */
func BuildCommandFactories() {
	// Meta Command
	CommandSetFactory(META, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseMetaCommand).Init(dict)
	}))
	// Document Command
	factory := NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseDocumentCommand).Init(dict)
	})
	CommandSetFactory(DOCUMENT, factory)
	CommandSetFactory("profile", factory)
	CommandSetFactory("visa", factory)
	CommandSetFactory("bulletin", factory)

	// Group Commands
	CommandSetFactory("group", NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseGroupCommand).Init(dict)
	}))
	CommandSetFactory(INVITE, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(InviteGroupCommand).Init(dict)
	}))
	CommandSetFactory(EXPEL, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(ExpelGroupCommand).Init(dict)
	}))
	CommandSetFactory(JOIN, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(JoinGroupCommand).Init(dict)
	}))
	CommandSetFactory(QUIT, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(QuitGroupCommand).Init(dict)
	}))
	CommandSetFactory(QUERY, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(QueryGroupCommand).Init(dict)
	}))
	CommandSetFactory(RESET, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(ResetGroupCommand).Init(dict)
	}))
}
