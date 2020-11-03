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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
)

const (
	//-------- command names begin --------
	META      = "meta"
	PROFILE   = "profile"
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
 *      extra   : info   // command parameters
 *  }
 */
type Command struct {
	Content
}

func (cmd *Command) Init(dictionary map[string]interface{}) *Command {
	if cmd.Content.Init(dictionary) != nil {
		// init
	}
	return cmd
}

func (cmd *Command) InitWithType(t ContentType, command string) *Command {
	if cmd.Content.InitWithType(t) != nil {
		cmd.SetCommand(command)
	}
	return cmd
}

func (cmd *Command) InitWithCommand(command string) *Command {
	return cmd.InitWithType(COMMAND, command)
}

//-------- setter/getter --------

/**
 *  Get command name
 *
 * @return command name string
 */
func (cmd *Command) GetCommand() string {
	command := cmd.Get("command")
	return command.(string)
}

func (cmd *Command) SetCommand(command string) {
	cmd.Set("command", command)
}
