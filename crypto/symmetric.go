/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
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
package crypto

import (
	. "github.com/dimchat/mkm-go/types"
)

type BaseSymmetricKey struct {
	Dictionary
}

func (key *BaseSymmetricKey) Init(dict StringKeyMap) {
	key.Dictionary.Init(dict)
}

//// Override
//func (key *BaseSymmetricKey) Equal(other interface{}) bool {
//	if other == nil {
//		return key.IsEmpty()
//	} else if other == key {
//		// same object
//		return true
//	}
//	// check targeted value
//	target := ObjectTargetValue(other)
//	if target == nil {
//		return key.IsEmpty()
//	}
//	// check value types
//	switch v := target.(type) {
//	case SymmetricKey:
//		self := ObjectTargetValue(key).(SymmetricKey)
//		return SymmetricKeysEqual(v, self)
//	case Mapper:
//		other = v.Map()
//	case StringKeyMap:
//		other = v
//	default:
//		// other types
//		other = FetchMap(v)
//	}
//	return reflect.DeepEqual(other, key.Map())
//}

// Override
func (key *BaseSymmetricKey) Algorithm() string {
	info := key.Map()
	return GetKeyAlgorithm(info)
}

//// Override
//func (key *BaseSymmetricKey) MatchEncryptKey(pKey EncryptKey) bool {
//	self := ObjectTargetValue(key).(SymmetricKey)
//	return MatchEncryptKey(pKey, self)
//}
