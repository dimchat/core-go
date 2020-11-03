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
	. "github.com/dimchat/mkm-go/format"
	"github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
	"unsafe"
)

/**
 *  User/Group Meta data
 *  ~~~~~~~~~~~~~~~~~~~~
 *  This class is used to generate entity ID
 *
 *      data format: {
 *          version: 1,         // algorithm version
 *          key: {public key},  // PK = secp256k1(SK);
 *          seed: "moKy",       // user/group name
 *          fingerprint: "..."  // CT = sign(seed, SK);
 *      }
 *
 *      algorithm:
 *          fingerprint = sign(seed, SK);
 */
type Meta struct {
	Dictionary
	mkm.Meta

	/**
	 *  Meta algorithm version
	 *
	 *      0x01 - username@address
	 *      0x02 - btc_address
	 *      0x03 - username@btc_address
	 */
	_version MetaType

	/**
	 *  Public key (used for signature)
	 *
	 *      RSA / ECC
	 */
	_key *PublicKey

	/**
	 *  Seed to generate fingerprint
	 *
	 *      Username / Group-X
	 */
	_seed string

	/**
	 *  Fingerprint to verify ID and public key
	 *
	 *      Build: fingerprint = sign(seed, privateKey)
	 *      Check: verify(seed, fingerprint, publicKey)
	 */
	_fingerprint []byte

	_status int8 // 1 for valid, -1 for invalid
}

func (meta *Meta) Equal(other interface{}) bool {
	ptr, ok := other.(*Meta)
	if !ok {
		obj, ok := other.(Meta)
		if !ok {
			return false
		}
		ptr = &obj
	}
	meta1 := (*mkm.Meta)(unsafe.Pointer(meta))
	meta2 := (*mkm.Meta)(unsafe.Pointer(ptr))
	return mkm.MetasEqual(meta1, meta2)
}

/**
 *  Check meta valid
 *  (must call this when received a new meta from network)
 *
 * @return true on valid
 */
func (meta *Meta) IsValid() bool {
	if meta._status == 0 {
		key := meta.Key()
		if key == nil {
			// meta.key should not be empty
			meta._status = -1
		} else if MetaTypeHasSeed(meta.Type()) {
			seed := meta.Seed()
			fingerprint := meta.Fingerprint()
			if seed == "" || fingerprint == nil {
				// seed and fingerprint should not be empty
				meta._status = -1
			} else if (*key).Verify(UTF8BytesFromString(seed), fingerprint) {
				// fingerprint matched
				meta._status = 1
			} else {
				// fingerprint not matched
				meta._status = -1
			}
		} else {
			meta._status = 1
		}
	}
	return meta._status == 1
}

func (meta *Meta) Type() MetaType {
	if meta._version == 0 {
		version := meta.Get("version")
		meta._version = MetaType(version.(uint8))
	}
	return meta._version
}

func (meta *Meta) Key() *PublicKey {
	// TODO: parse key
	return meta._key
}

func (meta *Meta) Seed() string {
	if meta._seed == "" {
		if MetaTypeHasSeed(meta.Type()) {
			seed := meta.Get("seed")
			meta._seed = seed.(string)
		}
	}
	return meta._seed
}

func (meta *Meta) Fingerprint() []byte {
	if meta._fingerprint == nil {
		if MetaTypeHasSeed(meta.Type()) {
			base64 := meta.Get("fingerprint")
			meta._fingerprint = Base64Decode(base64.(string))
		}
	}
	return meta._fingerprint
}

func (meta *Meta) MatchID(identifier *mkm.ID) bool {
	ptr := (*mkm.Meta)(unsafe.Pointer(meta))
	return mkm.MetaMatchID(ptr, identifier)
}

func (meta *Meta) MatchAddress(address *mkm.Address) bool {
	ptr := (*mkm.Meta)(unsafe.Pointer(meta))
	return mkm.MetaMatchAddress(ptr, address)
}

func (meta *Meta) MatchKey(key *PublicKey) bool {
	ptr := (*mkm.Meta)(unsafe.Pointer(meta))
	return mkm.MetaMatchKey(ptr, key)
}

func (meta *Meta) GenerateID(network NetworkType) *mkm.ID {
	var name string
	if MetaTypeHasSeed(meta.Type()) {
		name = meta.Seed()
	} else {
		name = ""
	}
	address := meta.GenerateAddress(network)
	return mkm.CreateID(name, address, "")
}

func (meta *Meta) GenerateAddress(network NetworkType) *mkm.Address {
	// TODO: generate address with network ID
	return nil
}
