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
package core

import (
	. "github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Entity Handler
 *  ~~~~~~~~~~~~~~
 */
type EntityHandler interface {
	EntityDelegate

	CreateUser(identifier ID) User
	CreateGroup(identifier ID) Group

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return users with private key
	 */
	GetLocalUsers() []User
}

type IBarrack interface {
	EntityHandler
	EntityDataSource
}

/**
 *  Delegate for Entity
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Barrack struct {
	IBarrack

	_handler EntityHandler
	_source EntityDataSource
}

func (barrack *Barrack) Init() *Barrack {
	barrack._handler = nil
	barrack._source = nil
	return barrack
}

func (barrack *Barrack) SetHandler(handler EntityHandler) {
	barrack._handler = handler
}
func (barrack *Barrack) Handler() EntityHandler {
	return barrack._handler
}

func (barrack *Barrack) SetSource(source EntityDataSource) {
	barrack._source = source
}
func (barrack *Barrack) Source() EntityDataSource {
	return barrack._source
}

//-------- EntityHandler

func (barrack *Barrack) CreateUser(identifier ID) User {
	return barrack.Handler().CreateUser(identifier)
}

func (barrack *Barrack) CreateGroup(identifier ID) Group {
	return barrack.Handler().CreateGroup(identifier)
}

func (barrack *Barrack) GetLocalUsers() []User {
	return barrack.Handler().GetLocalUsers()
}

//-------- EntityDelegate

func (barrack *Barrack) SelectLocalUser(receiver ID) User {
	return barrack.Handler().SelectLocalUser(receiver)
}

func (barrack *Barrack) GetUser(identifier ID) User {
	return barrack.Handler().GetUser(identifier)
}

func (barrack *Barrack) GetGroup(identifier ID) Group {
	return barrack.Handler().GetGroup(identifier)
}

//-------- EntityDataSource

func (barrack *Barrack) GetMeta(identifier ID) Meta {
	return barrack.Source().GetMeta(identifier)
}

func (barrack *Barrack) GetDocument(identifier ID, docType string) Document {
	return barrack.Source().GetDocument(identifier, docType)
}

//-------- UserDataSource

func (barrack *Barrack) GetContacts(user ID) []ID {
	return barrack.Source().GetContacts(user)
}

func (barrack *Barrack) GetPublicKeyForEncryption(user ID) EncryptKey {
	return barrack.Source().GetPublicKeyForEncryption(user)
}

func (barrack *Barrack) GetPublicKeysForVerification(user ID) []VerifyKey {
	return barrack.Source().GetPublicKeysForVerification(user)
}

func (barrack *Barrack) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return barrack.Source().GetPrivateKeysForDecryption(user)
}

func (barrack *Barrack) GetPrivateKeyForSignature(user ID) SignKey {
	return barrack.Source().GetPrivateKeyForSignature(user)
}

func (barrack *Barrack) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return barrack.Source().GetPrivateKeyForVisaSignature(user)
}

//-------- GroupDataSource

func (barrack *Barrack) GetFounder(group ID) ID {
	return barrack.Source().GetFounder(group)
}

func (barrack *Barrack) GetOwner(group ID) ID {
	return barrack.Source().GetOwner(group)
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	return barrack.Source().GetMembers(group)
}

func (barrack *Barrack) GetAssistants(group ID) []ID {
	return barrack.Source().GetAssistants(group)
}
