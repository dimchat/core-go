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
 *  Base Message Content
 */

type BaseContent struct {
	//Content
	*Dictionary

	// message type: text, image, ...
	msgType MessageType

	// serial number: random number to identify message content
	sn SerialNumberType

	// message time
	time Time
}

func NewBaseContent(dict StringKeyMap, msgType MessageType) *BaseContent {
	if dict != nil {
		// init content with map
		return &BaseContent{
			Dictionary: NewDictionary(dict),
			// lazy load
			msgType: "",
			sn:      0,
			time:    nil,
		}
	}
	// new message content
	when := TimeNow()
	sn := GenerateSerialNumber(msgType, when)
	// prepare content info
	dict = NewMap()
	dict["type"] = msgType
	dict["sn"] = sn
	dict["time"] = TimeToFloat64(when)
	return &BaseContent{
		Dictionary: NewDictionary(nil),
		msgType:    msgType,
		sn:         sn,
		time:       when,
	}
}

// Override
func (content *BaseContent) Type() MessageType {
	msgType := content.msgType
	if msgType == "" {
		helper := GetGeneralMessageHelper()
		msgType = helper.GetContentType(content.Map(), "")
		content.msgType = msgType
	}
	return msgType
}

// Override
func (content *BaseContent) SN() SerialNumberType {
	sn := content.sn
	if sn == 0 {
		sn = content.GetUInt64("sn", 0)
		content.sn = sn
	}
	return sn
}

// Override
func (content *BaseContent) Time() Time {
	when := content.time
	if when == nil {
		when = content.GetTime("time", nil)
		if when == nil {
			when = TimeNil()
		}
		content.time = when
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
 */

type BaseCommand struct {
	//Command
	*BaseContent
}

func NewBaseCommand(dict StringKeyMap, msgType MessageType, cmd string) *BaseCommand {
	if dict != nil {
		// init command with map
		return &BaseCommand{
			BaseContent: NewBaseContent(dict, ""),
		}
	}
	// new command content
	if msgType == "" {
		msgType = ContentType.COMMAND
	}
	content := &BaseCommand{
		BaseContent: NewBaseContent(nil, msgType),
	}
	content.Set("command", cmd)
	return content
}

// Override
func (content *BaseCommand) CMD() string {
	helper := GetGeneralCommandHelper()
	return helper.GetCMD(content.Map(), "")
}
