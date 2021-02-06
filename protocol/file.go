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
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
)

/**
 *  File message: {
 *      type : 0x10,
 *      sn   : 123,
 *
 *      URL      : "http://", // upload to CDN
 *      data     : "...",     // if (!URL) base64_encode(fileContent)
 *      filename : "..."
 *  }
 */
type FileContent struct {
	BaseContent

	_data []byte       // file data (plaintext)
	_key SymmetricKey  // symmetric key to decrypt the encrypted data from URL
}

func (content *FileContent) Init(dict map[string]interface{}) *FileContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._data = nil
		content._key = nil
	}
	return content
}

func (content *FileContent) InitWithType(msgType uint8, filename string, data []byte) *FileContent {
	if msgType == 0 {
		msgType = FILE
	}
	if content.BaseContent.InitWithType(msgType) != nil {
		content.SetFilename(filename)
		content.SetData(data)
		content._key = nil
	}
	return content
}

//-------- setter/getter --------

func (content *FileContent) GetURL() string {
	url := content.Get("URL")
	if url == nil {
		return ""
	}
	return url.(string)
}
func (content *FileContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *FileContent) GetData() []byte {
	if content._data == nil {
		b64 := content.Get("data")
		if b64 != nil {
			content._data = Base64Decode(b64.(string))
		}
	}
	return content._data
}
func (content *FileContent) SetData(data []byte) {
	if data == nil {
		content.Set("data", nil)
	} else {
		b64 := Base64Encode(data)
		content.Set("data", b64)
	}
	content._data = data
}

func (content *FileContent) GetFilename() string {
	filename := content.Get("filename")
	if filename == nil {
		return ""
	}
	return filename.(string)
}
func (content *FileContent) SetFilename(filename string) {
	content.Set("filename", filename)
}

func (content *FileContent) GetPassword() SymmetricKey {
	if content._key == nil {
		dict := content.Get("password")
		content._key = SymmetricKeyParse(dict)
	}
	return content._key
}

func (content *FileContent) SetPassword(password SymmetricKey) {
	if password == nil {
		content.Set("password", nil)
	} else {
		content.Set("password", password.GetMap(false))
	}
	content._key = password
}

/**
 *  Image message: {
 *      type : 0x12,
 *      sn   : 123,
 *
 *      URL       : "http://", // upload to CDN
 *      data      : "...",     // if (!URL) base64_encode(image)
 *      thumbnail : "...",     // base64_encode(smallImage)
 *      filename  : "..."
 *  }
 */
type ImageContent struct {
	FileContent

	_thumbnail []byte
}

func (content *ImageContent) Init(dict map[string]interface{}) *ImageContent {
	if content.FileContent.Init(dict) != nil {
		// lazy load
		content._thumbnail = nil
	}
	return content
}

func (content *ImageContent) InitWithFilename(filename string, data []byte) *ImageContent {
	if content.FileContent.InitWithType(IMAGE, filename, data) != nil {
		content._thumbnail = nil
	}
	return content
}

func (content *ImageContent) GetThumbnail() []byte {
	if content._thumbnail == nil {
		b64 := content.Get("thumbnail")
		if b64 != nil {
			content._thumbnail = Base64Decode(b64.(string))
		}
	}
	return content._thumbnail
}

func (content *ImageContent) SetThumbnail(thumbnail []byte) {
	if thumbnail == nil {
		content.Set("thumbnail", nil)
	} else {
		b64 := Base64Encode(thumbnail)
		content.Set("thumbnail", b64)
	}
	content._thumbnail = thumbnail
}

/**
 *  Audio message: {
 *      type : 0x14,
 *      sn   : 123,
 *
 *      URL      : "http://", // upload to CDN
 *      data     : "...",     // if (!URL) base64_encode(audio)
 *      text     : "...",     // Automatic Speech Recognition
 *      filename : "..."
 *  }
 */
type AudioContent struct {
	FileContent
}

func (content *AudioContent) Init(dict map[string]interface{}) *AudioContent {
	if content.FileContent.Init(dict) != nil {
	}
	return content
}

func (content *AudioContent) InitWithFilename(filename string, data []byte) *AudioContent {
	if content.FileContent.InitWithType(AUDIO, filename, data) != nil {
		// init
	}
	return content
}

/**
 *  Video message: {
 *      type : 0x16,
 *      sn   : 123,
 *
 *      URL      : "http://", // upload to CDN
 *      data     : "...",     // if (!URL) base64_encode(video)
 *      snapshot : "...",     // base64_encode(smallImage)
 *      filename : "..."
 *  }
 */
type VideoContent struct {
	FileContent

	_snapshot []byte
}

func (content *VideoContent) Init(dict map[string]interface{}) *VideoContent {
	if content.FileContent.Init(dict) != nil {
		// lazy load
		content._snapshot = nil
	}
	return content
}

func (content *VideoContent) InitWithFilename(filename string, data []byte) *VideoContent {
	if content.FileContent.InitWithType(VIDEO, filename, data) != nil {
		content._snapshot = nil
	}
	return content
}

func (content *VideoContent) GetSnapshot() []byte {
	if content._snapshot == nil {
		b64 := content.Get("snapshot")
		if b64 != nil {
			content._snapshot = Base64Decode(b64.(string))
		}
	}
	return content._snapshot
}

func (content *VideoContent) SetSnapshot(snapshot []byte) {
	if snapshot == nil {
		content.Set("snapshot", nil)
	} else {
		b64 := Base64Encode(snapshot)
		content.Set("snapshot", b64)
	}
	content._snapshot = snapshot
}
