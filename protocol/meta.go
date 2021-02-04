/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
package protocol

import (
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Meta command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "meta", // command name
 *      ID      : "{ID}", // contact's ID
 *      meta    : {...}   // when meta is empty, means query meta for ID
 *  }
 */
type MetaCommand struct {
	BaseCommand

	_identifier ID
	_meta map[string]interface{}
}

func (cmd *MetaCommand) Init(dict map[string]interface{}) *MetaCommand {
	if cmd.BaseCommand.Init(dict) != nil {
		// lazy load
		cmd._identifier = nil
		cmd._meta = nil
	}
	return cmd
}

func (cmd *MetaCommand) InitWithCommand(command string, id ID, meta Meta) *MetaCommand {
	if cmd.BaseCommand.InitWithCommand(command) != nil {
		// ID
		cmd._identifier = id
		cmd.Set("ID", id.String())
		// meta
		if meta == nil {
			cmd._meta = nil
		} else {
			cmd._meta = meta.GetMap(false)
		}
		cmd.Set("meta", cmd._meta)
	}
	return cmd
}

/**
 *  Response Meta
 *
 * @param identifier - entity ID
 * @param meta - entity Meta
 */
func (cmd *MetaCommand) InitWithMeta(id ID, meta Meta) *MetaCommand {
	return cmd.InitWithCommand(META, id, meta)
}

/**
 *  Query Meta
 *
 * @param identifier - entity ID
 */
func (cmd *MetaCommand) InitWithID(id ID) *MetaCommand {
	return cmd.InitWithCommand(META, id, nil)
}

//-------- setter/getter --------

/*
 *  Entity Meta
 */
func (cmd *MetaCommand) GetMeta() map[string]interface{} {
	if cmd._meta == nil {
		meta := cmd.Get("meta")
		if meta == nil {
			return nil
		}
		cmd._meta = meta.(map[string]interface{})
	}
	return cmd._meta
}

//-------- factories

func MetaCommandQuery(id ID) *MetaCommand {
	return new(MetaCommand).InitWithID(id)
}

func MetaCommandRespond(id ID, meta Meta) *MetaCommand {
	return new(MetaCommand).InitWithMeta(id, meta)
}
