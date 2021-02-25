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

type BaseMetaCommand struct {
	BaseCommand
	IMetaCommand

	_identifier ID
	_meta Meta
}

func (cmd *BaseMetaCommand) Init(dict map[string]interface{}) *BaseMetaCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._identifier = nil
		cmd._meta = nil
	}
	return cmd
}

func (cmd *BaseMetaCommand) InitWithCommand(command string, id ID, meta Meta) *BaseMetaCommand {
	if cmd.BaseCommand.InitWithCommand(command) != nil {
		// ID
		cmd._identifier = id
		cmd.Set("ID", id.String())
		// meta
		cmd._meta = meta
		if meta != nil {
			cmd.Set("meta", meta.GetMap(false))
		}
	}
	return cmd
}

/**
 *  Response Meta
 *
 * @param identifier - entity ID
 * @param meta - entity Meta
 */
func (cmd *BaseMetaCommand) InitWithMeta(id ID, meta Meta) *BaseMetaCommand {
	return cmd.InitWithCommand(META, id, meta)
}

/**
 *  Query Meta
 *
 * @param identifier - entity ID
 */
func (cmd *BaseMetaCommand) InitWithID(id ID) *BaseMetaCommand {
	return cmd.InitWithCommand(META, id, nil)
}

//-------- IMetaCommand

func (cmd *BaseMetaCommand) ID() ID {
	return IDParse(cmd.Get("ID"))
}

func (cmd *BaseMetaCommand) Meta() Meta {
	if cmd._meta == nil {
		meta := cmd.Get("meta")
		cmd._meta = MetaParse(meta)
	}
	return cmd._meta
}

//
//  Factories
//

func MetaCommandQuery(id ID) MetaCommand {
	return new(BaseMetaCommand).InitWithID(id)
}

func MetaCommandRespond(id ID, meta Meta) MetaCommand {
	return new(BaseMetaCommand).InitWithMeta(id, meta)
}
