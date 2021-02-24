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
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  User Data Source
 *  ~~~~~~~~~~~~~~~~
 *
 *  (Encryption/decryption)
 *  1. public key for encryption
 *     if visa.key not exists, means it is the same key with meta.key
 *  2. private keys for decryption
 *     the private keys paired with [visa.key, meta.key]
 *
 *  (Signature/Verification)
 *  3. private key for signature
 *     the private key paired with visa.key or meta.key
 *  4. public keys for verification
 *     [visa.key, meta.key]
 *
 *  (Visa Document)
 *  5. private key for visa signature
 *     the private key pared with meta.key
 *  6. public key for visa verification
 *     meta.key only
 */
type UserDataSource interface {

	/**
	 *  Get contacts list
	 *
	 * @param user - user ID
	 * @return contacts list (ID)
	 */
	GetContacts(user ID) []ID

	/**
	 *  Get user's public key for encryption
	 *  (visa.key or meta.key)
	 *
	 * @param user - user ID
	 * @return visa.key or meta.key
	 */
	GetPublicKeyForEncryption(user ID) EncryptKey

	/**
	 *  Get user's public keys for verification
	 *  [visa.key, meta.key]
	 *
	 * @param user - user ID
	 * @return public keys
	 */
	GetPublicKeysForVerification(user ID) []VerifyKey

	/**
	 *  Get user's private keys for decryption
	 *  (which paired with [visa.key, meta.key])
	 *
	 * @param user - user ID
	 * @return private keys
	 */
	GetPrivateKeysForDecryption(user ID) []DecryptKey

	/**
	 *  Get user's private key for signature
	 *  (which paired with visa.key or meta.key)
	 *
	 * @param user - user ID
	 * @return private key
	 */
	GetPrivateKeyForSignature(user ID) SignKey

	/**
	 *  Get user's private key for signing visa
	 *
	 * @param user - user ID
	 * @return private key
	 */
	GetPrivateKeyForVisaSignature(user ID) SignKey
}

/**
 *  Group Data Source
 *  ~~~~~~~~~~~~~~~~~
 */
type GroupDataSource interface {

	/**
	 *  Get group founder
	 *
	 * @param group - group ID
	 * @return fonder ID
	 */
	GetFounder(group ID) ID

	/**
	 *  Get group owner
	 *
	 * @param group - group ID
	 * @return owner ID
	 */
	GetOwner(group ID) ID

	/**
	 *  Get group members list
	 *
	 * @param group - group ID
	 * @return members list (ID)
	 */
	GetMembers(group ID) []ID

	/**
	 *  Get assistants for this group
	 *
	 * @param group - group ID
	 * @return robot ID list
	 */
	GetAssistants(group ID) []ID
}

/**
 *  Entity Data Source
 *  ~~~~~~~~~~~~~~~~~~
 */
type EntityDataSource interface {
	UserDataSource
	GroupDataSource

	/**
	 *  Get meta for entity ID
	 *
	 * @param identifier - entity ID
	 * @return meta object
	 */
	GetMeta(identifier ID) Meta

	/**
	 *  Get document for entity ID
	 *
	 * @param identifier - entity ID
	 * @param docType - document type
	 * @return document object
	 */
	GetDocument(identifier ID, docType string) Document
}

/**
 *  Entity Delegate
 *  ~~~~~~~~~~~~~~~
 *
 *  1. Create User/Group
 *  2. Select a local user as receiver
 */
type EntityFactory interface {

	/**
	 *  Select local user for receiver
	 *
	 * @param receiver - user/group ID
	 * @return local user
	 */
	SelectLocalUser(receiver ID) User

	/**
	 *  Create user with ID
	 *
	 * @param identifier - user ID
	 * @return user
	 */
	GetUser(identifier ID) User

	/**
	 *  Create group with ID
	 *
	 * @param identifier - group ID
	 * @return group
	 */
	GetGroup(identifier ID) Group
}

/**
 *  Cipher Key Delegate
 *
 *  1. get symmetric key for sending message;
 *  2. cache symmetric key for reusing.
 */
type CipherKeyDelegate interface {

	/**
	 *  Get cipher key for encrypt message from 'sender' to 'receiver'
	 *
	 * @param sender - from where (user or contact ID)
	 * @param receiver - to where (contact or user/group ID)
	 * @param generate - generate when key not exists
	 * @return cipher key
	 */
	GetCipherKey(sender, receiver ID, generate bool) SymmetricKey

	/**
	 *  Cache cipher key for reusing, with the direction (from 'sender' to 'receiver')
	 *
	 * @param sender - from where (user or contact ID)
	 * @param receiver - to where (contact or user/group ID)
	 * @param key - cipher key
	 */
	CacheCipherKey(sender, receiver ID, key SymmetricKey)
}
