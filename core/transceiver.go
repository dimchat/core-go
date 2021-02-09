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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

type Transceiver struct {
	dimp.Transceiver
	InstantMessageDelegate
	ReliableMessageDelegate

	_entityDelegate dimp.EntityDelegate
	_keyDelegate dimp.CipherKeyDelegate
	_transformer dimp.Transformer
	_processor dimp.Processor
	_packer dimp.Packer
}

func (transceiver *Transceiver) Init() *Transceiver {
	transceiver._entityDelegate = nil
	transceiver._keyDelegate = nil
	transceiver._transformer = nil
	transceiver._processor = nil
	transceiver._packer = nil
	return transceiver
}

/**
 *  Delegate for User/Group
 *
 * @param barrack - entity delegate
 */
func (transceiver *Transceiver) SetEntityDelegate(barrack dimp.EntityDelegate) {
	transceiver._entityDelegate = barrack
}
func (transceiver *Transceiver) EntityDelegate() dimp.EntityDelegate {
	return transceiver._entityDelegate
}

/**
 *  Delegate for Cipher Key
 *
 * @param keyCache - key store
 */
func (transceiver *Transceiver) SetCipherKeyDelegate(keyCache dimp.CipherKeyDelegate) {
	transceiver._keyDelegate = keyCache
}
func (transceiver *Transceiver) CipherKeyDelegate() dimp.CipherKeyDelegate {
	return transceiver._keyDelegate
}

/**
 *  Delegate for Transforming Message
 *
 * @param transformer - message transformer
 */
func (transceiver *Transceiver) SetTransformer(transformer dimp.Transformer) {
	transceiver._transformer = transformer
}
func (transceiver *Transceiver) Transformer() dimp.Transformer {
	return transceiver._transformer
}

/**
 *  Delegate for Processing Message
 *
 * @param processor - message processor
 */
func (transceiver *Transceiver) SetProcessor(processor dimp.Processor) {
	transceiver._processor = processor
}
func (transceiver *Transceiver) Processor() dimp.Processor {
	return transceiver._processor
}

/**
 *  Delegate for Packing Message
 *
 * @param packer - message packer
 */
func (transceiver *Transceiver) SetPacker(packer dimp.Packer) {
	transceiver._packer = packer
}
func (transceiver *Transceiver) Packer() dimp.Packer {
	return transceiver._packer
}

//
//  Interfaces for User/Group
//
func (transceiver *Transceiver) SelectLocalUser(receiver ID) dimp.User {
	return transceiver.EntityDelegate().SelectLocalUser(receiver)
}

func (transceiver *Transceiver) GetUser(identifier ID) dimp.User {
	return transceiver.EntityDelegate().GetUser(identifier)
}

func (transceiver *Transceiver) GetGroup(identifier ID) dimp.Group {
	return transceiver.EntityDelegate().GetGroup(identifier)
}

//
//  Interfaces for Cipher Key
//
func (transceiver *Transceiver) GetCipherKey(sender, receiver ID, generate bool) SymmetricKey {
	return transceiver.CipherKeyDelegate().GetCipherKey(sender, receiver, generate)
}

func (transceiver *Transceiver) CacheCipherKey(sender, receiver ID, key SymmetricKey) {
	transceiver.CipherKeyDelegate().CacheCipherKey(sender, receiver, key)
}

//
//  Interfaces for Packing Message
//
func (transceiver *Transceiver) GetOvertGroup(content Content) ID {
	return transceiver.Packer().GetOvertGroup(content)
}

func (transceiver *Transceiver) EncryptMessage(iMsg InstantMessage) SecureMessage {
	return transceiver.Packer().EncryptMessage(iMsg)
}

func (transceiver *Transceiver) SignMessage(sMsg SecureMessage) ReliableMessage {
	return transceiver.Packer().SignMessage(sMsg)
}

func (transceiver *Transceiver) SerializeMessage(rMsg ReliableMessage) []byte {
	return transceiver.Packer().SerializeMessage(rMsg)
}

func (transceiver *Transceiver) DeserializeMessage(data []byte) ReliableMessage {
	return transceiver.Packer().DeserializeMessage(data)
}

func (transceiver *Transceiver) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	return transceiver.Packer().VerifyMessage(rMsg)
}

func (transceiver *Transceiver) DecryptMessage(sMsg SecureMessage) InstantMessage {
	return transceiver.Packer().DecryptMessage(sMsg)
}

//
//  Interfaces for Processing Message
//
func (transceiver *Transceiver) ProcessData(data []byte) []byte {
	return transceiver.Processor().ProcessData(data)
}

func (transceiver *Transceiver) ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage {
	return transceiver.Processor().ProcessReliableMessage(rMsg)
}

func (transceiver *Transceiver) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage {
	return transceiver.Processor().ProcessSecureMessage(sMsg, rMsg)
}

func (transceiver *Transceiver) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	return transceiver.Processor().ProcessInstantMessage(iMsg, rMsg)
}

func (transceiver *Transceiver) ProcessContent(content Content, rMsg ReliableMessage) Content {
	return transceiver.Processor().ProcessContent(content, rMsg)
}

//-------- InstantMessageDelegate

func (transceiver *Transceiver) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().SerializeContent(content, password, iMsg)
}

func (transceiver *Transceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().EncryptContent(data, password, iMsg)
}

func (transceiver *Transceiver) EncodeData(data []byte, iMsg InstantMessage) string {
	return transceiver.Transformer().EncodeData(data, iMsg)
}

func (transceiver *Transceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().SerializeKey(password, iMsg)
}

func (transceiver *Transceiver) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	return transceiver.Transformer().EncryptKey(data, receiver, iMsg)
}

func (transceiver *Transceiver) EncodeKey(data []byte, iMsg InstantMessage) string {
	return transceiver.Transformer().EncodeKey(data, iMsg)
}

//-------- SecureMessageDelegate

func (transceiver *Transceiver) DecodeKey(key string, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecodeKey(key, sMsg)
}

func (transceiver *Transceiver) DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecryptKey(key, sender, receiver, sMsg)
}

func (transceiver *Transceiver) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	return transceiver.Transformer().DeserializeKey(key, sender, receiver, sMsg)
}

func (transceiver *Transceiver) DecodeData(data string, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecodeData(data, sMsg)
}

func (transceiver *Transceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecryptContent(data, password, sMsg)
}

func (transceiver *Transceiver) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	return transceiver.Transformer().DeserializeContent(data, password, sMsg)
}

func (transceiver *Transceiver) SignData(data []byte, sender ID, sMsg SecureMessage) []byte {
	return transceiver.Transformer().SignData(data, sender, sMsg)
}

func (transceiver *Transceiver) EncodeSignature(signature []byte, sMsg SecureMessage) string {
	return transceiver.Transformer().EncodeSignature(signature, sMsg)
}

//-------- ReliableMessageDelegate

func (transceiver *Transceiver) DecodeSignature(signature string, rMsg ReliableMessage) []byte {
	return transceiver.Transformer().DecodeSignature(signature, rMsg)
}

func (transceiver *Transceiver) VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool {
	return transceiver.Transformer().VerifyDataSignature(data, signature, sender, rMsg)
}
