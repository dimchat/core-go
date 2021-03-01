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
	. "github.com/dimchat/mkm-go/types"
	"reflect"
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
type Entity interface {
	Object

	ID() ID
	Type() uint8

	Meta() Meta
	GetDocument(docType string) Document

	DataSource() EntityDataSource
	SetDataSource(delegate EntityDataSource)
}

/**
 *  Base Entity
 *  ~~~~~~~~~~~
 */
type BaseEntity struct {
	BaseObject
	Entity

	_identifier ID

	_delegate EntityDataSource
}

func (entity *BaseEntity) Init(identifier ID) *BaseEntity {
	if entity.BaseObject.Init() != nil {
		entity.setID(identifier)
		entity._delegate = nil
	}
	return entity
}

func (entity *BaseEntity) Release() int {
	cnt := entity.BaseObject.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		entity.setID(nil)
	}
	return cnt
}

func (entity *BaseEntity) setID(identifier ID) {
	if identifier != entity._identifier {
		ObjectRetain(identifier)
		ObjectRelease(entity._identifier)
		entity._identifier = identifier
	}
}

func (entity *BaseEntity) Equal(other interface{}) bool {
	e, ok := other.(Entity)
	if ok {
		value := reflect.ValueOf(other)
		if value.Kind() == reflect.Ptr {
			// compare pointers
			if entity == other {
				return true
			}
			other = value.Elem().Interface()
		}
		// compare values
		if *entity == other {
			return true
		}
		// compare ids
		return entity._identifier.Equal(e.ID())
	} else {
		return entity._identifier.Equal(other)
	}
}

func (entity *BaseEntity) DataSource() EntityDataSource {
	return entity._delegate
}

func (entity *BaseEntity) SetDataSource(delegate EntityDataSource) {
	entity._delegate = delegate
}

func (entity *BaseEntity) ID() ID {
	return entity._identifier
}

/**
 *  Get entity type
 *
 * @return ID(address) type as entity type
 */
func (entity *BaseEntity) Type() uint8 {
	return entity.ID().Type()
}

func (entity *BaseEntity) Meta() Meta {
	delegate := entity.DataSource()
	if delegate == nil {
		return nil
	}
	return delegate.GetMeta(entity.ID())
}

func (entity *BaseEntity) GetDocument(docType string) Document {
	delegate := entity.DataSource()
	if delegate == nil {
		return nil
	}
	return delegate.GetDocument(entity.ID(), docType)
}
