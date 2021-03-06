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
package dkd

import (
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/dkd-go/dkd"
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
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
type BaseFileContent struct {
	BaseContent
	IFileContent

	_data []byte       // file data (plaintext)
	_key SymmetricKey  // symmetric key to decrypt the encrypted data from URL
}

func NewFileContent(msgType uint8, filename string, data []byte) FileContent {
	return new(BaseFileContent).InitWithType(msgType, filename, data)
}

/* designated initializer */
func (content *BaseFileContent) Init(dict map[string]interface{}) *BaseFileContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._data = nil
		content._key = nil
	}
	return content
}

/* designated initializer */
func (content *BaseFileContent) InitWithType(msgType uint8, filename string, data []byte) *BaseFileContent {
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

//-------- IFileContent

func (content *BaseFileContent) URL() string {
	url, ok := content.Get("URL").(string)
	if ok {
		return url
	} else {
		return ""
	}
}
func (content *BaseFileContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *BaseFileContent) Data() []byte {
	if content._data == nil {
		b64, ok := content.Get("data").(string)
		if ok {
			content._data = Base64Decode(b64)
		}
	}
	return content._data
}
func (content *BaseFileContent) SetData(data []byte) {
	if ValueIsNil(data) {
		content.Remove("data")
	} else {
		b64 := Base64Encode(data)
		content.Set("data", b64)
	}
	content._data = data
}

func (content *BaseFileContent) Filename() string {
	text, ok := content.Get("filename").(string)
	if ok {
		return text
	} else {
		return ""
	}
}
func (content *BaseFileContent) SetFilename(filename string) {
	content.Set("filename", filename)
}

func (content *BaseFileContent) Password() SymmetricKey {
	if content._key == nil {
		dict := content.Get("password")
		content._key = SymmetricKeyParse(dict)
	}
	return content._key
}

func (content *BaseFileContent) SetPassword(password SymmetricKey) {
	if ValueIsNil(password) {
		content.Remove("password")
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

type ImageFileContent struct {
	BaseFileContent
	IImageContent

	_thumbnail []byte
}

func NewImageContent(filename string, data []byte) ImageContent {
	return new(ImageFileContent).InitWithFilename(filename, data)
}

func (content *ImageFileContent) Init(dict map[string]interface{}) *ImageFileContent {
	if content.BaseFileContent.Init(dict) != nil {
		// lazy load
		content._thumbnail = nil
	}
	return content
}

func (content *ImageFileContent) InitWithFilename(filename string, data []byte) *ImageFileContent {
	if content.BaseFileContent.InitWithType(IMAGE, filename, data) != nil {
		content._thumbnail = nil
	}
	return content
}

//-------- IImageContent

func (content *ImageFileContent) Thumbnail() []byte {
	if content._thumbnail == nil {
		b64, ok := content.Get("thumbnail").(string)
		if ok {
			content._thumbnail = Base64Decode(b64)
		}
	}
	return content._thumbnail
}

func (content *ImageFileContent) SetThumbnail(thumbnail []byte) {
	if ValueIsNil(thumbnail) {
		content.Remove("thumbnail")
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
type AudioFileContent struct {
	BaseFileContent
	IAudioContent
}

func NewAudioContent(filename string, data []byte) AudioContent {
	return new(AudioFileContent).InitWithFilename(filename, data)
}

func (content *AudioFileContent) Init(dict map[string]interface{}) *AudioFileContent {
	if content.BaseFileContent.Init(dict) != nil {
	}
	return content
}

func (content *AudioFileContent) InitWithFilename(filename string, data []byte) *AudioFileContent {
	if content.BaseFileContent.InitWithType(AUDIO, filename, data) != nil {
	}
	return content
}

//-------- IAudioContent

func (content *AudioFileContent) Duration() int {
	duration, ok := content.Get("duration").(int)
	if ok {
		return duration
	} else {
		return 0
	}
}
func (content *AudioFileContent) SetDuration(duration int) {
	content.Set("duration", duration)
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
type VideoFileContent struct {
	BaseFileContent
	IVideoContent

	_snapshot []byte
}

func NewVideoContent(filename string, data []byte) VideoContent {
	return new(VideoFileContent).InitWithFilename(filename, data)
}

func (content *VideoFileContent) Init(dict map[string]interface{}) *VideoFileContent {
	if content.BaseFileContent.Init(dict) != nil {
		// lazy load
		content._snapshot = nil
	}
	return content
}

func (content *VideoFileContent) InitWithFilename(filename string, data []byte) *VideoFileContent {
	if content.BaseFileContent.InitWithType(VIDEO, filename, data) != nil {
		content._snapshot = nil
	}
	return content
}

//-------- IVideoContent

func (content *VideoFileContent) Snapshot() []byte {
	if content._snapshot == nil {
		b64, ok := content.Get("snapshot").(string)
		if ok {
			content._snapshot = Base64Decode(b64)
		}
	}
	return content._snapshot
}

func (content *VideoFileContent) SetSnapshot(snapshot []byte) {
	if ValueIsNil(snapshot) {
		content.Remove("snapshot")
	} else {
		b64 := Base64Encode(snapshot)
		content.Set("snapshot", b64)
	}
	content._snapshot = snapshot
}
