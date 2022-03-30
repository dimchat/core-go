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
package dimp

import (
	. "github.com/dimchat/core-go/mkm"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type IBarrack interface {
	EntityDelegate
	EntityDataSource
	UserDataSource
	GroupDataSource
}

/**
 *  Entity Manager
 *  ~~~~~~~~~~~~~~
 *  Manage meta/document for all users/groups
 *
 * @abstract:
 *      // IEntityDelegate
 *      GetUser(identifier ID) User
 *      GetGroup(identifier ID) Group
 *      // IEntityDataSource
 *      GetMeta(identifier ID) Meta
 *      GetDocument(identifier ID, docType string) Document
 *      // IUserDataSource
 *      GetContacts(user ID) []ID
 *      GetPrivateKeysForDecryption(user ID) []DecryptKey
 *      GetPrivateKeyForSignature(user ID) SignKey
 *      GetPrivateKeyForVisaSignature(user ID) SignKey
 */
type Barrack struct {

	_source EntityDataSource  // facebook
}

func (barrack *Barrack) Init() *Barrack {
	barrack._source = nil
	return barrack
}

func (barrack *Barrack) Source() EntityDataSource {
	return barrack._source
}
func (barrack *Barrack) SetSource(self EntityDataSource) {
	barrack._source = self
}

func (barrack *Barrack) getVisaKey(user ID) EncryptKey {
	self := barrack.Source()
	doc := self.GetDocument(user, VISA)
	if doc == nil || !doc.IsValid() {
		return nil
	}
	visa, ok := doc.(Visa)
	if ok {
		return visa.Key()
	} else {
		return nil
	}
}
func (barrack *Barrack) getMetaKey(user ID) VerifyKey {
	self := barrack.Source()
	meta := self.GetMeta(user)
	if meta == nil {
		return nil
	} else {
		return meta.Key()
	}
}

//-------- IEntityDelegate

//func (barrack *Barrack) GetUser(identifier ID) User {
//	panic("not implemented")
//}
//
//func (barrack *Barrack) GetGroup(identifier ID) Group {
//	panic("not implemented")
//}

//-------- IEntityDataSource

//func (barrack *Barrack) GetMeta(identifier ID) Meta {
//	panic("not implemented")
//}
//
//func (barrack *Barrack) GetDocument(identifier ID, docType string) Document {
//	panic("not implemented")
//}

//-------- IUserDataSource

//func (barrack *Barrack) GetContacts(user ID) []ID {
//	panic("not implemented")
//}

func (barrack *Barrack) GetPublicKeyForEncryption(user ID) EncryptKey {
	// 1. get key from visa
	visaKey := barrack.getVisaKey(user)
	if visaKey != nil {
		// if visa.key exists, use it for encryption
		return visaKey
	}
	// 2. get key from meta
	metaKey := barrack.getMetaKey(user)
	key, ok := metaKey.(EncryptKey)
	if ok {
		// if profile.key not exists and meta.key is encrypt key,
		// use it for encryption
		return key
	} else {
		//panic("failed to get encrypt key for user: " + user.String())
		return nil
	}
}

func (barrack *Barrack) GetPublicKeysForVerification(user ID) []VerifyKey {
	keys := make([]VerifyKey, 0, 2)
	// 1. get key from visa
	visaKey := barrack.getVisaKey(user)
	key, ok := visaKey.(VerifyKey)
	if ok {
		// the sender may use communication key to sign message.data,
		// so try to verify it with visa.key here
		keys = append(keys, key)
	}
	// 2. get key from meta
	metaKey := barrack.getMetaKey(user)
	if metaKey != nil {
		// the sender may use identity key to sign message.data,
		// try to verify it with meta.key
		keys = append(keys, key)
	}
	return keys
}

//func (barrack *Barrack) GetPrivateKeysForDecryption(user ID) []DecryptKey {
//	panic("not implemented")
//}
//
//func (barrack *Barrack) GetPrivateKeyForSignature(user ID) SignKey {
//	panic("not implemented")
//}
//
//func (barrack *Barrack) GetPrivateKeyForVisaSignature(user ID) SignKey {
//	panic("not implemented")
//}

//-------- IGroupDataSource

func (barrack *Barrack) GetFounder(group ID) ID {
	// check broadcast group
	if group.IsBroadcast() {
		// founder of broadcast group
		return BroadcastGroupGetFounder(group)
	}
	self := barrack.Source()
	// check group meta
	gMeta := self.GetMeta(group)
	if gMeta == nil {
		// FIXME: when group profile was arrived but the meta still on the way,
		//        here will cause founder not found
		return nil
	}
	// check each member's public key with group meta
	members := self.GetMembers(group)
	if members != nil {
		var mMeta Meta
		for _, item := range members {
			mMeta = self.GetMeta(item)
			if mMeta == nil {
				// failed to get member's meta
				continue
			}
			if MetaMatchKey(gMeta, mMeta.Key()) {
				// if the member's public key matches with the group's meta,
				// it means this meta was generated by the member's private key
				return item
			}
		}
	}
	// TODO: load founder from database
	return nil
}

func (barrack *Barrack) GetOwner(group ID) ID {
	// check broadcast group
	if group.IsBroadcast() {
		// owner of broadcast group
		return BroadcastGroupGetOwner(group)
	}
	// check group type
	if group.Type() == POLYLOGUE {
		// Polylogue owner is its founder
		self := barrack.Source()
		return self.GetFounder(group)
	}
	// TODO: load owner from database
	return nil
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	// check broadcast group
	if group.IsBroadcast() {
		// members of broadcast group
		return BroadcastGroupGetMembers(group)
	}
	// TODO: load members from database
	return nil
}

func (barrack *Barrack) GetAssistants(group ID) []ID {
	self := barrack.Source()
	doc := self.GetDocument(group, BULLETIN)
	if doc == nil || !doc.IsValid() {
		return nil
	}
	bulletin, ok := doc.(Bulletin)
	if ok {
		return bulletin.Assistants()
	}
	// TODO: get group bots from SP configuration
	return nil
}

//
//  Functions for Broadcast Group
//
func getGroupSeed(group ID) string {
	name := group.Name()
	length := len(name)
	if length == 0 || (length == 8 && name == Everyone) {
		return ""
	}
	return name
}

func BroadcastGroupGetFounder(group ID) ID {
	name := getGroupSeed(group)
	if name == "" {
		// Consensus: the founder of group 'everyone@everywhere'
		//            'Albert Moky'
		return FOUNDER
	} else {
		// DISCUSS: who should be the founder of group 'xxx@everywhere'?
		//          'anyone@anywhere', or 'xxx.founder@anywhere'
		return IDParse(name + ".founder@anywhere")
	}
}
func BroadcastGroupGetOwner(group ID) ID {
	name := getGroupSeed(group)
	if name == "" {
		// Consensus: the owner of group 'everyone@everywhere'
		//            'anyone@anywhere'
		return ANYONE
	} else {
		// DISCUSS: who should be the owner of group 'xxx@everywhere'?
		//          'anyone@anywhere', or 'xxx.owner@anywhere'
		return IDParse(name + ".owner@anywhere")
	}
}
func BroadcastGroupGetMembers(group ID) []ID {
	name := getGroupSeed(group)
	if name == "" {
		// Consensus: the member of group 'everyone@everywhere'
		//            'anyone@anywhere'
		return []ID{ANYONE}
	} else {
		// DISCUSS: who should be the member of group 'xxx@everywhere'?
		//          'anyone@anywhere', or 'xxx.member@anywhere'
		owner := IDParse(name + ".owner@anywhere")
		member := IDParse(name + ".member@anywhere")
		return []ID{owner, member}
	}
}
