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
	"github.com/dimchat/core-go/dimp"
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
type Processor struct {
	dimp.Processor

	_transceiver dimp.Transceiver
}

func (processor *Processor) Init(transceiver dimp.Transceiver) *Processor {
	processor._transceiver = transceiver
	return processor
}

func (processor *Processor) Transceiver() dimp.Transceiver {
	return processor._transceiver
}

func (processor *Processor) ProcessData(data []byte) []byte {
	// 1. deserialize message
	rMsg := processor.Transceiver().DeserializeMessage(data)
	if rMsg == nil {
		// no valid message received
		return nil
	}
	// 2. process message
	rMsg = processor.Transceiver().ProcessReliableMessage(rMsg)
	if rMsg == nil {
		// nothing to respond
		return nil
	}
	// 3. serialize message
	return processor.Transceiver().SerializeMessage(rMsg)
}

func (processor *Processor) ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage {
	// NOTICE: override to check broadcast message before calling it

	// 1. verify message
	sMsg := processor.Transceiver().VerifyMessage(rMsg)
	if sMsg == nil {
		// waiting for sender's meta if not exists
		return nil
	}
	// 2. process message
	sMsg = processor.Transceiver().ProcessSecureMessage(sMsg, rMsg)
	if sMsg == nil {
		// nothing to respond
		return nil
	}
	return processor.Transceiver().SignMessage(sMsg)

	// NOTICE: override to deliver to the receiver when catch exception "receiver error ..."
}

func (processor *Processor) ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage {
	// 1. decrypt message
	iMsg := processor.Transceiver().DecryptMessage(sMsg)
	if iMsg == nil {
		// cannot decrypt this message, not for you?
		// delivering message to other receiver?
		return nil
	}
	// 2. process message
	iMsg = processor.Transceiver().ProcessInstantMessage(iMsg, rMsg)
	if iMsg == nil {
		// nothing to respond
		return nil
	}
	// 3. encrypt message
	return processor.Transceiver().EncryptMessage(iMsg)
}

func (processor *Processor) ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage {
	// 1. process content
	response := processor.Transceiver().ProcessContent(iMsg.Content(), rMsg)
	if response == nil {
		// nothing to respond
		return nil
	}

	// 2. select a local user to build message
	sender := iMsg.Sender()
	receiver := iMsg.Receiver()
	user := processor.Transceiver().SelectLocalUser(receiver)

	// 3. pack message
	env := EnvelopeCreate(user.ID(), sender, time.Time{})
	return InstantMessageCreate(env, response)
}
