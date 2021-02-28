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

type BaseDocumentCommand struct {
	BaseMetaCommand
	IDocumentCommand

	_doc Document
}

func (cmd *BaseDocumentCommand) Init(this DocumentCommand, dict map[string]interface{}) *BaseDocumentCommand {
	if cmd.BaseMetaCommand.Init(this, dict) != nil {
		// lazy load
		cmd.setDocument(nil)
	}
	return cmd
}

func (cmd *BaseDocumentCommand) InitWithMeta(this DocumentCommand, id ID, meta Meta, doc Document) *BaseDocumentCommand {
	if cmd.BaseMetaCommand.InitWithCommand(this, DOCUMENT, id, meta) != nil {
		// document
		if doc != nil {
			cmd.Set("document", doc.GetMap(false))
		}
		cmd.setDocument(doc)
	}
	return cmd
}

func (cmd *BaseDocumentCommand) InitWithDocument(this DocumentCommand, id ID, doc Document) *BaseDocumentCommand {
	return cmd.InitWithMeta(this, id, nil, doc)
}

/**
 *  Query Meta
 *
 * @param identifier - entity ID
 */
func (cmd *BaseDocumentCommand) InitWithID(this DocumentCommand, id ID) *BaseDocumentCommand {
	return cmd.InitWithMeta(this, id, nil, nil)
}

/**
 *  Query Entity Document for updating with current signature
 *
 * @param identifier - entity ID
 * @param signature - document signature
 */
func (cmd *BaseDocumentCommand) InitWithSignature(this DocumentCommand, id ID, signature string) *BaseDocumentCommand {
	if cmd.InitWithID(this, id) != nil {
		if signature != "" {
			cmd.Set("signature", signature)
		}
	}
	return cmd
}

func (cmd *BaseDocumentCommand) Release() int {
	cnt := cmd.BaseMetaCommand.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		cmd.setDocument(nil)
	}
	return cnt
}

func (cmd *BaseDocumentCommand) setDocument(doc Document) {
	if doc != cmd._doc {
		ObjectRetain(doc)
		ObjectRelease(cmd._doc)
		cmd._doc = doc
	}
}

//-------- IDocumentCommand

/*
 *  Entity Document
 */
func (cmd *BaseDocumentCommand) Document() Document {
	if cmd._doc == nil {
		document := cmd.Get("document")
		if document == nil {
			// compatible with v1.0
			document = cmd.Get("profile")
		}
		cmd.setDocument(DocumentParse(document))
	}
	return cmd._doc
}

func (cmd *BaseDocumentCommand) Signature() string {
	signature := cmd.Get("signature")
	if signature == nil {
		return ""
	}
	return signature.(string)
}

//
//  Factories
//

func DocumentCommandQuery(id ID, signature string) DocumentCommand {
	cmd := new(BaseDocumentCommand)
	cmd.InitWithSignature(cmd, id, signature).AutoRelease()
	return cmd
}

func DocumentCommandRespond(id ID, meta Meta, doc Document) DocumentCommand {
	cmd := new(BaseDocumentCommand)
	cmd.InitWithMeta(cmd, id, meta, doc).AutoRelease()
	return cmd
}
