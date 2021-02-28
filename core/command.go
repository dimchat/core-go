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
 *  Command Factory
 */
type GeneralCommandFactory struct {
	ContentFactory
	CommandFactory

	_parser CommandParser
}

func NewGeneralCommandFactory(parser CommandParser) *GeneralCommandFactory {
	return new(GeneralCommandFactory).Init(parser)
}

func (factory *GeneralCommandFactory) Init(parser CommandParser) *GeneralCommandFactory {
	factory._parser = parser
	return factory
}

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

func (factory *GeneralCommandFactory) ParseCommand(cmd map[string]interface{}) Command {
	return factory._parser(cmd)
}

/**
 *  History Command Factory
 */
type HistoryCommandFactory struct {
	GeneralCommandFactory
}

func NewHistoryCommandFactory(parser CommandParser) *HistoryCommandFactory {
	return new(HistoryCommandFactory).Init(parser)
}

func (factory *HistoryCommandFactory) Init(parser CommandParser) *HistoryCommandFactory {
	if factory.GeneralCommandFactory.Init(parser) != nil {
	}
	return factory
}

/**
 *  Group Command Factory
 */
type GroupCommandFactory struct {
	GeneralCommandFactory
}

func NewGroupCommandFactory(parser CommandParser) *GroupCommandFactory {
	return new(GroupCommandFactory).Init(parser)
}

func (factory *GroupCommandFactory) Init(parser CommandParser) *GroupCommandFactory {
	if factory.GeneralCommandFactory.Init(parser) != nil {
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
	CommandRegister(META, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseMetaCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	// Document Command
	CommandRegister(DOCUMENT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseDocumentCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))

	// Group Commands
	CommandRegister("group", NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(INVITE, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(InviteGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(EXPEL, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(ExpelGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(JOIN, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(JoinGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(QUIT, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(QuitGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(QUERY, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(QueryGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
	CommandRegister(RESET, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(ResetGroupCommand)
		cmd.Init(cmd, dict).AutoRelease()
		return cmd
	}))
}

func init() {
	BuildCommandFactories()
}
