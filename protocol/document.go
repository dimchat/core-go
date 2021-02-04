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
 *  Document Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command   : "profile", // command name
 *      ID        : "{ID}",    // entity ID
 *      meta      : {...},     // only for handshaking with new friend
 *      profile   : {...},     // when profile is empty, means query for ID
 *      signature : "..."      // old profile's signature for querying
 *  }
 */
type DocumentCommand struct {
	MetaCommand

	_profile map[string]interface{}
}

func (cmd *DocumentCommand) Init(dict map[string]interface{}) *DocumentCommand {
	if cmd.MetaCommand.Init(dict) != nil {
		// lazy load
		cmd._profile = nil
	}
	return cmd
}
func (cmd *DocumentCommand) InitWithMeta(id ID, meta Meta, doc Document) *DocumentCommand {
	if cmd.MetaCommand.InitWithCommand(PROFILE, id, meta) != nil {
		// document
		if doc == nil {
			cmd._profile = nil
		} else {
			cmd._profile = doc.GetMap(false)
		}
		cmd.Set("profile", cmd._profile)
	}
	return cmd
}

func (cmd *DocumentCommand) InitWithDocument(id ID, doc Document) *DocumentCommand {
	return cmd.InitWithMeta(id, nil, doc)
}

/**
 *  Query Meta
 *
 * @param identifier - entity ID
 */
func (cmd *DocumentCommand) InitWithID(id ID) *DocumentCommand {
	return cmd.InitWithMeta(id, nil, nil)
}

//-------- setter/getter --------

/*
 *  Entity Document
 */
func (cmd *DocumentCommand) GetDocument() map[string]interface{} {
	doc := cmd.Get("document")
	if doc == nil {
		// compatible with v1.0
		doc = cmd.Get("profile")
		if doc == nil {
			return nil
		}
	}
	return doc.(map[string]interface{})
}

//-------- factories

func DocumentCommandQuery(id ID) *DocumentCommand {
	return new(DocumentCommand).InitWithID(id)
}

func DocumentCommandRespond(id ID, meta Meta, doc Document) *DocumentCommand {
	return new(DocumentCommand).InitWithMeta(id, meta, doc)
}
