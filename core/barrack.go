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
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	"reflect"
)

/**
 *  Entity Delegate
 *  ~~~~~~~~~~~~~~~
 *
 *  1. Create User/Group
 *  2. Select a local user as receiver
 */
type EntityCreator interface {

	CreateUser(identifier ID) User
	CreateGroup(identifier ID) Group

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return users with private key
	 */
	GetLocalUsers() []User
}

type IBarrack interface {
	EntityDataSource
	EntityFactory
	EntityCreator
}

/**
 *  Delegate for Entity
 *  ~~~~~~~~~~~~~~~~~~~
 *
 * @abstract:
 *      // EntityDataSource
 *      GetMeta(identifier ID) Meta
 *      GetDocument(identifier ID, docType string) Document
 *      // UserDataSource
 *      GetContacts(user ID) []ID
 *      GetPrivateKeysForDecryption(user ID) []DecryptKey
 *      GetPrivateKeyForSignature(user ID) SignKey
 *      GetPrivateKeyForVisaSignature(user ID) SignKey
 *
 *      // EntityCreator
 *      CreateUser(identifier ID) User
 *      CreateGroup(identifier ID) Group
 *      GetLocalUsers() []User
 */
type Barrack struct {
	BaseObject
	IBarrack

	// memory caches
	_users map[ID]User
	_groups map[ID]Group
}

func (barrack *Barrack) Init() *Barrack {
	if barrack.BaseObject.Init() != nil {
		barrack._users = make(map[ID]User)
		barrack._groups = make(map[ID]Group)
	}
	return barrack
}

func (barrack *Barrack) self() IBarrack {
	return barrack.BaseObject.Self().(IBarrack)
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (barrack *Barrack) ReduceMemory() int {
	finger := 0
	finger = thanos(barrack._users, finger)
	finger = thanos(barrack._groups, finger)
	return finger >> 1
}

func thanos(planet interface{}, finger int) int {
	if reflect.TypeOf(planet).Kind() != reflect.Map {
		panic(planet)
		return finger
	}
	dict := reflect.ValueOf(planet)
	keys := dict.MapKeys()
	for _, key := range keys {
		finger++
		if (finger & 1) == 1 {
			// kill it
			dict.SetMapIndex(key, reflect.Value{})
		}
		// let it go
	}
	return finger
}

func (barrack *Barrack) cacheUser(user User) {
	if user.DataSource() == nil {
		user.SetDataSource(barrack.self())
	}
	barrack._users[user.ID()] = user
}

func (barrack *Barrack) cacheGroup(group Group) {
	if group.DataSource() == nil {
		group.SetDataSource(barrack.self())
	}
	barrack._groups[group.ID()] = group
}

//-------- EntityFactory

func (barrack *Barrack) SelectLocalUser(receiver ID) User {
	self := barrack.self()
	users := self.GetLocalUsers()
	if users == nil || len(users) == 0 {
		panic("local users should not be empty")
	} else if receiver.IsBroadcast() {
		// broadcast message can decrypt by anyone, so just return current user
		return users[0]
	}
	if receiver.IsGroup() {
		// group message (recipient not designated)
		members := self.GetMembers(receiver)
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

func (barrack *Barrack) GetUser(identifier ID) User {
	// 1. get from user cache
	user := barrack._users[identifier]
	if user == nil {
		// 2. create user and cache it
		user = barrack.self().CreateUser(identifier)
		if user != nil {
			barrack.cacheUser(user)
		}
	}
	return user
}

func (barrack *Barrack) GetGroup(identifier ID) Group {
	// 1. get from group cache
	// 1. get from user cache
	group := barrack._groups[identifier]
	if group == nil {
		// 2. create group and cache it
		group = barrack.self().CreateGroup(identifier)
		if group != nil {
			barrack.cacheGroup(group)
		}
	}
	return group
}
