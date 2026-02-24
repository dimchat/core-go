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
	. "github.com/dimchat/mkm-go/types"
)

// Command name constants for system commands
//
// These values are used as the "command" field in Command message content
const (
	META      = "meta"
	DOCUMENTS = "documents"
	RECEIPT   = "receipt"
)

// Command defines the interface for system command message content
//
// Extends the base Content interface for structured command messages (e.g., meta query, document exchange)
//
//	Standard data structure: {
//	    "type"    : i2s(0x88),
//	    "sn"      : 123,
//
//	    "command" : "...",  // Command name (e.g., "meta", "documents", "receipt", ...)
//	    "extra"   : info    // Optional command-specific parameters
//	}
type Command interface {
	Content

	// CMD returns the name of the command (method / declaration)
	CMD() string
}

// CommandFactory defines the factory interface for Command
//
// Provides type-safe conversion from raw map data to concrete Command implementations
type CommandFactory interface {

	// ParseCommand converts a StringKeyMap into a Command instance
	//
	// Expects the map to follow the Command data structure format
	//
	// Parameters:
	//   - content: Raw command data in StringKeyMap format
	// Returns: Parsed Command instance (nil if parsing/validation fails)
	ParseCommand(content StringKeyMap) Command
}

//
//  Factory method
//

func ParseCommand(content interface{}) Command {
	helper := GetCommandHelper()
	return helper.ParseCommand(content)
}

func GetCommandFactory(cmd string) CommandFactory {
	helper := GetCommandHelper()
	return helper.GetCommandFactory(cmd)
}

func SetCommandFactory(cmd string, factory CommandFactory) {
	helper := GetCommandHelper()
	helper.SetCommandFactory(cmd, factory)
}
