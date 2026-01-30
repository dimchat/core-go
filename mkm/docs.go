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

/**
 *  User Document
 *  ~~~~~~~~~~~~~
 *
 *  This interface is defined for authorizing other apps to login,
 *  which can generate a temporary asymmetric key pair for messaging.
 */
type BaseVisa struct {
	BaseDocument

	// Public Key for encryption
	// ~~~~~~~~~~~~~~~~~~~~~~~~~
	// For safety considerations, the visa.key which used to encrypt message data
	// should be different with meta.key
	_key EncryptKey

	// Avatar URL
	_image TransportableFile
}

func (doc *BaseVisa) Init(dict StringKeyMap) Visa {
	if doc.BaseDocument.Init(dict) != nil {
		// lazy load
		doc._key = nil
		doc._image = nil
	}
	return doc
}

// load visa info from local database
func (doc *BaseVisa) InitWithData(data string, signature TransportableData) Visa {
	if doc.BaseDocument.InitWithType(VISA, data, signature) != nil {
		// lazy load
		doc._key = nil
		doc._image = nil
	}
	return doc
}

// create empty visa
func (doc *BaseVisa) InitEmptyVisa() Visa {
	if doc.BaseDocument.InitEmptyDocument(VISA) != nil {
		doc._key = nil
		doc._image = nil
	}
	return doc
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
	key := doc._key
	if key == nil {
		info := doc.GetProperty("key")
		pubKey := ParsePublicKey(info)
		if pubKey == nil {
			//panic("public key error")
			return nil
		}
		encKey, ok := pubKey.(EncryptKey)
		if ok {
			key = encKey
			doc._key = encKey
		}
	}
	return key
}

// Override
func (doc *BaseVisa) SetPublicKey(key EncryptKey) {
	if key == nil {
		doc.SetProperty("key", nil)
	} else {
		doc.SetProperty("key", key.Map())
	}
	doc._key = key
}

// Override
func (doc *BaseVisa) Avatar() TransportableFile {
	img := doc._image
	if img == nil {
		url := doc.GetProperty("avatar")
		img = ParseTransportableFile(url)
		doc._image = img
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
	doc._image = img
}

/**
 *  Group Document
 *  ~~~~~~~~~~~~~~
 */
type BaseBulletin struct {
	BaseDocument
}

func (doc *BaseBulletin) Init(dict StringKeyMap) Bulletin {
	if doc.BaseDocument.Init(dict) != nil {
	}
	return doc
}

func (doc *BaseBulletin) InitWithData(data string, signature TransportableData) Bulletin {
	if doc.BaseDocument.InitWithType(BULLETIN, data, signature) != nil {
	}
	return doc
}

func (doc *BaseBulletin) InitEmptyBulletin() Bulletin {
	if doc.BaseDocument.InitEmptyDocument(BULLETIN) != nil {
	}
	return doc
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
