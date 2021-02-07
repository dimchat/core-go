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

	_doc Document
}

func (cmd *DocumentCommand) Init(dict map[string]interface{}) *DocumentCommand {
	if cmd.MetaCommand.Init(dict) != nil {
		// lazy load
		cmd._doc = nil
	}
	return cmd
}

func (cmd *DocumentCommand) InitWithMeta(id ID, meta Meta, doc Document) *DocumentCommand {
	if cmd.MetaCommand.InitWithCommand(DOCUMENT, id, meta) != nil {
		// document
		cmd._doc = doc
		if doc != nil {
			cmd.Set("document", doc.GetMap(false))
		}
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

/**
 *  Query Entity Document for updating with current signature
 *
 * @param identifier - entity ID
 * @param signature - document signature
 */
func (cmd *DocumentCommand) InitWithSignature(id ID, signature string) *DocumentCommand {
	if cmd.InitWithID(id) != nil {
		if signature != "" {
			cmd.Set("signature", signature)
		}
	}
	return cmd
}

//-------- setter/getter --------

/*
 *  Entity Document
 */
func (cmd *DocumentCommand) Document() Document {
	if cmd._doc == nil {
		document := cmd.Get("document")
		if document == nil {
			// compatible with v1.0
			document = cmd.Get("profile")
		}
		cmd._doc = DocumentParse(document)
	}
	return cmd._doc
}

func (cmd *DocumentCommand) Signature() string {
	signature := cmd.Get("signature")
	if signature == nil {
		return ""
	}
	return signature.(string)
}

//-------- factories

func DocumentCommandQuery(id ID, signature string) *DocumentCommand {
	return new(DocumentCommand).InitWithSignature(id, signature)
}

func DocumentCommandRespond(id ID, meta Meta, doc Document) *DocumentCommand {
	return new(DocumentCommand).InitWithMeta(id, meta, doc)
}
