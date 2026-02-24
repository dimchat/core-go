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
	return NewBaseCommand(dict, "", "")
}

func NewHistoryCommandWithMap(dict StringKeyMap) Command {
	return NewBaseHistoryCommand(dict, "")
}

/**
 *  Meta Command
 */

func NewCommandForQueryMeta(did ID) MetaCommand {
	return NewBaseMetaCommand(nil, META, did, nil)
}

func NewCommandForRespondMeta(did ID, meta Meta) MetaCommand {
	return NewBaseMetaCommand(nil, META, did, meta)
}

func NewMetaCommandWithMap(dict StringKeyMap) Command {
	return NewBaseMetaCommand(dict, "", nil, nil)
}

/**
 *  Document Command
 */

func NewCommandForQueryDocuments(did ID, lastTime Time) DocumentCommand {
	return NewBaseDocumentCommand(nil, did, nil, nil, lastTime)
}

func NewCommandForRespondDocuments(did ID, meta Meta, docs []Document) DocumentCommand {
	return NewBaseDocumentCommand(nil, did, meta, docs, nil)
}

func NewCommandForRespondDocument(did ID, meta Meta, document Document) DocumentCommand {
	docs := []Document{document}
	return NewBaseDocumentCommand(nil, did, meta, docs, nil)
}

func NewDocumentCommandWithMap(dict StringKeyMap) Command {
	return NewBaseDocumentCommand(dict, nil, nil, nil, nil)
}

/**
 *  Receipt Command
 */

func NewReceiptCommand(text string, head Envelope, body Content) ReceiptCommand {
	origin := PurifyForReceipt(head, body)
	return NewBaseReceiptCommand(nil, text, origin)
}

func NewReceiptCommandWithMap(dict StringKeyMap) Command {
	return NewBaseReceiptCommand(dict, "", nil)
}

/**
 *  Group History
 */

func NewGroupCommand(cmd string, group ID, members []ID) GroupCommand {
	return NewBaseGroupCommand(nil, cmd, group, members)
}

func NewGroupCommandWithMap(dict StringKeyMap) Command {
	return NewBaseGroupCommand(dict, "", nil, nil)
}

// Invite

func NewInviteCommand(group ID, members []ID) InviteCommand {
	return NewInviteGroupCommand(nil, group, members)
}

func NewInviteCommandWithMap(dict StringKeyMap) Command {
	return NewInviteGroupCommand(dict, nil, nil)
}

// Expel

func NewExpelCommand(group ID, members []ID) ExpelCommand {
	return NewExpelGroupCommand(nil, group, members)
}

func NewExpelCommandWithMap(dict StringKeyMap) Command {
	return NewExpelGroupCommand(dict, nil, nil)
}

// Join

func NewJoinCommand(group ID) JoinCommand {
	return NewJoinGroupCommand(nil, group)
}

func NewJoinCommandWithMap(dict StringKeyMap) Command {
	return NewJoinGroupCommand(dict, nil)
}

// Quit

func NewQuitCommand(group ID) QuitCommand {
	return NewQuitGroupCommand(nil, group)
}

func NewQuitCommandWithMap(dict StringKeyMap) Command {
	return NewQuitGroupCommand(dict, nil)
}

// Reset

func NewResetCommand(group ID, members []ID) ResetCommand {
	return NewResetGroupCommand(nil, group, members)
}

func NewResetCommandWithMap(dict StringKeyMap) Command {
	return NewResetGroupCommand(dict, nil, nil)
}
