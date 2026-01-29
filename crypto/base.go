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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/types"
)

type BaseKey struct {
	Dictionary
	//CryptographyKey
}

func (key *BaseKey) Init(dict StringKeyMap) {
	key.Dictionary.Init(dict)
}

// Override
func (key *BaseKey) Algorithm() string {
	info := key.Map()
	return GetKeyAlgorithm(info)
}

//
//  Conveniences
//

func GetKeyAlgorithm(key StringKeyMap) string {
	helper := GetGeneralCryptoHelper()
	return helper.GetKeyAlgorithm(key, "")
}

func MatchEncryptKey(pKey EncryptKey, sKey DecryptKey) bool {
	return SymmetricKeysMatch(pKey, sKey)
}

func MatchSignKey(sKey SignKey, pKey VerifyKey) bool {
	return AsymmetricKeysMatch(sKey, pKey)
}

func SymmetricKeysEqual(a, b SymmetricKey) bool {
	if a == nil || b == nil {
		return a == b
	} else if a == b {
		// same object
		return true
	}
	// compare by encryption
	return MatchEncryptKey(a, b)
}

func PrivateKeysEqual(a, b PrivateKey) bool {
	if a == nil || b == nil {
		return a == b
	} else if a == b {
		// same object
		return true
	}
	// compare by signature
	return MatchSignKey(a, b.PublicKey())
}
