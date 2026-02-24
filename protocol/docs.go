/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
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
package protocol

import (
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

// Visa defines the interface for user profile documents (user "visa")
//
// This document authorizes third-party applications to log in on behalf of the user
// and can generate temporary asymmetric key pairs for secure messaging.
type Visa interface {
	Document

	// Name returns the user's display name/nickname
	Name() string
	SetName(nickname string)

	// PublicKey returns the public encryption key for secure messaging
	//
	// Other users use this key to encrypt messages sent to the user (maps to "visa.key" field)
	PublicKey() EncryptKey
	SetPublicKey(publicKey EncryptKey)

	// Avatar returns the user's avatar (PNF format, typically a URL)
	Avatar() TransportableFile
	SetAvatar(img TransportableFile)
}

// Bulletin defines the interface for group profile documents (group "bulletin")
//
// Contains core metadata for group entities (name, founder, etc.)
type Bulletin interface {
	Document

	// Name returns the group's display name/title
	Name() string
	SetName(title string)

	// Founder returns the unique ID of the group's founder/creator
	//
	// This ID represents the original owner of the group
	Founder() ID
}
