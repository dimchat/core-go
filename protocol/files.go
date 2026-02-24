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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// FileContent defines the base interface for file attachment message content
//
// Extends the base Content interface for file transfer (supports direct data or CDN URL)
// Files uploaded to public CDNs should be encrypted with a symmetric key for security
//
//	Data structure: {
//	    "type"     : i2s(0x10),
//	    "sn"       : 123,
//
//	    "data"     : "...",        // Base64-encoded file content (direct file data)
//	    "filename" : "photo.png",  // Name of the file (including extension)
//
//	    "URL"      : "http://...", // CDN download URL (alternative to direct data)
//	    "key"      : {             // Symmetric key to decrypt CDN-hosted file data (optional)
//	        "algorithm" : "AES",   // Encryption algorithm (e.g., DES, AES)
//	        "data"      : "{BASE64_ENCODE}",
//	        ...
//	    }
//	}
type FileContent interface {
	Content

	// Data returns the direct file content (Base64-encoded TransportableData)
	Data() TransportableData
	SetData(data TransportableData)

	// Filename returns the name of the file (including extension)
	Filename() string
	SetFilename(filename string)

	// URL returns the CDN download URL (alternative to direct file data)
	URL() URL
	SetURL(url URL)

	// Password returns the symmetric decrypt key for CDN-hosted encrypted files
	// (maps to the "key" field in the data structure)
	Password() DecryptKey
	SetPassword(key DecryptKey)
}

// ImageContent defines the interface for image file message content
//
// Extends FileContent with thumbnail support for image previews
//
//	Data structure: {
//	    "type"     : i2s(0x12),
//	    "sn"       : 123,
//
//	    "data"     : "...",        // Base64-encoded image content
//	    "filename" : "photo.png",  // Image file name
//
//		"URL"      : "http://...", // CDN download URL (optional)
//		"key"      : {             // Symmetric decrypt key for CDN image (optional)
//		    "algorithm" : "AES",    // Encryption algorithm
//		    "data"      : "{BASE64_ENCODE}"
//		},
//		"thumbnail": "data:image/jpeg;base64,..."
//	}
type ImageContent interface {
	FileContent

	// Thumbnail returns the image thumbnail (PNF format)
	Thumbnail() TransportableFile
	SetThumbnail(img TransportableFile)
}

// AudioContent defines the interface for audio file message content
//
// Extends FileContent with duration and speech-to-text support
//
//	Data structure: {
//	    "type"     : i2s(0x14),
//	    "sn"       : 123,
//
//	    "data"     : "...",        // Base64-encoded audio content
//	    "filename" : "voice.mp4",  // Audio file name
//
//	    "URL"      : "http://...", // CDN download URL (optional)
//	    "key"      : {             // Symmetric decrypt key for CDN audio (optional)
//	        "algorithm" : "AES",   // Encryption algorithm
//	        "data"      : "{BASE64_ENCODE}"
//	    },
//	    "duration" : 123.45,       // Audio duration in seconds (float)
//	    "text"     : "..."         // Transcribed text (Automatic Speech Recognition/ASR)
//	}
type AudioContent interface {
	FileContent

	// Duration returns the audio duration in seconds
	Duration() float64
	SetDuration(duration float64)

	// Text returns the transcribed text from the audio (ASR result)
	Text() string
	SetText(text string)
}

// VideoContent defines the interface for video file message content
//
// Extends FileContent with snapshot support for video previews
//
//	Data structure: {
//	    "type"     : i2s(0x16),
//	    "sn"       : 123,
//
//	    "data"     : "...",        // Base64-encoded video content
//	    "filename" : "movie.mp4",  // Video file name
//
//	    "URL"      : "http://...", // CDN download URL (optional)
//	    "key"      : {             // Symmetric decrypt key for CDN video (optional)
//	        "algorithm" : "AES",   // Encryption algorithm
//	        "data"      : "{BASE64_ENCODE}"
//	    },
//	    "snapshot" : "data:image/jpeg;base64,..."
//	}
type VideoContent interface {
	FileContent

	// Snapshot returns the video snapshot/thumbnail (PNF format)
	Snapshot() TransportableFile
	SetSnapshot(img TransportableFile)
}
