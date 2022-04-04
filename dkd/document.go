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

	_doc Document
}

func (cmd *BaseDocumentCommand) Init(dict map[string]interface{}) *BaseDocumentCommand {
	if cmd.BaseMetaCommand.Init(dict) != nil {
		// lazy load
		cmd._doc = nil
	}
	return cmd
}

func (cmd *BaseDocumentCommand) InitWithMeta(id ID, meta Meta, doc Document) *BaseDocumentCommand {
	if cmd.BaseMetaCommand.InitWithCommand(DOCUMENT, id, meta) != nil {
		// document
		if !ValueIsNil(doc) {
			cmd.Set("document", doc.GetMap(false))
		}
		cmd._doc = doc
	}
	return cmd
}

func (cmd *BaseDocumentCommand) InitWithDocument(id ID, doc Document) *BaseDocumentCommand {
	if cmd.InitWithMeta(id, nil, doc) != nil {
	}
	return cmd
}

/**
 *  Query Meta
 *
 * @param identifier - entity ID
 */
func (cmd *BaseDocumentCommand) InitWithID(id ID) *BaseDocumentCommand {
	if cmd.InitWithMeta(id, nil, nil) != nil {
	}
	return cmd
}

/**
 *  Query Entity Document for updating with current signature
 *
 * @param identifier - entity ID
 * @param signature - document signature
 */
func (cmd *BaseDocumentCommand) InitWithSignature(id ID, signature string) *BaseDocumentCommand {
	if cmd.InitWithID(id) != nil {
		if signature != "" {
			cmd.Set("signature", signature)
		}
	}
	return cmd
}

//-------- IDocumentCommand

/*
 *  Entity Document
 */
func (cmd *BaseDocumentCommand) Document() Document {
	if cmd._doc == nil {
		cmd._doc = DocumentParse(cmd.Get("document"))
	}
	return cmd._doc
}

func (cmd *BaseDocumentCommand) Signature() string {
	b64 := cmd.Get("signature")
	if b64 == nil {
		return ""
	}
	return b64.(string)
}

//
//  Factories
//

func DocumentCommandQuery(id ID, signature string) DocumentCommand {
	cmd := new(BaseDocumentCommand)
	cmd.InitWithSignature(id, signature)
	return cmd
}

func DocumentCommandRespond(id ID, meta Meta, doc Document) DocumentCommand {
	cmd := new(BaseDocumentCommand)
	cmd.InitWithMeta(id, meta, doc)
	return cmd
}
