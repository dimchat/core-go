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
package core

import (
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type IBarrack interface {
	EntityDelegate

	CreateUser(identifier ID) *User
	CreateGroup(identifier ID) *Group

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return users with private key
	 */
	GetLocalUsers() []*User
}

type Barrack struct {
	IBarrack
	UserDataSource
	GroupDataSource

	// memory caches
	_users map[ID]*User
	_groups map[ID]*Group
}

func (barrack *Barrack) Init() *Barrack {
	barrack._users = make(map[ID]*User)
	barrack._groups = make(map[ID]*Group)
	return barrack
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (barrack *Barrack) ReduceMemory() int {
	finger := 0
	//finger = thanos(barrack._users, finger)
	//finger = thanos(barrack._groups, finger)
	return finger >> 1
}

func thanos(dict map[ID]interface{}, finger int) int {
	keys := keys(dict)
	for _, key := range keys {
		finger++
		if (finger & 1) == 1 {
			// kill it
			delete(dict, key)
		}
		// let it go
	}
	return finger
}

func keys(dict map[ID]interface{}) []ID {
	index := 0
	keys := make([]ID, len(dict))
	for key := range dict {
		keys[index] = key
		index++
	}
	return keys
}

func (barrack *Barrack) cacheUser(user *User) {
	if user.DataSource() == nil {
		user.SetDataSource(barrack)
	}
	barrack._users[user.ID()] = user
}

func (barrack *Barrack) cacheGroup(group *Group) {
	if group.DataSource() == nil {
		group.SetDataSource(barrack)
	}
	barrack._groups[group.ID()] = group
}

//-------- EntityDelegate

func (barrack *Barrack) SelectLocalUser(receiver ID) *User {
	users := barrack.GetLocalUsers()
	if users == nil || len(users) == 0 {
		panic("local users should not be empty")
	} else if receiver.IsBroadcast() {
		// broadcast message can decrypt by anyone, so just return current user
		return users[0]
	}
	if receiver.IsGroup() {
		// group message (recipient not designated)
		members := barrack.GetMembers(receiver)
		if members == nil || len(members) == 0 {
			// TODO: group not ready, waiting for group info
			return nil
		}
		for _, item := range users {
			if item == nil {
				continue
			}
			for _, member := range members {
				if item.ID().Equal(member) {
					// DISCUSS: set this item to be current user?
					return item
				}
			}
		}
	} else {
		// 1. personal message
		// 2. split group message
		for _, item := range users {
			if item == nil {
				continue
			}
			if item.ID().Equal(receiver) {
				// DISCUSS: set this item to be current user?
				return item
			}
		}
	}
	return nil
}

func (barrack *Barrack) GetUser(identifier ID) *User {
	// 1. get from user cache
	user := barrack._users[identifier]
	if user == nil {
		// 2. create user and cache it
		user = barrack.CreateUser(identifier)
		if user != nil {
			barrack.cacheUser(user)
		}
	}
	return user
}

func (barrack *Barrack) GetGroup(identifier ID) *Group {
	// 1. get from group cache
	// 1. get from user cache
	group := barrack._groups[identifier]
	if group == nil {
		// 2. create group and cache it
		group = barrack.CreateGroup(identifier)
		if group != nil {
			barrack.cacheGroup(group)
		}
	}
	return group
}

//-------- UserDataSource

func (barrack *Barrack) getVisaKey(user ID) EncryptKey {
	doc := barrack.GetDocument(user, VISA)
	visa, ok := doc.(Visa)
	if ok && visa.IsValid() {
		return visa.Key()
	}
	return nil
}
func (barrack *Barrack) getMetaKey(user ID) VerifyKey {
	meta := barrack.GetMeta(user)
	if meta == nil {
		return nil
	}
	return meta.Key()
}

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
	}
	//panic("failed to get encrypt key for user: " + user.String())
	return nil
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

//-------- GroupDataSource

func getIDName(group ID) string {
	name := group.Name()
	length := len(name)
	if length == 0 || (length == 8 && name == Everyone) {
		return ""
	}
	return name
}

func (barrack *Barrack) GetBroadcastFounder(group ID) ID {
	name := getIDName(group)
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
func (barrack *Barrack) GetBroadcastOwner(group ID) ID {
	name := getIDName(group)
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
func (barrack *Barrack) GetBroadcastMembers(group ID) []ID {
	name := getIDName(group)
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

func (barrack *Barrack) GetFounder(group ID) ID {
	// check broadcast group
	if group.IsBroadcast() {
		// founder of broadcast group
		return barrack.GetBroadcastFounder(group)
	}
	// check group meta
	gMeta := barrack.GetMeta(group)
	if gMeta == nil {
		// FIXME: when group profile was arrived but the meta still on the way,
		//        here will cause founder not found
		return nil
	}
	// check each member's public key with group meta
	members := barrack.GetMembers(group)
	if members != nil {
		var mMeta Meta
		for _, item := range members {
			mMeta = barrack.GetMeta(item)
			if mMeta == nil {
				// failed to get member's meta
				continue
			}
			if gMeta.MatchKey(mMeta.Key()) {
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
		return barrack.GetBroadcastOwner(group)
	}
	// check group type
	if group.Type() == POLYLOGUE {
		// Polylogue's owner is its founder
		return barrack.GetFounder(group)
	}
	// TODO: load owner from database
	return nil
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	// check broadcast group
	if group.IsBroadcast() {
		// members of broadcast group
		return barrack.GetBroadcastMembers(group)
	}
	// TODO: load members from database
	return nil
}

func (barrack *Barrack) GetAssistants(group ID) []ID {
	doc := barrack.GetDocument(group, BULLETIN)
	bulletin, ok := doc.(Bulletin)
	if ok && bulletin.IsValid() {
		return bulletin.Assistants()
	}
	// TODO: get group bots from SP configuration
	return nil
}
