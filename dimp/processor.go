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
package dimp

import (
	. "github.com/dimchat/dkd-go/protocol"
)

/**
 *  Message Processor
 *  ~~~~~~~~~~~~~~~~~
 */
type Processor interface {

	/**
	 *  Process data package
	 *
	 * @param data - data to be processed
	 * @return response data
	 */
	ProcessData(data []byte) []byte

	/**
	 *  Process network message
	 *
	 * @param rMsg - message to be processed
	 * @return response message
	 */
	ProcessReliableMessage(rMsg ReliableMessage) ReliableMessage

	/**
	 *  Process encrypted message
	 *
	 * @param sMsg - message to be processed
	 * @param rMsg - message received
	 * @return response message
	 */
	ProcessSecureMessage(sMsg SecureMessage, rMsg ReliableMessage) SecureMessage

	/**
	 *  Process plain message
	 *
	 * @param iMsg - message to be processed
	 * @param rMsg - message received
	 * @return response message
	 */
	ProcessInstantMessage(iMsg InstantMessage, rMsg ReliableMessage) InstantMessage

	/**
	 *  Process message content
	 *
	 * @param content - content to be processed
	 * @param rMsg - message received
	 * @return response content
	 */
	ProcessContent(content Content, rMsg ReliableMessage) Content
}
