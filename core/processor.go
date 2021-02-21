/* license: https://mit-license.org
 *
 *  DIMP : Decentralized Instant Messaging Protocol
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
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
	"time"
)

/**
 *  Message Processor Implementations
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *  Abstract method:
 *      ProcessContent(content Content, rMsg ReliableMessage) Content
 */
type MessageProcessor struct {
	Processor

	_transceiver Transceiver
}

func (processor *MessageProcessor) Init(transceiver Transceiver) *MessageProcessor {
	processor._transceiver = transceiver
	return processor
}

func (processor *MessageProcessor) Transceiver() Transceiver {
	return processor._transceiver
}

func (processor *MessageProcessor) ProcessData(data []byte) []byte {
	transceiver := processor.Transceiver()
	// 1. deserialize message
	rMsg := transceiver.DeserializeMessage(data)
	if rMsg == nil {
		// no valid message received
		return nil
	}
	// 2. process message
	rMsg = transceiver.ProcessReliableMessage(rMsg)
	if rMsg == nil {
		// nothing to respond
		return nil
	}
	// 3. serialize message
	return transceiver.SerializeMessage(rMsg)
}

func (processor *MessageProcessor) ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage {
	// NOTICE: override to check broadcast message before calling it
	transceiver := processor.Transceiver()

	// 1. verify message
	sMsg := transceiver.VerifyMessage(rMsg)
	if sMsg == nil {
		// waiting for sender's meta if not exists
		return nil
	}
	// 2. process message
	sMsg = transceiver.ProcessSecureMessage(sMsg, rMsg)
	if sMsg == nil {
		// nothing to respond
		return nil
	}
	return transceiver.SignMessage(sMsg)

	// NOTICE: override to deliver to the receiver when catch exception "receiver error ..."
}

func (processor *MessageProcessor) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage {
	transceiver := processor.Transceiver()
	// 1. decrypt message
	iMsg := transceiver.DecryptMessage(sMsg)
	if iMsg == nil {
		// cannot decrypt this message, not for you?
		// delivering message to other receiver?
		return nil
	}
	// 2. process message
	iMsg = transceiver.ProcessInstantMessage(iMsg, rMsg)
	if iMsg == nil {
		// nothing to respond
		return nil
	}
	// 3. encrypt message
	return transceiver.EncryptMessage(iMsg)
}

func (processor *MessageProcessor) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	transceiver := processor.Transceiver()
	// 1. process content
	response := transceiver.ProcessContent(iMsg.Content(), rMsg)
	if response == nil {
		// nothing to respond
		return nil
	}

	// 2. select a local user to build message
	sender := iMsg.Sender()
	receiver := iMsg.Receiver()
	user := transceiver.EntityDelegate().SelectLocalUser(receiver)

	// 3. pack message
	env := EnvelopeCreate(user.ID(), sender, time.Time{})
	return InstantMessageCreate(env, response)
}
