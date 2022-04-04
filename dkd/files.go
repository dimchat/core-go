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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/types"
)

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

	_thumbnail []byte
}

func NewImageContent(filename string, data []byte) ImageContent {
	content := new(ImageFileContent)
	content.InitWithFilename(filename, data)
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
}

func NewAudioContent(filename string, data []byte) AudioContent {
	content := new(AudioFileContent)
	content.InitWithFilename(filename, data)
	return content
}

//func (content *AudioFileContent) Init(dict map[string]interface{}) *AudioFileContent {
//	if content.BaseFileContent.Init(dict) != nil {
//	}
//	return content
//}

func (content *AudioFileContent) InitWithFilename(filename string, data []byte) *AudioFileContent {
	if content.BaseFileContent.InitWithType(AUDIO, filename, data) != nil {
	}
	return content
}

//-------- IAudioContent

func (content *AudioFileContent) Duration() float64 {
	duration := content.Get("duration")
	if duration == nil {
		return 0.0
	}
	return duration.(float64)
}
func (content *AudioFileContent) SetDuration(duration float64) {
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

	_snapshot []byte
}

func NewVideoContent(filename string, data []byte) VideoContent {
	content := new(VideoFileContent)
	content.InitWithFilename(filename, data)
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
	if ValueIsNil(snapshot) {
		content.Remove("snapshot")
	} else {
		b64 := Base64Encode(snapshot)
		content.Set("snapshot", b64)
	}
	content._snapshot = snapshot
}
