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
package core

import (
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/protocol"
	"reflect"
)

/**
 *  Delegate for Barrack
 *  ~~~~~~~~~~~~~~~~~~~~
 *
 *  Abstract methods:
 *      CreateUser(identifier ID) User
 *      CreateGroup(identifier ID) Group
 *      GetLocalUsers() []User
 */
type BarrackHandler struct {
	EntityHandler

	// memory caches
	_users map[ID]User
	_groups map[ID]Group

	_barrack IBarrack
}

func (shadow *BarrackHandler) Init(facebook IBarrack) *BarrackHandler {
	shadow._users = make(map[ID]User)
	shadow._groups = make(map[ID]Group)
	shadow._barrack = facebook
	return shadow
}

func (shadow *BarrackHandler) Barrack() IBarrack {
	return shadow._barrack
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (shadow *BarrackHandler) ReduceMemory() int {
	finger := 0
	finger = thanos(shadow._users, finger)
	finger = thanos(shadow._groups, finger)
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

func (shadow *BarrackHandler) cacheUser(user User) {
	if user.DataSource() == nil {
		user.SetDataSource(shadow.Barrack())
	}
	shadow._users[user.ID()] = user
}

func (shadow *BarrackHandler) cacheGroup(group Group) {
	if group.DataSource() == nil {
		group.SetDataSource(shadow.Barrack())
	}
	shadow._groups[group.ID()] = group
}

//-------- EntityDelegate

func (shadow *BarrackHandler) SelectLocalUser(receiver ID) User {
	users := shadow.Barrack().GetLocalUsers()
	if users == nil || len(users) == 0 {
		panic("local users should not be empty")
	} else if receiver.IsBroadcast() {
		// broadcast message can decrypt by anyone, so just return current user
		return users[0]
	}
	if receiver.IsGroup() {
		// group message (recipient not designated)
		members := shadow.Barrack().GetMembers(receiver)
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

func (shadow *BarrackHandler) GetUser(identifier ID) User {
	// 1. get from user cache
	user := shadow._users[identifier]
	if user == nil {
		// 2. create user and cache it
		user = shadow.Barrack().CreateUser(identifier)
		if user != nil {
			shadow.cacheUser(user)
		}
	}
	return user
}

func (shadow *BarrackHandler) GetGroup(identifier ID) Group {
	// 1. get from group cache
	// 1. get from user cache
	group := shadow._groups[identifier]
	if group == nil {
		// 2. create group and cache it
		group = shadow.Barrack().CreateGroup(identifier)
		if group != nil {
			shadow.cacheGroup(group)
		}
	}
	return group
}
