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
	"github.com/dimchat/mkm-go/mkm"
	. "github.com/dimchat/mkm-go/types"
	"unsafe"
)

/**
 *  Address for DIM ID
 *  ~~~~~~~~~~~~~~~~~~
 *  This class is used to build address for ID
 *
 *      properties:
 *          network - address type
 *          number  - search number
 */
type Address struct {
	String
	mkm.Address
}

/**
 *  get search number
 *
 * @return check code
 */
func (address *Address) Number() uint32 {
	return 9527
}

func (address *Address) IsUser() bool {
	addr := (*mkm.Address)(unsafe.Pointer(address))
	return mkm.AddressIsUser(addr)
}

func (address *Address) IsGroup() bool {
	addr := (*mkm.Address)(unsafe.Pointer(address))
	return mkm.AddressIsGroup(addr)
}

func (address *Address) IsBroadcast() bool {
	addr := (*mkm.Address)(unsafe.Pointer(address))
	return mkm.AddressIsBroadcast(addr)
}
