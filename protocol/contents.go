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
package protocol

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// TextContent defines the interface for plain text message content
//
// Extends the base Content interface for simple text-based messages
//
//	Data structure: {
//	    "type" : i2s(0x01),
//	    "sn"   : 123,
//
//	    "text" : "..."
//	}
type TextContent interface {
	Content

	Text() string
}

// PageContent defines the interface for web page preview message content
//
// Extends the base Content interface for sharing web page information (title, URL, HTML, etc.)
//
//	Data structure: {
//	    "type"     : i2s(0x20),
//	    "sn"       : 123,
//
//	    "title"    : "...",        // Web page title
//	    "desc"     : "...",        // Short description of the web page
//	    "icon"     : "data:image/x-icon;base64,...",
//
//	    "URL"      : "https://github.com/moky/dimp",
//
//	    "HTML"     : "...",        // Optional raw HTML content of the page
//	    "mime_type": "text/html",  // MIME type of the HTML content (Content-Type)
//	    "encoding" : "utf8",       // Character encoding of the HTML content
//	    "base"     : "about:blank" // Base URL for resolving relative links in HTML
//	}
type PageContent interface {
	Content

	// Title returns the web page title
	Title() string
	SetTitle(title string)

	// Icon returns the web page icon (PNF format)
	Icon() TransportableFile
	SetIcon(img TransportableFile)

	// Description returns the short description of the web page (from "desc" field)
	Description() string
	SetDescription(desc string)

	// URL returns the full URL of the web page
	URL() URL
	SetURL(url URL)

	// HTML returns the raw HTML content of the web page
	HTML() string
	SetHTML(html string)
}

// NameCard defines the interface for contact name card message content
//
// Extends the base Content interface for sharing contact information (did, name, avatar)
//
//	Data structure: {
//	    "type"   : i2s(0x33),
//	    "sn"     : 123,
//
//	    "did"    : "{ID}",       // Contact's entity ID
//	    "name"   : "{nickname}", // Contact's display name/nickname
//	    "avatar" : "{URL}",      // Contact's avatar (PNF format, typically a URL)
//	    ...                     // Additional contact info (optional)
//	}
type NameCard interface {
	Content

	// ID returns the contact's entity ID (from "did" field)
	ID() ID

	// Name returns the contact's display name/nickname
	Name() string

	// Avatar returns the contact's avatar (PNF format)
	Avatar() TransportableFile
}
