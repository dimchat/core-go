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
	"github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Delegate for Barrack
 */
type BarrackDelegate struct {
	dimp.EntityDelegate

	// memory caches
	_users map[ID]dimp.User
	_groups map[ID]dimp.Group

	_barrack IBarrack
}

func (shadow *BarrackDelegate) Init(facebook IBarrack) *BarrackDelegate {
	shadow._users = make(map[ID]dimp.User)
	shadow._groups = make(map[ID]dimp.Group)
	shadow._barrack = facebook
	return shadow
}

func (shadow *BarrackDelegate) Barrack() IBarrack {
	return shadow._barrack
}

/**
 * Call it when received 'UIApplicationDidReceiveMemoryWarningNotification',
 * this will remove 50% of cached objects
 *
 * @return number of survivors
 */
func (shadow *BarrackDelegate) ReduceMemory() int {
	finger := 0
	//finger = thanos(barrack._users, finger)
	dict1 := shadow._users
	if len(dict1) > 0 {
		index := 0
		keys := make([]ID, len(dict1))
		for key := range dict1 {
			keys[index] = key
			index++
		}
		for _, key := range keys {
			finger++
			if (finger & 1) == 1 {
				// kill it
				delete(dict1, key)
			}
			// let it go
		}
	}
	//finger = thanos(barrack._groups, finger)
	dict2 := shadow._groups
	if len(dict2) > 0 {
		index := 0
		keys := make([]ID, len(dict2))
		for key := range dict2 {
			keys[index] = key
			index++
		}
		for _, key := range keys {
			finger++
			if (finger & 1) == 1 {
				// kill it
				delete(dict2, key)
			}
			// let it go
		}
	}
	return finger >> 1
}

//func thanos(dict map[ID]Entity, finger int) int {
//	keys := keys(dict)
//	for _, key := range keys {
//		finger++
//		if (finger & 1) == 1 {
//			// kill it
//			delete(dict, key)
//		}
//		// let it go
//	}
//	return finger
//}
//
//func keys(dict map[ID]Entity) []ID {
//	index := 0
//	keys := make([]ID, len(dict))
//	for key := range dict {
//		keys[index] = key
//		index++
//	}
//	return keys
//}

func (shadow *BarrackDelegate) cacheUser(user dimp.User) {
	if user.DataSource() == nil {
		user.SetDataSource(shadow.Barrack())
	}
	shadow._users[user.ID()] = user
}

func (shadow *BarrackDelegate) cacheGroup(group dimp.Group) {
	if group.DataSource() == nil {
		group.SetDataSource(shadow.Barrack())
	}
	shadow._groups[group.ID()] = group
}

//-------- EntityDelegate

func (shadow *BarrackDelegate) SelectLocalUser(receiver ID) dimp.User {
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

func (shadow *BarrackDelegate) GetUser(identifier ID) dimp.User {
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

func (shadow *BarrackDelegate) GetGroup(identifier ID) dimp.Group {
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
