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
	. "github.com/dimchat/core-go/protocol/cmd"
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

func NewCommandFactory(parser CommandParser) *GeneralCommandFactory {
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
	if factory._parser == nil {
		return NewCommand(cmd)
	}
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
	factory.GeneralCommandFactory.Init(parser)
	return factory
}

func (factory *HistoryCommandFactory) ParseCommand(cmd map[string]interface{}) Command {
	if factory._parser == nil {
		return NewHistoryCommand(cmd)
	}
	return factory._parser(cmd)
}

/**
 *  Group Command Factory
 */
type GroupCommandFactory struct {
	HistoryCommandFactory
}

func NewGroupCommandFactory(parser CommandParser) *GroupCommandFactory {
	return new(GroupCommandFactory).Init(parser)
}

func (factory *GroupCommandFactory) Init(parser CommandParser) *GroupCommandFactory {
	factory.HistoryCommandFactory.Init(parser)
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

func (factory *GroupCommandFactory) ParseCommand(cmd map[string]interface{}) Command {
	if factory._parser == nil {
		return NewGroupCommand(cmd)
	}
	return factory._parser(cmd)
}
