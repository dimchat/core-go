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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/ext"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  User/Group Meta info
 *  <p>
 *      This class is used to generate entity ID
 *  </p>
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type"        : 1,              // algorithm version
 *      "key"         : "{public key}", // PK = secp256k1(SK);
 *      "seed"        : "moKy",         // user/group name
 *      "fingerprint" : "..."           // CT = sign(seed, SK);
 *  }
 *
 *  algorithm:
 *      fingerprint = sign(seed, SK);
 *
 *  abstract method:
 *      - GenerateAddress(network EntityType) Address
 *  </pre></blockquote>
 */
type BaseMeta struct {
	Dictionary

	/**
	 *  Meta algorithm version
	 *
	 *  <pre>
	 *  1 = MKM : username@address (default)
	 *  2 = BTC : btc_address
	 *  4 = ETH : eth_address
	 *      ...
	 *  </pre>
	 */
	_type MetaType

	/**
	 *  Public key (used for signature)
	 *  <p>
	 *      RSA / ECC
	 *  </p>
	 */
	_key VerifyKey

	/**
	 *  Seed to generate fingerprint
	 *  <p>
	 *      Username / Group-X
	 *  </p>
	 */
	_seed string

	/**
	 *  Fingerprint to verify ID and public key
	 *
	 *  <pre>
	 *  Build: fingerprint = sign(seed, privateKey)
	 *  Check: verify(seed, fingerprint, publicKey)
	 *  </pre>
	 */
	_fingerprint TransportableData

	_status int8 // 1 for valid, -1 for invalid

	// protected
	HasSeed bool
}

func (meta *BaseMeta) Init(dict StringKeyMap) Meta {
	if meta.Dictionary.Init(dict) != nil {
		// meta info from network, waiting to verify.
		meta._status = 0
		// lazy load
		meta._type = ""
		meta._key = nil
		meta._seed = ""
		meta._fingerprint = nil
	}
	return meta
}

func (meta *BaseMeta) InitWithType(version MetaType, key VerifyKey, seed string, fingerprint TransportableData) Meta {
	if meta.Dictionary.Init(NewMap()) != nil {
		// meta type
		meta.Set("type", version)
		meta._type = version
		// meta key
		meta.Set("key", key.Map())
		meta._key = key
		// seed
		if seed != "" {
			meta.Set("seed", seed)
		}
		meta._seed = seed
		// fingerprint
		if fingerprint != nil {
			meta.Set("fingerprint", fingerprint.Serialize())
		}
		meta._fingerprint = fingerprint

		// generated meta, or loaded from local storage,
		// no need to verify again.
		meta._status = 1
	}
	return meta
}

//-------- IMeta

// Override
func (meta *BaseMeta) Type() MetaType {
	version := meta._type
	if version == "" {
		helper := GetGeneralAccountHelper()
		version = helper.GetMetaType(meta.Map(), "")
		meta._type = version
	}
	return version
}

// Override
func (meta *BaseMeta) PublicKey() VerifyKey {
	key := meta._key
	if key == nil {
		info := meta.Get("key")
		key = ParsePublicKey(info)
		meta._key = key
	}
	return key
}

// Override
func (meta *BaseMeta) Seed() string {
	seed := meta._seed
	if seed == "" && meta.HasSeed {
		seed = meta.GetString("seed", "")
		meta._seed = seed
	}
	return seed
}

// Override
func (meta *BaseMeta) Fingerprint() TransportableData {
	ted := meta._fingerprint
	if ted == nil && meta.HasSeed {
		base64 := meta.Get("fingerprint")
		ted = ParseTransportableData(base64)
		meta._fingerprint = ted
	}
	return ted
}

//
//  Validation
//

// Override
func (meta *BaseMeta) IsValid() bool {
	if meta._status == 0 {
		// meta from network, try to verify
		if meta.checkValid() {
			// correct
			meta._status = 1
		} else {
			// error
			meta._status = -1
		}
	}
	return meta._status > 0
}

// protected
func (meta *BaseMeta) checkValid() bool {
	key := meta.PublicKey()
	if key == nil {
		return false
	} else if !meta.HasSeed {
		// this meta has no seed, so
		// it should not contain 'seed' or 'fingerprint'
		txt := meta.Get("seed")
		b64 := meta.Get("fingerprint")
		return txt == nil && b64 == nil
	}
	seed := meta.Seed()
	fingerprint := meta.Fingerprint()
	// check meta seed & signature
	if fingerprint == nil || fingerprint.IsEmpty() || seed == "" {
		return false
	}
	// verify fingerprint
	data := UTF8Encode(seed)
	return key.Verify(data, fingerprint.Bytes())
}

// Override
func (meta *BaseMeta) GenerateAddress(network EntityType) Address {
	panic("BaseMeta::GenerateAddress(network) > implement me!")
}
