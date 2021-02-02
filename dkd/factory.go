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
	. "github.com/dimchat/core-go/protocol/cmd"
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

/**
 *  Register core content parsers
 */
func BuildContentFactories() {
	// Top-Secret
	ContentRegister(FORWARD, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(ForwardContent).Init(dict)
	}))
	// Text
	ContentRegister(TEXT, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(TextContent).Init(dict)
	}))

	// File
	ContentRegister(FILE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(FileContent).Init(dict)
	}))
	// Image
	ContentRegister(IMAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(ImageContent).Init(dict)
	}))
	// Audio
	ContentRegister(AUDIO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(AudioContent).Init(dict)
	}))
	// Video
	ContentRegister(VIDEO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(VideoContent).Init(dict)
	}))

	// Web Page
	ContentRegister(PAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(PageContent).Init(dict)
	}))

	// Command
	ContentRegister(COMMAND, NewCommandFactory(nil))
	// History Command
	ContentRegister(HISTORY, NewHistoryCommandFactory(nil))

	// unknown content type
	ContentRegister(0, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return NewContent(dict)
	}))
}

/**
 *  Register core command parsers
 */
func BuildCommandFactories() {
	// Meta Command
	CommandRegister(META, NewCommandFactory(func(dict map[string]interface{}) Command {
		return new(MetaCommand).Init(dict)
	}))
	// Document Command
	docParser := NewCommandFactory(func(dict map[string]interface{}) Command {
		return new(DocumentCommand).Init(dict)
	})
	CommandRegister(DOCUMENT, docParser)
	CommandRegister("profile", docParser)

	// Group Commands
	CommandRegister("group", NewGroupCommandFactory(nil))
	CommandRegister(INVITE, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(InviteCommand).Init(dict)
	}))
	CommandRegister(EXPEL, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(ExpelCommand).Init(dict)
	}))
	CommandRegister(JOIN, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(JoinCommand).Init(dict)
	}))
	CommandRegister(QUIT, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(QuitCommand).Init(dict)
	}))
	CommandRegister(QUERY, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(QueryCommand).Init(dict)
	}))
	CommandRegister(RESET, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		return new(ResetCommand).Init(dict)
	}))
}
