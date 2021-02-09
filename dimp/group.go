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
	. "github.com/dimchat/mkm-go/protocol"
)

type Group interface {
	Entity

	/**
	 *  Get document for group name, assistants
	 *
	 * @return visa document
	 */
	Bulletin() Bulletin

	Founder() ID
	Owner() ID
	Members() []ID
	Assistants() []ID
}

/**
 *  Base Group
 *  ~~~~~~~~~~
 */
type BaseGroup struct {
	BaseEntity
	Group

	_founder ID
}

func (group *BaseGroup) Init(identifier ID) *BaseGroup {
	if group.BaseEntity.Init(identifier) != nil {
		// lazy load
		group._founder = nil
	}
	return group
}

func (group *BaseGroup) Equal(other interface{}) bool {
	return group.BaseEntity.Equal(other)
}

//-------- Entity

func (group *BaseGroup) DataSource() EntityDataSource {
	return group.BaseEntity.DataSource()
}

func (group *BaseGroup) SetDataSource(delegate EntityDataSource) {
	group.BaseEntity.SetDataSource(delegate)
}

func (group *BaseGroup) ID() ID {
	return group.BaseEntity.ID()
}

func (group *BaseGroup) Type() uint8 {
	return group.BaseEntity.Type()
}

func (group *BaseGroup) Meta() Meta {
	return group.BaseEntity.Meta()
}

func (group *BaseGroup) GetDocument(docType string) Document {
	return group.BaseEntity.GetDocument(docType)
}

//-------- Group

func (group *BaseGroup) Bulletin() Bulletin {
	doc := group.GetDocument(BULLETIN)
	if doc != nil {
		bulletin, ok := doc.(Bulletin)
		if ok {
			return bulletin
		}
	}
	return nil
}

func (group *BaseGroup) Founder() ID {
	if group._founder == nil {
		group._founder = group.DataSource().GetFounder(group.ID())
	}
	return group._founder
}

func (group *BaseGroup) Owner() ID {
	return group.DataSource().GetOwner(group.ID())
}

// NOTICE: the owner must be a member
//         (usually the first one)
func (group *BaseGroup) Members() []ID {
	return group.DataSource().GetMembers(group.ID())
}

func (group *BaseGroup) Assistants() []ID {
	return group.DataSource().GetAssistants(group.ID())
}
