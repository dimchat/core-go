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
	content := new(BaseFileContent).InitWithType(msgType, filename, data)
	ObjectRetain(content)
	return content
}

/* designated initializer */
func (content *BaseFileContent) Init(dict map[string]interface{}) *BaseFileContent {
	if content.BaseContent.Init(dict) != nil {
		// lazy load
		content._data = nil
		content.setPassword(nil)
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
		content.setPassword(nil)
	}
	return content
}

func (content *BaseFileContent) Release() int {
	cnt := content.BaseContent.Release()
	if cnt == 0 {
		// this object is going to be destroyed,
		// release children
		content.setPassword(nil)
	}
	return cnt
}

func (content *BaseFileContent) setPassword(key SymmetricKey) {
	if key != content._key {
		ObjectRetain(key)
		ObjectRelease(content._key)
		content._key = key
	}
}

//-------- IFileContent

func (content *BaseFileContent) URL() string {
	url := content.Get("URL")
	if url == nil {
		return ""
	}
	return url.(string)
}
func (content *BaseFileContent) SetURL(url string) {
	content.Set("URL", url)
}

func (content *BaseFileContent) Data() []byte {
	if content._data == nil {
		b64 := content.Get("data")
		if b64 != nil {
			content._data = Base64Decode(b64.(string))
		}
	}
	return content._data
}
func (content *BaseFileContent) SetData(data []byte) {
	if data == nil {
		content.Set("data", nil)
	} else {
		b64 := Base64Encode(data)
		content.Set("data", b64)
	}
	content._data = data
}

func (content *BaseFileContent) Filename() string {
	filename := content.Get("filename")
	if filename == nil {
		return ""
	}
	return filename.(string)
}
func (content *BaseFileContent) SetFilename(filename string) {
	content.Set("filename", filename)
}

func (content *BaseFileContent) Password() SymmetricKey {
	if content._key == nil {
		dict := content.Get("password")
		content.setPassword(SymmetricKeyParse(dict))
	}
	return content._key
}

func (content *BaseFileContent) SetPassword(password SymmetricKey) {
	if password == nil {
		content.Set("password", nil)
	} else {
		content.Set("password", password.GetMap(false))
	}
	content.setPassword(password)
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
	content := new(ImageFileContent).InitWithFilename(filename, data)
	ObjectRetain(content)
	return content
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
		b64 := content.Get("thumbnail")
		if b64 != nil {
			content._thumbnail = Base64Decode(b64.(string))
		}
	}
	return content._thumbnail
}

func (content *ImageFileContent) SetThumbnail(thumbnail []byte) {
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
type AudioFileContent struct {
	BaseFileContent
	IAudioContent
}

func NewAudioContent(filename string, data []byte) AudioContent {
	content := new(AudioFileContent).InitWithFilename(filename, data)
	ObjectRetain(content)
	return content
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
	duration := content.Get("duration")
	if duration == nil {
		return 0
	} else {
		return duration.(int)
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
	content := new(VideoFileContent).InitWithFilename(filename, data)
	ObjectRetain(content)
	return content
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
		b64 := content.Get("snapshot")
		if b64 != nil {
			content._snapshot = Base64Decode(b64.(string))
		}
	}
	return content._snapshot
}

func (content *VideoFileContent) SetSnapshot(snapshot []byte) {
	if snapshot == nil {
		content.Set("snapshot", nil)
	} else {
		b64 := Base64Encode(snapshot)
		content.Set("snapshot", b64)
	}
	content._snapshot = snapshot
}
