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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/protocol"
)

/**
 *  Message Transceiver Implementations
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 */
type MessageTransceiver struct {
	Transceiver

	// shadows
	_transformer Transformer
	_processor Processor
	_packer Packer

	// delegates
	_entityFactory EntityFactory
	_keyDelegate CipherKeyDelegate
}

func (transceiver *MessageTransceiver) Init() *MessageTransceiver {
	// shadows
	transceiver._transformer = nil
	transceiver._processor = nil
	transceiver._packer = nil
	// delegates
	transceiver._entityFactory = nil
	transceiver._keyDelegate = nil
	return transceiver
}

/**
 *  Delegate for User/Group
 *
 * @param barrack - entity delegate
 */
func (transceiver *MessageTransceiver) SetEntityFactory(barrack EntityFactory) {
	transceiver._entityFactory = barrack
}
func (transceiver *MessageTransceiver) EntityFactory() EntityFactory {
	return transceiver._entityFactory
}

/**
 *  Delegate for Cipher Key
 *
 * @param keyCache - key store
 */
func (transceiver *MessageTransceiver) SetCipherKeyDelegate(keyCache CipherKeyDelegate) {
	transceiver._keyDelegate = keyCache
}
func (transceiver *MessageTransceiver) CipherKeyDelegate() CipherKeyDelegate {
	return transceiver._keyDelegate
}

/**
 *  Delegate for Transforming Message
 *
 * @param transformer - message transformer
 */
func (transceiver *MessageTransceiver) SetTransformer(transformer Transformer) {
	transceiver._transformer = transformer
}
func (transceiver *MessageTransceiver) Transformer() Transformer {
	return transceiver._transformer
}

/**
 *  Delegate for Processing Message
 *
 * @param processor - message processor
 */
func (transceiver *MessageTransceiver) SetProcessor(processor Processor) {
	transceiver._processor = processor
}
func (transceiver *MessageTransceiver) Processor() Processor {
	return transceiver._processor
}

/**
 *  Delegate for Packing Message
 *
 * @param packer - message packer
 */
func (transceiver *MessageTransceiver) SetPacker(packer Packer) {
	transceiver._packer = packer
}
func (transceiver *MessageTransceiver) Packer() Packer {
	return transceiver._packer
}

//
//  Interfaces for Packing Message
//
func (transceiver *MessageTransceiver) GetOvertGroup(content Content) ID {
	return transceiver.Packer().GetOvertGroup(content)
}

func (transceiver *MessageTransceiver) EncryptMessage(iMsg InstantMessage) SecureMessage {
	return transceiver.Packer().EncryptMessage(iMsg)
}

func (transceiver *MessageTransceiver) SignMessage(sMsg SecureMessage) ReliableMessage {
	return transceiver.Packer().SignMessage(sMsg)
}

func (transceiver *MessageTransceiver) SerializeMessage(rMsg ReliableMessage) []byte {
	return transceiver.Packer().SerializeMessage(rMsg)
}

func (transceiver *MessageTransceiver) DeserializeMessage(data []byte) ReliableMessage {
	return transceiver.Packer().DeserializeMessage(data)
}

func (transceiver *MessageTransceiver) VerifyMessage(rMsg ReliableMessage) SecureMessage {
	return transceiver.Packer().VerifyMessage(rMsg)
}

func (transceiver *MessageTransceiver) DecryptMessage(sMsg SecureMessage) InstantMessage {
	return transceiver.Packer().DecryptMessage(sMsg)
}

//
//  Interfaces for Processing Message
//
func (transceiver *MessageTransceiver) ProcessData(data []byte) []byte {
	return transceiver.Processor().ProcessData(data)
}

func (transceiver *MessageTransceiver) ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage {
	return transceiver.Processor().ProcessReliableMessage(rMsg)
}

func (transceiver *MessageTransceiver) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage {
	return transceiver.Processor().ProcessSecureMessage(sMsg, rMsg)
}

func (transceiver *MessageTransceiver) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	return transceiver.Processor().ProcessInstantMessage(iMsg, rMsg)
}

func (transceiver *MessageTransceiver) ProcessContent(content Content, rMsg ReliableMessage) Content {
	return transceiver.Processor().ProcessContent(content, rMsg)
}

//-------- InstantMessageDelegate

func (transceiver *MessageTransceiver) SerializeContent(content Content, password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().SerializeContent(content, password, iMsg)
}

func (transceiver *MessageTransceiver) EncryptContent(data []byte, password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().EncryptContent(data, password, iMsg)
}

func (transceiver *MessageTransceiver) EncodeData(data []byte, iMsg InstantMessage) string {
	return transceiver.Transformer().EncodeData(data, iMsg)
}

func (transceiver *MessageTransceiver) SerializeKey(password SymmetricKey, iMsg InstantMessage) []byte {
	return transceiver.Transformer().SerializeKey(password, iMsg)
}

func (transceiver *MessageTransceiver) EncryptKey(data []byte, receiver ID, iMsg InstantMessage) []byte {
	return transceiver.Transformer().EncryptKey(data, receiver, iMsg)
}

func (transceiver *MessageTransceiver) EncodeKey(data []byte, iMsg InstantMessage) string {
	return transceiver.Transformer().EncodeKey(data, iMsg)
}

//-------- SecureMessageDelegate

func (transceiver *MessageTransceiver) DecodeKey(key string, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecodeKey(key, sMsg)
}

func (transceiver *MessageTransceiver) DecryptKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecryptKey(key, sender, receiver, sMsg)
}

func (transceiver *MessageTransceiver) DeserializeKey(key []byte, sender ID, receiver ID, sMsg SecureMessage) SymmetricKey {
	return transceiver.Transformer().DeserializeKey(key, sender, receiver, sMsg)
}

func (transceiver *MessageTransceiver) DecodeData(data string, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecodeData(data, sMsg)
}

func (transceiver *MessageTransceiver) DecryptContent(data []byte, password SymmetricKey, sMsg SecureMessage) []byte {
	return transceiver.Transformer().DecryptContent(data, password, sMsg)
}

func (transceiver *MessageTransceiver) DeserializeContent(data []byte, password SymmetricKey, sMsg SecureMessage) Content {
	return transceiver.Transformer().DeserializeContent(data, password, sMsg)
}

func (transceiver *MessageTransceiver) SignData(data []byte, sender ID, sMsg SecureMessage) []byte {
	return transceiver.Transformer().SignData(data, sender, sMsg)
}

func (transceiver *MessageTransceiver) EncodeSignature(signature []byte, sMsg SecureMessage) string {
	return transceiver.Transformer().EncodeSignature(signature, sMsg)
}

//-------- ReliableMessageDelegate

func (transceiver *MessageTransceiver) DecodeSignature(signature string, rMsg ReliableMessage) []byte {
	return transceiver.Transformer().DecodeSignature(signature, rMsg)
}

func (transceiver *MessageTransceiver) VerifyDataSignature(data []byte, signature []byte, sender ID, rMsg ReliableMessage) bool {
	return transceiver.Transformer().VerifyDataSignature(data, signature, sender, rMsg)
}
