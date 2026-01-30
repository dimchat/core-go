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
	"strings"

	. "github.com/dimchat/mkm-go/types"
)

/**
 *  RFC 2397
 *  ~~~~~~~~
 *  https://www.rfc-editor.org/rfc/rfc2397
 *
 *      data:[<mime type>][;charset=<charset>][;<encoding>],<encoded data>
 */
type DataURI interface {
	Head() DataHeader // "mime-type", "charset", "encoding"
	Body() string     // encoded data

	IsEmpty() bool

	// extra header value
	HeaderValue(name string) string
	Charset() string
	Filename() string

	String() string
}

/**
 *  Head of data URI
 *  ~~~~~~~~~~~~~~~~
 */
type DataHeader interface {
	MimeType() string // default is "text/plain"
	Encoding() string // default is URL Escaped Encoding (RFC 2396)

	ExtraKeys() []string
	// charset: default is "us-ascii"
	// filename: "avatar.png"
	ExtraValue(name string) string

	String() string
}

/**
 *  Split text string for data URI
 */
func ParseDataURI(uri string) DataURI {
	if uri == "" {
		return nil
	} else if !strings.HasPrefix(uri, "data:") {
		return nil
	}
	pos := strings.IndexByte(uri, ',')
	if pos < 0 {
		// data URI error
		return nil
	}
	head := splitHeader(uri, pos)
	body := uri[pos+1:]
	return NewDataURI(head, body)
}

func NewDataURI(head DataHeader, body string) DataURI {
	return &BaseURI{
		_head:      head,
		_body:      body,
		_uriString: "",
	}
}

func NewDataHeader(mimeType string, encoding string, extra StringKeyMap) DataHeader {
	return &BaseHeader{
		_mimeType:     mimeType,
		_encoding:     encoding,
		_extra:        extra,
		_headerString: "",
	}
}

// samples:
//    "data:,A%20simple%20text"
//    "data:text/html,<p>Hello, World!</p>"
//    "data:text/plain;charset=iso-8859-7,%be%fg%be"
//    "data:image/png;base64,{BASE64_ENCODE}"
//    "data:text/plain;charset=utf-8;base64,SGVsbG8sIHdvcmxkIQ=="

/**
 *  Split headers between 'data:' and first ',' from URI string
 */
func splitHeader(uri string, end int) DataHeader {
	if end < 6 {
		// header empty
		return NewDataHeader("", "", nil)
	}
	array := strings.Split(uri[5:end], ";")
	// split main info
	mimeType := ""
	encoding := ""
	// split extra info
	extra := NewMap()
	var pos int
	var name string
	var value string
	for _, item := range array {
		if item == "" {
			//panic(fmt.Sprintf("header error: %s", uri))
			continue
		}
		//
		//  2. extra info: 'charset' or 'filename'
		//
		pos = strings.IndexByte(item, '=')
		if pos >= 0 {
			name = strings.ToLower(item[:pos])
			value = item[pos+1:]
			extra[name] = value
			continue
		}
		//
		//  1. 'mime-type'
		//
		pos = strings.IndexByte(item, '/')
		if pos >= 0 {
			mimeType = item
			continue
		}
		//
		//  3. 'encoding'
		//
		encoding = item
	}
	return NewDataHeader(mimeType, encoding, extra)
}
