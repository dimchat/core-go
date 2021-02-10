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
	"github.com/dimchat/core-go/dimp"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type IBarrack interface {
	dimp.EntityDataSource
	dimp.EntityDelegate

	CreateUser(identifier ID) dimp.User
	CreateGroup(identifier ID) dimp.Group

	/**
	 *  Get all local users (for decrypting received message)
	 *
	 * @return users with private key
	 */
	GetLocalUsers() []dimp.User
}

/**
 *  Delegate for entity
 *  ~~~~~~~~~~~~~~~~~~~
 *
 *  Abstract methods:
 *      CreateUser(identifier ID) dimp.User
 *      CreateGroup(identifier ID) dimp.Group
 *      GetLocalUsers() []dimp.User
 */
type Barrack struct {
	IBarrack

	_delegate dimp.EntityDelegate
	_source dimp.EntityDataSource
}

func (barrack *Barrack) Init() *Barrack {
	barrack._delegate = nil
	barrack._source = nil
	return barrack
}

func (barrack *Barrack) SetDelegate(delegate dimp.EntityDelegate) {
	barrack._delegate = delegate
}
func (barrack *Barrack) Delegate() dimp.EntityDelegate {
	return barrack._delegate
}

func (barrack *Barrack) SetDataSource(source dimp.EntityDataSource) {
	barrack._source = source
}
func (barrack *Barrack) DataSource() dimp.EntityDataSource {
	return barrack._source
}

//-------- EntityDelegate

func (barrack *Barrack) SelectLocalUser(receiver ID) dimp.User {
	return barrack.Delegate().SelectLocalUser(receiver)
}

func (barrack *Barrack) GetUser(identifier ID) dimp.User {
	return barrack.Delegate().GetUser(identifier)
}

func (barrack *Barrack) GetGroup(identifier ID) dimp.Group {
	return barrack.Delegate().GetGroup(identifier)
}

//-------- EntityDataSource

func (barrack *Barrack) GetMeta(identifier ID) Meta {
	return barrack.DataSource().GetMeta(identifier)
}

func (barrack *Barrack) GetDocument(identifier ID, docType string) Document {
	return barrack.DataSource().GetDocument(identifier, docType)
}

//-------- UserDataSource

func (barrack *Barrack) GetContacts(user ID) []ID {
	return barrack.DataSource().GetContacts(user)
}

func (barrack *Barrack) GetPublicKeyForEncryption(user ID) EncryptKey {
	return barrack.DataSource().GetPublicKeyForEncryption(user)
}

func (barrack *Barrack) GetPublicKeysForVerification(user ID) []VerifyKey {
	return barrack.DataSource().GetPublicKeysForVerification(user)
}

func (barrack *Barrack) GetPrivateKeysForDecryption(user ID) []DecryptKey {
	return barrack.DataSource().GetPrivateKeysForDecryption(user)
}

func (barrack *Barrack) GetPrivateKeyForSignature(user ID) SignKey {
	return barrack.DataSource().GetPrivateKeyForSignature(user)
}

func (barrack *Barrack) GetPrivateKeyForVisaSignature(user ID) SignKey {
	return barrack.DataSource().GetPrivateKeyForVisaSignature(user)
}

//-------- GroupDataSource

func (barrack *Barrack) GetFounder(group ID) ID {
	return barrack.DataSource().GetFounder(group)
}

func (barrack *Barrack) GetOwner(group ID) ID {
	return barrack.DataSource().GetOwner(group)
}

func (barrack *Barrack) GetMembers(group ID) []ID {
	return barrack.DataSource().GetMembers(group)
}

func (barrack *Barrack) GetAssistants(group ID) []ID {
	return barrack.DataSource().GetAssistants(group)
}
