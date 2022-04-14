/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2022 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	. "github.com/dimchat/dkd-go/protocol"
)

type ContentParser func(map[string]interface{})Content
type CommandParser func(map[string]interface{})Command

/**
 *  General Content Factory
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type GeneralContentFactory struct {

	_parse ContentParser
}

func NewGeneralContentFactory(fn ContentParser) ContentFactory {
	factory := new(GeneralContentFactory)
	factory.Init(fn)
	return factory
}

func (factory *GeneralContentFactory) Init(fn ContentParser) ContentFactory {
	factory._parse = fn
	return factory
}

//-------- IContentFactory

func (factory *GeneralContentFactory) ParseContent(content map[string]interface{}) Content {
	return factory._parse(content)
}

/**
 *  General Command Factory
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type GeneralCommandFactory struct {

	_parse CommandParser
}

func NewGeneralCommandFactory(fn CommandParser) CommandFactory {
	factory := new(GeneralCommandFactory)
	factory.Init(fn)
	return factory
}

func (factory *GeneralCommandFactory) Init(fn CommandParser) CommandFactory {
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
 *  ~~~~~~~~~~~~~~~~~~~~~~~
 */
type HistoryCommandFactory struct {
	GeneralCommandFactory
}

func NewHistoryCommandFactory(fn CommandParser) CommandFactory {
	factory := new(HistoryCommandFactory)
	factory.Init(fn)
	return factory
}

//func (factory *HistoryCommandFactory) Init(fn CommandParser) CommandFactory {
//	if factory.GeneralCommandFactory.Init(fn) != nil {
//	}
//	return factory
//}

/**
 *  Group Command Factory
 *  ~~~~~~~~~~~~~~~~~~~~~
 */
type GroupCommandFactory struct {
	GeneralCommandFactory
}

func NewGroupCommandFactory(fn CommandParser) CommandFactory {
	factory := new(GroupCommandFactory)
	factory.Init(fn)
	return factory
}

//func (factory *GroupCommandFactory) Init(fn CommandParser) CommandFactory {
//	if factory.GeneralCommandFactory.Init(fn) != nil {
//	}
//	return factory
//}

//-------- IContentFactory

func (factory *GroupCommandFactory) ParseContent(content map[string]interface{}) Content {
	command := CommandGetName(content)
	// get factory by command name
	cmdFactory := CommandGetFactory(command)
	if cmdFactory == nil {
		cmdFactory = factory
	}
	return cmdFactory.ParseCommand(content)
}
