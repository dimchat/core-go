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

type BaseMetaCommand struct {
	BaseCommand
	IMetaCommand

	_identifier ID
	_meta Meta
}

func (cmd *BaseMetaCommand) Init(dict map[string]interface{}) *BaseMetaCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd.setID(nil)
		cmd.setMeta(nil)
	}
	return cmd
}

func (cmd *BaseMetaCommand) InitWithCommand(command string, id ID, meta Meta) *BaseMetaCommand {
	if cmd.BaseCommand.InitWithCommand(command) != nil {
		// ID
		cmd.Set("ID", id.String())
		cmd.setID(id)
		// meta
		if meta != nil {
			cmd.Set("meta", meta.GetMap(false))
		}
		cmd.setMeta(meta)
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

func (cmd *BaseMetaCommand) Release() int {
	cnt := cmd.BaseCommand.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		cmd.setID(nil)
		cmd.setMeta(nil)
	}
	return cnt
}

func (cmd *BaseMetaCommand) setID(identifier ID) {
	if identifier != cmd._identifier {
		ObjectRetain(identifier)
		ObjectRelease(cmd._identifier)
		cmd._identifier = identifier
	}
}

func (cmd *BaseMetaCommand) setMeta(meta Meta) {
	if meta != cmd._meta {
		ObjectRetain(meta)
		ObjectRelease(cmd._meta)
		cmd._meta = meta
	}
}

//-------- IMetaCommand

func (cmd *BaseMetaCommand) ID() ID {
	if cmd._identifier == nil {
		cmd.setID(IDParse(cmd.Get("ID")))
	}
	return cmd._identifier
}

func (cmd *BaseMetaCommand) Meta() Meta {
	if cmd._meta == nil {
		meta := cmd.Get("meta")
		cmd.setMeta(MetaParse(meta))
	}
	return cmd._meta
}

//
//  Factories
//

func MetaCommandQuery(id ID) MetaCommand {
	cmd := new(BaseMetaCommand).InitWithID(id)
	ObjectRetain(cmd)
	ObjectAutorelease(cmd)
	return cmd
}

func MetaCommandRespond(id ID, meta Meta) MetaCommand {
	cmd := new(BaseMetaCommand).InitWithMeta(id, meta)
	ObjectRetain(cmd)
	ObjectAutorelease(cmd)
	return cmd
}
