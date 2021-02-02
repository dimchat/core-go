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
package mkm

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Entity (User/Group)
 *  ~~~~~~~~~~~~~~~~~~~
 *  Base class of User and Group, ...
 *
 *  properties:
 *      identifier - entity ID
 *      type       - entity type
 *      meta       - meta for generate ID
 *      document   - entity document
 */
type Entity struct {
	Object

	_identifier ID

	_delegate EntityDataSource
}

func (entity *Entity) Init(identifier ID) *Entity {
	entity._identifier = identifier
	entity._delegate = nil
	return entity
}

func (entity Entity) Equal(other interface{}) bool {
	var identifier ID
	value := ObjectValue(other)
	switch value.(type) {
	case Entity:
		identifier = value.(Entity).ID()
	case ID:
		identifier = value.(ID)
	default:
		return false
	}
	// check by ID
	return entity.ID().Equal(identifier)
}

func (entity Entity) DataSource() EntityDataSource {
	return entity._delegate
}

func (entity *Entity) SetDataSource(delegate interface{}) {
	ds, ok := delegate.(EntityDataSource)
	if ok {
		entity._delegate = ds
	} else {
		panic("entity data source error")
	}
}

func (entity Entity) ID() ID {
	return entity._identifier
}

/**
 *  Get entity type
 *
 * @return ID(address) type as entity type
 */
func (entity Entity) Type() uint8 {
	return entity.ID().Type()
}

func (entity Entity) GetMeta() Meta {
	delegate := entity.DataSource()
	if delegate == nil {
		return nil
	}
	return delegate.GetMeta(entity.ID())
}

func (entity Entity) GetDocument(docType string) Document {
	delegate := entity.DataSource()
	if delegate == nil {
		return nil
	}
	return delegate.GetDocument(entity.ID(), docType)
}
