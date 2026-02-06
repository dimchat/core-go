/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
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
package dkd

import (
	. "github.com/dimchat/core-go/ext"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/ext"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Message Content
 *  <p>
 *      This class is for creating message content
 *  </p>
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type"    : i2s(0),         // message type
 *      "sn"      : 0,              // serial number
 *
 *      "time"    : 123,            // message time
 *      "group"   : "{GroupID}",    // for group message
 *
 *      //-- message info
 *      "text"    : "text",         // for text message
 *      "command" : "Command Name"  // for system command
 *      //...
 *  }
 *  </pre></blockquote>
 */
type BaseContent struct {
	//Content
	Dictionary

	// message type: text, image, ...
	_type MessageType

	// serial number: random number to identify message content
	_sn SerialNumberType

	// message time
	_time Time
}

/* designated initializer */
func (content *BaseContent) InitWithMap(dict StringKeyMap) Content {
	if content.Dictionary.InitWithMap(dict) != nil {
		// lazy load
		content._type = ""
		content._sn = 0
		content._time = nil
	}
	return content
}

/* designated initializer */
func (content *BaseContent) InitWithType(msgType MessageType) Content {
	now := TimeNow()
	sn := GenerateSerialNumber(msgType, now)
	if content.Dictionary.Init() != nil {
		content._type = msgType
		content._sn = sn
		content._time = now
		content.Set("type", msgType)
		content.Set("sn", sn)
		content.SetTime("time", now)
	}
	return content
}

// Override
func (content *BaseContent) Type() MessageType {
	msgType := content._type
	if msgType == "" {
		helper := GetGeneralMessageHelper()
		msgType = helper.GetContentType(content.Map(), "")
		content._type = msgType
	}
	return msgType
}

// Override
func (content *BaseContent) SN() SerialNumberType {
	sn := content._sn
	if sn == 0 {
		sn = content.GetUInt64("sn", 0)
		content._sn = sn
	}
	return sn
}

// Override
func (content *BaseContent) Time() Time {
	when := content._time
	if when == nil {
		when = content.GetTime("time", nil)
		if when == nil {
			when = TimeNil()
		}
		content._time = when
	}
	return when
}

// Override
func (content *BaseContent) Group() ID {
	group := content.Get("group")
	return ParseID(group)
}

// Override
func (content *BaseContent) SetGroup(group ID) {
	content.SetStringer("group", group)
}

/**
 *  Base Command Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x88),
 *      "sn"   : 123,
 *
 *      "command" : "...", // command name
 *      "extra"   : info   // command parameters
 *  }
 *  </pre></blockquote>
 */
type BaseCommand struct {
	//Command
	BaseContent
}

/* designated initializer */
func (content *BaseCommand) InitWithType(msgType MessageType, cmd string) Command {
	if content.BaseContent.InitWithType(msgType) != nil {
		content.Set("command", cmd)
	}
	return content
}

func (content *BaseCommand) Init(cmd string) Command {
	return content.InitWithType(ContentType.COMMAND, cmd)
}

// Override
func (content *BaseCommand) CMD() string {
	helper := GetGeneralCommandHelper()
	return helper.GetCMD(content.Map(), "")
}
