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

/**
 *  Meta Command
 */

type BaseMetaCommand struct {
	//MetaCommand
	*BaseCommand

	meta Meta
}

func NewBaseMetaCommand(dict StringKeyMap, cmd string, did ID, meta Meta) *BaseMetaCommand {
	if dict != nil {
		// init meta command with map
		return &BaseMetaCommand{
			BaseCommand: NewBaseCommand(dict, "", ""),
			// lazy load
			meta: nil,
		}
	}
	// new meta command
	if cmd == "" {
		cmd = META
	}
	content := &BaseMetaCommand{
		BaseCommand: NewBaseCommand(nil, "", cmd),
		meta:        meta,
	}
	// ID
	content.Set("did", did.String())
	// meta
	if meta != nil {
		content.Set("meta", meta.Map())
	}
	return content
}

// Override
func (content *BaseMetaCommand) ID() ID {
	did := content.Get("did")
	return ParseID(did)
}

// Override
func (content *BaseMetaCommand) Meta() Meta {
	meta := content.meta
	if meta == nil {
		info := content.Get("meta")
		meta = ParseMeta(info)
		content.meta = meta
	}
	return meta
}

/**
 *  Document Command
 */

type BaseDocumentCommand struct {
	//DocumentCommand
	*BaseMetaCommand

	docs []Document
}

func NewBaseDocumentCommand(dict StringKeyMap, did ID, meta Meta, docs []Document, lastTime Time) *BaseDocumentCommand {
	if dict != nil {
		// init document command with map
		return &BaseDocumentCommand{
			BaseMetaCommand: NewBaseMetaCommand(dict, "", nil, nil),
			// lazy load
			docs: nil,
		}
	}
	// new document command
	content := &BaseDocumentCommand{
		BaseMetaCommand: NewBaseMetaCommand(nil, DOCUMENTS, did, meta),
		docs:            docs,
	}
	// documents
	if docs != nil {
		content.Set("documents", DocumentRevert(docs))
	}
	// last document time
	if lastTime != nil {
		content.SetLastTime(lastTime)
	}
	return content
}

// Override
func (content *BaseDocumentCommand) Documents() []Document {
	docs := content.docs
	if docs == nil {
		array := content.Get("documents")
		if array != nil {
			docs = DocumentConvert(array)
		} else {
			//docs = []Document{}
		}
		content.docs = docs
	}
	return docs
}

// Override
func (content *BaseDocumentCommand) LastTime() Time {
	return content.GetTime("last_time", nil)
}

func (content *BaseDocumentCommand) SetLastTime(docTime Time) {
	content.SetTime("last_time", docTime)
}
