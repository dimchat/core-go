/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
package rfc

import (
	"fmt"
	"strings"

	. "github.com/dimchat/mkm-go/types"
)

type BaseURI struct {
	//DataURI

	head DataHeader
	body string

	uriString string // built string
}

// Override
func (uri *BaseURI) Head() DataHeader {
	return uri.head
}

// Override
func (uri *BaseURI) Body() string {
	return uri.body
}

// Override
func (uri *BaseURI) IsEmpty() bool {
	return uri.body == ""
}

// Override
func (uri *BaseURI) HeaderValue(name string) string {
	head := uri.head
	value := head.ExtraValue(name)
	if value != "" {
		// charset
		// filename
		return value
	}
	name = strings.ToLower(name)
	switch name {
	case "encoding":
		return head.Encoding()
	case "mime-type", "content-type":
		return head.MimeType()
	default:
		return ""
	}
}

// Override
func (uri *BaseURI) Charset() string {
	return uri.head.ExtraValue("charset")
}

// Override
func (uri *BaseURI) Filename() string {
	return uri.head.ExtraValue("filename")
}

// Override
func (uri *BaseURI) String() string {
	text := uri.uriString
	if text == "" {
		header := uri.head.String()
		if header == "" {
			text = fmt.Sprintf("data:,%s", uri.body)
		} else {
			text = fmt.Sprintf("data:%s,%s", header, uri.body)
		}
		uri.uriString = text
	}
	return text
}

/**
 *  Head of data URI
 *  ~~~~~~~~~~~~~~~~
 */
type BaseHeader struct {
	//DataHeader

	mimeType string
	encoding string

	extra StringKeyMap

	headerString string // built string
}

// Override
func (header *BaseHeader) MimeType() string {
	return header.mimeType
}

// Override
func (header *BaseHeader) Encoding() string {
	return header.encoding
}

// Override
func (header *BaseHeader) ExtraKeys() []string {
	extra := header.extra
	if extra == nil {
		return nil
	}
	return MapKeys(extra)
}

// Override
func (header *BaseHeader) ExtraValue(name string) string {
	extra := header.extra
	if extra == nil {
		// extra info is empty
		return ""
	} else if name == "" {
		panic("header name should not be empty")
		return ""
	}
	name = strings.ToLower(name)
	value, exists := extra[name]
	if !exists {
		// extra key not exists
		return ""
	}
	return ConvertString(value, "")
}

// Override
func (header *BaseHeader) String() string {
	text := header.headerString
	if text == "" {
		mimeType := header.mimeType
		encoding := header.encoding
		extra := header.extra
		items := make([]string, 0, 3)
		//
		//  1. 'mime-type'
		//
		if mimeType != "" {
			items = append(items, mimeType)
		} else if encoding != "" {
			// make sure 'mime-type' is the first header
			items = append(items, MIMEType.TEXT_PLAIN)
		} else if len(extra) > 0 {
			// make sure 'mime-type' is the first header
			items = append(items, MIMEType.TEXT_PLAIN)
		}
		//
		//  2. extra info: 'charset' & 'filename'
		//
		if extra != nil {
			for key, value := range extra {
				items = append(items, fmt.Sprintf("%s=%s", key, value))
			}
		}
		//
		//  3. 'encoding'
		//
		if encoding != "" {
			items = append(items, encoding)
		}
		// build header
		if len(items) > 0 {
			text = strings.Join(items, ";")
		} else {
			text = ""
		}
		header.headerString = text
	}
	return text
}
