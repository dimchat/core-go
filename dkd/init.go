/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2022 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Albert Moky
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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

/**
 *  Register core content parsers
 */
func RegisterContentFactories() {
	// Top-Secret
	ContentSetFactory(FORWARD, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(SecretContent)
		content.Init(dict)
		return content
	}))
	// Text
	ContentSetFactory(TEXT, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(BaseTextContent)
		content.Init(dict)
		return content
	}))

	// File
	ContentSetFactory(FILE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(BaseFileContent)
		content.Init(dict)
		return content
	}))
	// Image
	ContentSetFactory(IMAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(ImageFileContent)
		content.Init(dict)
		return content
	}))
	// Audio
	ContentSetFactory(AUDIO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(AudioFileContent)
		content.Init(dict)
		return content
	}))
	// Video
	ContentSetFactory(VIDEO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(VideoFileContent)
		content.Init(dict)
		return content
	}))

	// Web Page
	ContentSetFactory(PAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(WebPageContent)
		content.Init(dict)
		return content
	}))

	// Money
	ContentSetFactory(MONEY, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(BaseMoneyContent)
		content.Init(dict)
		return content
	}))
	ContentSetFactory(TRANSFER, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(TransferMoneyContent)
		content.Init(dict)
		return content
	}))

	// Command
	ContentSetFactory(COMMAND, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		content := new(BaseCommand)
		content.Init(dict)
		return content
	}))
	// History Command
	ContentSetFactory(HISTORY, NewHistoryCommandFactory(func(dict map[string]interface{}) Command {
		content := new(BaseHistoryCommand)
		content.Init(dict)
		return content
	}))

	// unknown content type
	ContentSetFactory(0, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		content := new(BaseContent)
		content.Init(dict)
		return content
	}))
}

/**
 *  Register core command parsers
 */
func RegisterCommandFactories() {
	// Meta Command
	CommandSetFactory(META, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseMetaCommand)
		cmd.Init(dict)
		return cmd
	}))
	// Document Command
	CommandSetFactory(DOCUMENT, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseDocumentCommand)
		cmd.Init(dict)
		return cmd
	}))

	// Group Commands
	CommandSetFactory("group", NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(BaseGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(INVITE, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(InviteGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(EXPEL, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(ExpelGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(JOIN, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(JoinGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(QUIT, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(QuitGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(QUERY, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(QueryGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
	CommandSetFactory(RESET, NewGroupCommandFactory(func(dict map[string]interface{}) Command {
		cmd := new(ResetGroupCommand)
		cmd.Init(dict)
		return cmd
	}))
}

//func init() {
//	RegisterContentFactories()
//	RegisterCommandFactories()
//}
