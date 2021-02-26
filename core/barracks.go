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

type BaseBarrackShadow struct {

	_barrack IBarrack
}

func (shadow *BaseBarrackShadow) Init(barrack IBarrack) *BaseBarrackShadow {
	shadow._barrack = barrack
	return shadow
}

func (shadow *BaseBarrackShadow) Barrack() IBarrack {
	return shadow._barrack
}

/**
 *  Shadow for inheritable Barrack
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 * @abstract:
 *      // EntityDataSource
 *      GetMeta(identifier ID) Meta
 *      GetDocument(identifier ID, docType string) Document
 *      // UserDataSource
 *      GetContacts(user ID) []ID
 *      GetPrivateKeysForDecryption(user ID) []DecryptKey
 *      GetPrivateKeyForSignature(user ID) SignKey
 *      GetPrivateKeyForVisaSignature(user ID) SignKey
 */
type BarrackShadow struct {
	BaseBarrackShadow

	BarrackSource
	BarrackFactory
}

func (shadow *BarrackShadow) Init(barrack IBarrack) *BarrackShadow {
	if shadow.BaseBarrackShadow.Init(barrack) != nil {
		shadow.BarrackSource.Init(barrack)
		shadow.BarrackFactory.Init(barrack)
	}
	return shadow
}

func (shadow *BarrackShadow) Barrack() IBarrack {
	return shadow.BaseBarrackShadow.Barrack()
}
