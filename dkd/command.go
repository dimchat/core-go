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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
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
type BaseCommand struct {
	BaseContent
}

///* designated initializer */
//func (cmd *BaseCommand) Init(dict map[string]interface{}) Command {
//	if cmd.BaseContent.Init(dict) != nil {
//	}
//	return cmd
//}

/* designated initializer */
func (cmd *BaseCommand) InitWithType(msgType ContentType, command string) Command {
	if cmd.BaseContent.InitWithType(msgType) != nil {
		cmd.Set("command", command)
	}
	return cmd
}

func (cmd *BaseCommand) InitWithCommand(command string) Command {
	if cmd.InitWithType(COMMAND, command) != nil {
	}
	return cmd
}

//-------- ICommand

func (cmd *BaseCommand) CommandName() string {
	return CommandGetName(cmd.Map())
}
