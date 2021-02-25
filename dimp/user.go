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
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  User account for communication
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *  This class is for creating user account
 *
 *  functions:
 *      (User)
 *      1. verify(data, signature) - verify (encrypted content) data and signature
 *      2. encrypt(data)           - encrypt (symmetric key) data
 *      (LocalUser)
 *      3. sign(data)    - calculate signature of (encrypted content) data
 *      4. decrypt(data) - decrypt (symmetric key) data
 */
type User interface {
	Entity
	IUser
}
type IUser interface {

	/**
	 *  Get visa document for nickname, avatar, public key
	 *
	 * @return visa document
	 */
	Visa() Visa

	/**
	 *  Get all contacts of the user
	 *
	 * @return contact list
	 */
	Contacts() []ID

	/**
	 *  Verify data and signature with user's public keys
	 *
	 * @param data - message data
	 * @param signature - message signature
	 * @return true on correct
	 */
	Verify(data []byte, signature []byte) bool

	/**
	 *  Encrypt data, try visa.key first, if not found, use meta.key
	 *
	 * @param plaintext - message data
	 * @return encrypted data
	 */
	Encrypt(plaintext []byte) []byte

	/**
	 *  Sign data with user's private key
	 *
	 * @param data - message data
	 * @return signature
	 */
	Sign(data []byte) []byte

	/**
	 *  Decrypt data with user's private key(s)
	 *
	 * @param ciphertext - encrypted data
	 * @return plain text
	 */
	Decrypt(ciphertext []byte) []byte

	SignVisa(visa Visa) Visa
	VerifyVisa(visa Visa) bool
}

/**
 *  Base User
 *  ~~~~~~~~~
 */
type BaseUser struct {
	BaseEntity
	IUser
}

func (user *BaseUser) Init(identifier ID) *BaseUser {
	if user.BaseEntity.Init(identifier) != nil {
	}
	return user
}

//-------- IUser

func (user *BaseUser) Visa() Visa {
	doc := user.GetDocument(VISA)
	if doc != nil {
		visa, ok := doc.(Visa)
		if ok {
			return visa
		}
	}
	return nil
}

func (user *BaseUser) Contacts() []ID {
	return user.DataSource().GetContacts(user.ID())
}

func (user *BaseUser) Verify(data []byte, signature []byte) bool {
	// NOTICE: I suggest using the private key paired with meta.key to sign message
	//         so here should return the meta.key
	keys := user.DataSource().GetPublicKeysForVerification(user.ID())
	if keys != nil {
		for _, key := range keys {
			if key.Verify(data, signature) {
				// matched!
				return true
			}
		}
	}
	return false
}

func (user *BaseUser) Encrypt(plaintext []byte) []byte {
	// NOTICE: meta.key will never changed, so use visa.key to encrypt message
	//         is a better way
	key := user.DataSource().GetPublicKeyForEncryption(user.ID())
	if key == nil {
		return nil
	}
	return key.Encrypt(plaintext)
}

//
//  Interfaces for Local User
//

func (user *BaseUser) Sign(data []byte) []byte {
	// NOTICE: I suggest use the private key which paired to visa.key
	//         to sign message
	key := user.DataSource().GetPrivateKeyForSignature(user.ID())
	if key == nil {
		panic("failed to get key for signing: " + user.ID().String())
		return nil
	}
	return key.Sign(data)
}

func (user *BaseUser) Decrypt(ciphertext []byte) []byte {
	// NOTICE: if you provide a public key in visa for encryption,
	//         here you should return the private key paired with visa.key
	keys := user.DataSource().GetPrivateKeysForDecryption(user.ID())
	if keys == nil {
		panic("failed to get keys for decryption: " + user.ID().String())
		return nil
	}
	var plaintext []byte
	for _, key := range keys {
		plaintext = key.Decrypt(ciphertext)
		if plaintext != nil {
			// OK
			return plaintext
		}
	}
	// decryption failed
	return nil
}

func (user *BaseUser) SignVisa(visa Visa) Visa {
	// NOTICE: only sign visa with the private key paired with your meta.key
	if user.ID().Equal(visa.ID()) == false {
		// visa ID not match
		return nil
	}
	key := user.DataSource().GetPrivateKeyForVisaSignature(user.ID())
	if key == nil {
		panic("failed to get sign key for visa: : " + user.ID().String())
		return nil
	}
	visa.Sign(key)
	return visa
}

func (user *BaseUser) VerifyVisa(visa Visa) bool {
	// NOTICE: only verify visa with meta.key
	if user.ID().Equal(visa.ID()) == false {
		// visa ID not match
		return false
	}
	key := user.Meta().Key()
	if key == nil {
		panic("failed to get verify key for visa: : " + user.ID().String())
		return false
	}
	return visa.Verify(key)
}
