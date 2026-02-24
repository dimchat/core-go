/* license: https://mit-license.org
 *
 *  Ming-Ke-Ming : Decentralized User Identity Authentication
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
package mkm

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// BaseVisa is the concrete implementation of the Visa interface (user profile document)
//
// Extends BaseDocument with user-specific fields (encryption public key, avatar)
// Core purpose: Authorize third-party app login and enable secure asymmetric messaging
//
// Security Note: The encryption key (visa.key) should be different from meta.key to enhance security
type BaseVisa struct {
	//Visa
	*BaseDocument

	// publicKey stores the public encryption key for secure messaging
	//
	// Used by other users to encrypt messages sent to this user (asymmetric encryption)
	// For security, this key SHOULD differ from the meta.key (signature verification key)
	publicKey EncryptKey

	// image stores the user's avatar (PNF/Portable Network File format)
	// Typically a URL
	image TransportableFile
}

func NewBaseVisa(dict StringKeyMap, data string, signature TransportableData) *BaseVisa {
	return &BaseVisa{
		BaseDocument: NewBaseDocument(dict, VISA, data, signature),
		// lazy load
		publicKey: nil,
		image:     nil,
	}
}

//-------- IVisa

// Override
func (doc *BaseVisa) Name() string {
	nickname := doc.GetProperty("name")
	return ConvertString(nickname, "")
}

// Override
func (doc *BaseVisa) SetName(nickname string) {
	doc.SetProperty("name", nickname)
}

// Override
func (doc *BaseVisa) PublicKey() EncryptKey {
	visaKey := doc.publicKey
	if visaKey == nil {
		info := doc.GetProperty("key")
		pubKey := ParsePublicKey(info)
		if encKey, ok := pubKey.(EncryptKey); ok {
			visaKey = encKey
			doc.publicKey = encKey
		}
	}
	return visaKey
}

// Override
func (doc *BaseVisa) SetPublicKey(key EncryptKey) {
	if key == nil {
		doc.SetProperty("key", nil)
	} else {
		doc.SetProperty("key", key.Map())
	}
	doc.publicKey = key
}

// Override
func (doc *BaseVisa) Avatar() TransportableFile {
	img := doc.image
	if img == nil {
		url := doc.GetProperty("avatar")
		img = ParseTransportableFile(url)
		doc.image = img
	}
	return img
}

// Override
func (doc *BaseVisa) SetAvatar(img TransportableFile) {
	if img == nil || img.IsEmpty() {
		doc.SetProperty("avatar", nil)
	} else {
		doc.SetProperty("avatar", img.Serialize())
	}
	doc.image = img
}

// BaseBulletin is the concrete implementation of the Bulletin interface (group profile document)
//
// Extends BaseDocument with core group metadata functionality
// Contains essential group profile information (name, founder, etc.)
type BaseBulletin struct {
	//Bulletin
	*BaseDocument
}

func NewBaseBulletin(dict StringKeyMap, data string, signature TransportableData) *BaseBulletin {
	return &BaseBulletin{
		BaseDocument: NewBaseDocument(dict, BULLETIN, data, signature),
	}
}

//-------- IBulletin

// Override
func (doc *BaseBulletin) Name() string {
	title := doc.GetProperty("name")
	return ConvertString(title, "")
}

// Override
func (doc *BaseBulletin) SetName(title string) {
	doc.SetProperty("name", title)
}

// Override
func (doc *BaseBulletin) Founder() ID {
	founder := doc.GetProperty("founder")
	return ParseID(founder)
}
