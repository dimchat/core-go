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
 *  Meta Command Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x88),
 *      "sn"   : 123,
 *
 *      "command" : "meta", // command name
 *      "did"     : "{ID}", // contact's ID
 *      "meta"    : {...}   // when meta is null, means query meta for ID
 *  }
 *  </pre></blockquote>
 */
type BaseMetaCommand struct {
	//MetaCommand
	BaseCommand

	_meta Meta
}

func (content *BaseMetaCommand) InitWithMap(dict StringKeyMap) MetaCommand {
	if content.BaseCommand.InitWithMap(dict) != nil {
		// lazy load
		content._meta = nil
	}
	return content
}

func (content *BaseMetaCommand) Init(cmd string, did ID, meta Meta) MetaCommand {
	if content.BaseCommand.Init(cmd) != nil {
		// ID
		content.Set("did", did.String())
		// meta
		if meta != nil {
			content.Set("meta", meta.Map())
		}
		content._meta = meta
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
	meta := content._meta
	if meta == nil {
		info := content.Get("meta")
		meta = ParseMeta(info)
		content._meta = meta
	}
	return meta
}

/**
 *  Document Command Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x88),
 *      "sn"   : 123,
 *
 *      "command"   : "documents", // command name
 *      "did"       : "{ID}",      // entity ID
 *      "meta"      : {...},       // only for handshaking with new friend
 *      "documents" : [...],       // when this is null, means to query
 *      "last_time" : 12345        // old document time for querying
 *  }
 *  </pre></blockquote>
 */
type BaseDocumentCommand struct {
	//DocumentCommand
	BaseMetaCommand

	_docs []Document
}

func (content *BaseDocumentCommand) InitWithMap(dict StringKeyMap) DocumentCommand {
	if content.BaseMetaCommand.InitWithMap(dict) != nil {
		// lazy load
		content._docs = nil
	}
	return content
}

func (content *BaseDocumentCommand) Init(did ID, meta Meta, docs []Document) DocumentCommand {
	if content.BaseMetaCommand.Init(DOCUMENTS, did, meta) != nil {
		// documents
		if docs != nil {
			content.Set("documents", DocumentRevert(docs))
		}
		content._docs = docs
	}
	return content
}

// Override
func (content *BaseDocumentCommand) Documents() []Document {
	docs := content._docs
	if docs == nil {
		array := content.Get("documents")
		if array != nil {
			docs = DocumentConvert(array)
		} else {
			docs = make([]Document, 0)
		}
		content._docs = docs
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
