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
	//Meta
	*Dictionary

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
	version MetaType

	/**
	 *  Public key (used for signature)
	 *  <p>
	 *      RSA / ECC
	 *  </p>
	 */
	publicKey VerifyKey

	/**
	 *  Seed to generate fingerprint
	 *  <p>
	 *      Username / Group-X
	 *  </p>
	 */
	seed string

	/**
	 *  Fingerprint to verify ID and public key
	 *
	 *  <pre>
	 *  Build: fingerprint = sign(seed, privateKey)
	 *  Check: verify(seed, fingerprint, publicKey)
	 *  </pre>
	 */
	fingerprint TransportableData

	status int8 // 1 for valid, -1 for invalid

	// protected
	HasSeed bool
}

func NewBaseMeta(dict StringKeyMap, metaType MetaType, publicKey VerifyKey, seed string, fingerprint TransportableData) *BaseMeta {
	var status int8
	if dict != nil {
		// meta info from network, waiting to verify.
		status = 0
	} else {
		dict = NewMap()
		// meta type
		dict["type"] = metaType
		// meta key
		dict["key"] = publicKey.Map()
		// seed
		if seed != "" {
			dict["seed"] = seed
		}
		// fingerprint
		if fingerprint != nil {
			dict["fingerprint"] = fingerprint.Serialize()
		}
		// generated meta, or loaded from local storage,
		// no need to verify again.
		status = 1
	}
	return &BaseMeta{
		Dictionary:  NewDictionary(dict),
		version:     metaType,
		publicKey:   publicKey,
		seed:        seed,
		fingerprint: fingerprint,
		status:      status,
		HasSeed:     false, // set by subclass
	}
}

//-------- IMeta

// Override
func (meta *BaseMeta) Type() MetaType {
	version := meta.version
	if version == "" {
		helper := GetGeneralAccountHelper()
		version = helper.GetMetaType(meta.Map(), "")
		meta.version = version
	}
	return version
}

// Override
func (meta *BaseMeta) PublicKey() VerifyKey {
	key := meta.publicKey
	if key == nil {
		info := meta.Get("key")
		key = ParsePublicKey(info)
		meta.publicKey = key
	}
	return key
}

// Override
func (meta *BaseMeta) Seed() string {
	seed := meta.seed
	if seed == "" && meta.HasSeed {
		seed = meta.GetString("seed", "")
		meta.seed = seed
	}
	return seed
}

// Override
func (meta *BaseMeta) Fingerprint() TransportableData {
	ted := meta.fingerprint
	if ted == nil && meta.HasSeed {
		base64 := meta.Get("fingerprint")
		ted = ParseTransportableData(base64)
		meta.fingerprint = ted
	}
	return ted
}

//
//  Validation
//

// Override
func (meta *BaseMeta) IsValid() bool {
	if meta.status == 0 {
		// meta from network, try to verify
		if meta.checkValid() {
			// correct
			meta.status = 1
		} else {
			// error
			meta.status = -1
		}
	}
	return meta.status > 0
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
