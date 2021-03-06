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
package core

import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

type ContentParser func(map[string]interface{})Content

type GeneralContentFactory struct {
	ContentFactory

	_parser ContentParser
}

func NewGeneralContentFactory(parser ContentParser) *GeneralContentFactory {
	return new(GeneralContentFactory).Init(parser)
}

func (factory *GeneralContentFactory) Init(parser ContentParser) *GeneralContentFactory {
	factory._parser = parser
	return factory
}

func (factory *GeneralContentFactory) ParseContent(content map[string]interface{}) Content {
	return factory._parser(content)
}

/**
 *  Register core content parsers
 */
func BuildContentFactories() {
	// Top-Secret
	ContentRegister(FORWARD, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(SecretContent).Init(dict)
	}))
	// Text
	ContentRegister(TEXT, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(BaseTextContent).Init(dict)
	}))

	// File
	ContentRegister(FILE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(BaseFileContent).Init(dict)
	}))
	// Image
	ContentRegister(IMAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(ImageFileContent).Init(dict)
	}))
	// Audio
	ContentRegister(AUDIO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(AudioFileContent).Init(dict)
	}))
	// Video
	ContentRegister(VIDEO, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(VideoFileContent).Init(dict)
	}))

	// Web Page
	ContentRegister(PAGE, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(WebPageContent).Init(dict)
	}))

	// Money
	ContentRegister(MONEY, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(BaseMoneyContent).Init(dict)
	}))
	ContentRegister(TRANSFER, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(TransferMoneyContent).Init(dict)
	}))

	// Command
	ContentRegister(COMMAND, NewGeneralCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseCommand).Init(dict)
	}))
	// History Command
	ContentRegister(HISTORY, NewHistoryCommandFactory(func(dict map[string]interface{}) Command {
		return new(BaseHistoryCommand).Init(dict)
	}))

	// unknown content type
	ContentRegister(0, NewGeneralContentFactory(func(dict map[string]interface{}) Content {
		return new(BaseContent).Init(dict)
	}))
}

func init() {
	BuildContentFactories()
}
