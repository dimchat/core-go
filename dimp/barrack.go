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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

type Barrack struct {
	EntityDelegate
	UserDataSource
	GroupDataSource

	// memory caches
	_ids map[string]ID
	//_users map[ID]User
	//_groups map[ID]Group
}

func (barrack *Barrack) Init() *Barrack {
	barrack._ids = make(map[string]ID)
	//barrack._users = make(map[ID]User)
	//barrack._groups = make(map[ID]Group)
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
	finger = thanos(barrack._ids, finger)
	//finger = thanos(barrack._users, finger)
	//finger = thanos(barrack._groups, finger)
	return finger >> 1
}

func thanos(dictionary map[string]ID, finger int) int {
	keys := keys(dictionary)
	for _, key := range keys {
		finger++
		if (finger & 1) == 1 {
			// kill it
			delete(dictionary, key)
		}
		// let it go
	}
	return finger
}

func keys(dictionary map[string]ID) []string {
	index := 0
	keys := make([]string, len(dictionary))
	for key := range dictionary {
		keys[index] = key
		index++
	}
	return keys
}

func (barrack *Barrack) CacheID(identifier ID) ID {
	if identifier == nil {
		return nil
	}
	barrack._ids[identifier.String()] = identifier
	return identifier
}

func (barrack *Barrack) CreateID(string string) ID {
	panic("override me!")
}

func (barrack *Barrack) CreateUser(identifier ID) User {
	panic("override me!")
}

func (barrack *Barrack) CreateGroup(identifier ID) Group {
	panic("override me!")
}

func (barrack *Barrack) getID(str string) ID {
	// 1. get from ID cache
	id := barrack._ids[str]
	if id != nil {
		return id
	}
	// 2. create and cache it
	id = barrack.CreateID(str)
	return barrack.CacheID(id)
}

//-------- EntityDelegate

func (barrack *Barrack) GetID(str interface{}) ID {
	if str == nil {
		return nil
	}
	str = ObjectValue(str)
	switch str.(type) {
	case ID:
		return str.(ID)
	case string:
		return barrack.getID(str.(string))
	default:
		return nil
	}
}

func (barrack *Barrack) GetUser(identifier ID) User {
	// 1. get from user cache
	// 2. create user and cache it
	return barrack.CreateUser(identifier)
}

func (barrack *Barrack) GetGroup(identifier ID) Group {
	// 1. get from group cache
	// 2. create group and cache it
	return barrack.CreateGroup(identifier)
}

//-------- UserDataSource

func (barrack *Barrack) GetPublicKeyForEncryption(user ID) EncryptKey {
	// get profile.key
	profile := barrack.GetProfile(user)
	if profile != nil {
		key := profile.GetKey()
		if key != nil {
			// if profile.key exists,
			//     use it for encryption
			return key
		}
	}
	// get meta.key
	meta := barrack.GetMeta(user)
	if meta != nil {
		mKey := meta.Key()
		key, ok := mKey.(EncryptKey)
		if ok {
			return key
		}
	}
	return nil
}

func (barrack *Barrack) GetPublicKeysForVerification(user ID) []VerifyKey {
	keys := make([]VerifyKey, 0)
	// get profile.key
	profile := barrack.GetProfile(user)
	if profile != nil {
		pKey := profile.GetKey()
		if pKey != nil {
			key, ok := pKey.(VerifyKey)
			if ok {
				// the sender may use communication key to sign message.data,
				// so try to verify it with profile.key here
				keys = append(keys, key)
			}
		}
	}
	// get meta.key
	meta := barrack.GetMeta(user)
	if meta != nil {
		key := meta.Key()
		if key != nil {
			keys = append(keys, key)
		}
	}
	return keys
}

//-------- GroupDataSource

func (barrack *Barrack) GetFounder(group ID) ID {
	// check for broadcast
	if AddressIsBroadcast(group.Address()) {
		name := group.Name()
		if name == "" || name == "everyone" {
			// Consensus: the founder of group 'everyone@everywhere'
			//            'Albert Moky'
			return barrack.GetID("moky@anywhere")
		} else {
			// DISCUSS: who should be the founder of group 'xxx@everywhere'?
			//          'anyone@anywhere', or 'xxx.founder@anywhere'
			return barrack.GetID(name + ".founder@anywhere")
		}
	}
	return nil
}

func (barrack *Barrack) GetOwner(group ID) ID {
	// check for broadcast
	if AddressIsBroadcast(group.Address()) {
		name := group.Name()
		if name == "" || name == "everyone" {
			// Consensus: the owner of group 'everyone@everywhere'
			//            'anyone@anywhere'
			return barrack.GetID("anyone@anywhere")
		} else {
			// DISCUSS: who should be the owner of group 'xxx@everywhere'?
			//          'anyone@anywhere', or 'xxx.owner@anywhere'
			return barrack.GetID(name + ".owner@anywhere")
		}
	}
	return nil
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	// check for broadcast
	if AddressIsBroadcast(group.Address()) {
		members := make([]ID, 0)
		// add owner first
		owner := barrack.GetOwner(group)
		if owner != nil {
			members = append(members, owner)
		}
		// check and add member
		var member string
		name := group.Name()
		if name == "" || name == "everyone" {
			// Consensus: the owner of group 'everyone@everywhere'
			//            'anyone@anywhere'
			member = "anyone@anywhere"
		} else {
			// DISCUSS: who should be the owner of group 'xxx@everywhere'?
			//          'anyone@anywhere', or 'xxx.owner@anywhere'
			member = name + ".member@anywhere"
		}
		id := barrack.GetID(member)
		if id != nil && !id.Equal(owner) {
			members = append(members, id)
		}
		return members
	}
	return nil
}
