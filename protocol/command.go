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
	. "github.com/dimchat/dkd-go/protocol"
)

const (
	//-------- command names begin --------
	META      = "meta"
	DOCUMENT  = "document"
	RECEIPT   = "receipt"
	HANDSHAKE = "handshake"
	LOGIN     = "login"
	//-------- command names end --------
)

/**
 *  Command message: {
 *      type : 0x88,
 *      sn   : 123,
 *
 *      command : "...", // command name
 *      // extra info
 *  }
 */
type Command interface {
	Content

	/**
	 *  Get command name
	 *
	 * @return command name string
	 */
	CommandName() string
}

func CommandGetName(cmd map[string]interface{}) string {
	text := cmd["command"]
	if text == nil {
		return ""
	}
	return text.(string)
}

/**
 *  Command Factory
 *  ~~~~~~~~~~~~~~~
 */
type CommandFactory interface {
	ContentFactory

	/**
	 *  Parse map object to command
	 *
	 * @param cmd - command info
	 * @return Command
	 */
	ParseCommand(cmd map[string]interface{}) Command
}

//
//  Instances of CommandFactory
//
var commandFactories = make(map[string]CommandFactory)

func CommandSetFactory(command string, factory CommandFactory) {
	commandFactories[command] = factory
}

func CommandGetFactory(command string) CommandFactory {
	return commandFactories[command]
}
