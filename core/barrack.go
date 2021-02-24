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
 *  Entity Delegate
 *  ~~~~~~~~~~~~~~~
 *
 *  1. Create User/Group
 *  2. Select a local user as receiver
 */
type EntityCreator interface {

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
	EntityDataSource
	EntityFactory
	EntityCreator
}

/**
 *  Delegate for Entity
 *  ~~~~~~~~~~~~~~~~~~~
 */
type Barrack struct {
	IBarrack

	// Shadow is a delegate doing really jobs of barrack,
	// it should not be equal to the barrack itself.
	_shadow IBarrack
}

func (barrack *Barrack) Init() *Barrack {
	barrack._shadow = nil
	return barrack
}

func (barrack *Barrack) SetShadow(shadow IBarrack) {
	barrack._shadow = shadow
}
func (barrack *Barrack) Shadow() IBarrack {
	return barrack._shadow
}

//-------- EntityCreator

func (barrack *Barrack) CreateUser(identifier ID) User {
	return barrack.Shadow().CreateUser(identifier)
}

func (barrack *Barrack) CreateGroup(identifier ID) Group {
	return barrack.Shadow().CreateGroup(identifier)
}

func (barrack *Barrack) GetLocalUsers() []User {
	return barrack.Shadow().GetLocalUsers()
}

//-------- EntityFactory

func (barrack *Barrack) SelectLocalUser(receiver ID) User {
	return barrack.Shadow().SelectLocalUser(receiver)
}

func (barrack *Barrack) GetUser(identifier ID) User {
	return barrack.Shadow().GetUser(identifier)
}

func (barrack *Barrack) GetGroup(identifier ID) Group {
	return barrack.Shadow().GetGroup(identifier)
}

//-------- EntityDataSource

func (barrack *Barrack) GetMeta(identifier ID) Meta {
	return barrack.Shadow().GetMeta(identifier)
}

func (barrack *Barrack) GetDocument(identifier ID, docType string) Document {
	return barrack.Shadow().GetDocument(identifier, docType)
}

//-------- UserDataSource

func (barrack *Barrack) GetContacts(user ID) []ID {
	return barrack.Shadow().GetContacts(user)
}

func (barrack *Barrack) GetPublicKeyForEncryption(user ID) EncryptKey {
	return barrack.Shadow().GetPublicKeyForEncryption(user)
}

func (barrack *Barrack) GetPublicKeysForVerification(user ID) []VerifyKey {
	return barrack.Shadow().GetPublicKeysForVerification(user)
}

func (barrack *Barrack) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return barrack.Shadow().GetPrivateKeysForDecryption(user)
}

func (barrack *Barrack) GetPrivateKeyForSignature(user ID) SignKey {
	return barrack.Shadow().GetPrivateKeyForSignature(user)
}

func (barrack *Barrack) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return barrack.Shadow().GetPrivateKeyForVisaSignature(user)
}

//-------- GroupDataSource

func (barrack *Barrack) GetFounder(group ID) ID {
	return barrack.Shadow().GetFounder(group)
}

func (barrack *Barrack) GetOwner(group ID) ID {
	return barrack.Shadow().GetOwner(group)
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	return barrack.Shadow().GetMembers(group)
}

func (barrack *Barrack) GetAssistants(group ID) []ID {
	return barrack.Shadow().GetAssistants(group)
}
