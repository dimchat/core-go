/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2026 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

func NewCommandWithMap(dict StringKeyMap) Command {
	content := &BaseCommand{}
	content.InitWithMap(dict)
	return content
}

func NewHistoryCommandWithMap(dict StringKeyMap) HistoryCommand {
	content := &BaseHistoryCommand{}
	content.InitWithMap(dict)
	return content
}

/**
 *  Meta Command
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
func NewCommandForQueryMeta(did ID) MetaCommand {
	content := &BaseMetaCommand{}
	return content.Init(META, did, nil)
}

func NewCommandForRespondMeta(did ID, meta Meta) MetaCommand {
	content := &BaseMetaCommand{}
	return content.Init(META, did, meta)
}

func NewMetaCommandWithMap(dict StringKeyMap) MetaCommand {
	content := &BaseMetaCommand{}
	return content.InitWithMap(dict)
}

/**
 *  Document Command
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
func NewCommandForQueryDocuments(did ID, lastTime Time) DocumentCommand {
	content := &BaseDocumentCommand{}
	if content.Init(did, nil, nil) != nil {
		content.SetLastTime(lastTime)
	}
	return content
}

func NewCommandForRespondDocuments(did ID, meta Meta, docs []Document) DocumentCommand {
	content := &BaseDocumentCommand{}
	return content.Init(did, meta, docs)
}

func NewCommandForRespondDocument(did ID, meta Meta, document Document) DocumentCommand {
	docs := []Document{
		document,
	}
	return NewCommandForRespondDocuments(did, meta, docs)
}

func NewDocumentCommandWithMap(dict StringKeyMap) DocumentCommand {
	content := &BaseDocumentCommand{}
	return content.InitWithMap(dict)
}

/**
 *  Receipt Command
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x88),
 *      "sn"   : 456,
 *
 *      "command" : "receipt",
 *      "text"    : "...",  // text message
 *      "origin"  : {       // original message envelope
 *          "sender"    : "...",
 *          "receiver"  : "...",
 *          "time"      : 0,
 *
 *          "sn"        : 123,
 *          "signature" : "..."
 *      }
 *  }
 *  </pre></blockquote>
 */
func NewReceiptCommand(text string, head Envelope, body Content) ReceiptCommand {
	var origin = PurifyForReceipt(head, body)
	content := &BaseReceiptCommand{}
	return content.Init(text, origin)
}

func NewReceiptCommandWithMap(dict StringKeyMap) ReceiptCommand {
	content := &BaseReceiptCommand{}
	return content.InitWithMap(dict)
}

/**
 *  Group History
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x89),
 *      "sn"   : 123,
 *
 *      "command" : "reset",   // "invite", "quit", ...
 *      "time"    : 123.456,   // command timestamp
 *
 *      "group"   : "{GROUP_ID}",
 *      "members" : ["{MEMBER_ID}",]
 *  }
 *  </pre></blockquote>
 */
func NewGroupCommand(cmd string, group ID, members []ID) GroupCommand {
	content := &BaseGroupCommand{}
	return content.Init(cmd, group, members)
}

func NewGroupCommandWithMap(dict StringKeyMap) GroupCommand {
	content := &BaseGroupCommand{}
	content.InitWithMap(dict)
	return content
}

// Invite

func NewInviteCommand(group ID, members []ID) InviteCommand {
	content := &InviteGroupCommand{}
	return content.Init(group, members)
}

func NewInviteCommandWithMap(dict StringKeyMap) InviteCommand {
	content := &InviteGroupCommand{}
	content.InitWithMap(dict)
	return content
}

// Expel

func NewExpelCommand(group ID, members []ID) ExpelCommand {
	content := &ExpelGroupCommand{}
	return content.Init(group, members)
}

func NewExpelCommandWithMap(dict StringKeyMap) ExpelCommand {
	content := &ExpelGroupCommand{}
	content.InitWithMap(dict)
	return content
}

// Join

func NewJoinCommand(group ID) JoinCommand {
	content := &JoinGroupCommand{}
	return content.Init(group)
}

func NewJoinCommandWithMap(dict StringKeyMap) JoinCommand {
	content := &JoinGroupCommand{}
	content.InitWithMap(dict)
	return content
}

// Quit

func NewQuitCommand(group ID) QuitCommand {
	content := &QuitGroupCommand{}
	return content.Init(group)
}

func NewQuitCommandWithMap(dict StringKeyMap) QuitCommand {
	content := &QuitGroupCommand{}
	content.InitWithMap(dict)
	return content
}

// Reset

func NewResetCommand(group ID, members []ID) ResetCommand {
	content := &ResetGroupCommand{}
	return content.Init(group, members)
}

func NewResetCommandWithMap(dict StringKeyMap) ResetCommand {
	content := &ResetGroupCommand{}
	content.InitWithMap(dict)
	return content
}
