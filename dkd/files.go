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
	. "github.com/dimchat/mkm-go/crypto"
	. "github.com/dimchat/mkm-go/format"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

/**
 *  Image File Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x12),
 *      "sn"   : 123,
 *
 *      "data"     : "...",        // base64_encode(fileContent)
 *      "filename" : "photo.png",
 *
 *      "URL"      : "http://...", // download from CDN
 *      // before fileContent uploaded to a public CDN,
 *      // it should be encrypted by a symmetric key
 *      "key"      : {             // symmetric key to decrypt file data
 *          "algorithm" : "AES",   // "DES", ...
 *          "data"      : "{BASE64_ENCODE}",
 *          ...
 *      },
 *      "thumbnail" : "data:image/jpeg;base64,..."
 *  }
 *  </pre></blockquote>
 */
type ImageFileContent struct {
	//ImageContent
	*BaseFileContent

	thumbnail TransportableFile
}

func NewImageFileContent(dict StringKeyMap,
	data TransportableData, filename string, url URL, password DecryptKey,
	thumbnail TransportableFile,
) *ImageFileContent {
	return &ImageFileContent{
		BaseFileContent: NewBaseFileContent(dict, ContentType.IMAGE, data, filename, url, password),
		thumbnail:       thumbnail,
	}
}

// Override
func (content *ImageFileContent) Map() StringKeyMap {
	// serialize 'thumbnail'
	img := content.thumbnail
	if img != nil && !content.Contains("thumbnail") {
		content.Set("thumbnail", img.Serialize())
	}
	// OK
	return content.BaseFileContent.Map()
}

// Override
func (content *ImageFileContent) Thumbnail() TransportableFile {
	img := content.thumbnail
	if img == nil {
		uri := content.Get("thumbnail")
		img = ParseTransportableFile(uri)
		content.thumbnail = img
	}
	return img
}

// Override
func (content *ImageFileContent) SetThumbnail(thumbnail TransportableFile) {
	content.Remove("thumbnail")
	//content.SetMapper("thumbnail", thumbnail)
	content.thumbnail = thumbnail
}

/**
 *  Audio File Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x14),
 *      "sn"   : 123,
 *
 *      "data"     : "...",        // base64_encode(fileContent)
 *      "filename" : "voice.mp4",
 *
 *      "URL"      : "http://...", // download from CDN
 *      // before fileContent uploaded to a public CDN,
 *      // it should be encrypted by a symmetric key
 *      "key"      : {             // symmetric key to decrypt file data
 *          "algorithm" : "AES",   // "DES", ...
 *          "data"      : "{BASE64_ENCODE}",
 *          ...
 *      },
 *      "duration" : 123.45,
 *      "text"     : "..."         // Automatic Speech Recognition
 *  }
 *  </pre></blockquote>
 */
type AudioFileContent struct {
	//AudioContent
	*BaseFileContent
}

func NewAudioFileContent(dict StringKeyMap,
	data TransportableData, filename string, url URL, password DecryptKey,
) *AudioFileContent {
	return &AudioFileContent{
		BaseFileContent: NewBaseFileContent(dict, ContentType.AUDIO, data, filename, url, password),
	}
}

// Override
func (content *AudioFileContent) Duration() float64 {
	return content.GetFloat64("duration", 0)
}

// Override
func (content *AudioFileContent) SetDuration(duration float64) {
	content.Set("duration", duration)
}

// Override
func (content *AudioFileContent) Text() string {
	return content.GetString("text", "")
}

// Override
func (content *AudioFileContent) SetText(text string) {
	if text == "" {
		content.Remove("text")
	} else {
		content.Set("text", text)
	}
}

/**
 *  Video File Content
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type" : i2s(0x16),
 *      "sn"   : 123,
 *
 *      "data"     : "...",        // base64_encode(fileContent)
 *      "filename" : "movie.mp4",
 *
 *      "URL"      : "http://...", // download from CDN
 *      // before fileContent uploaded to a public CDN,
 *      // it should be encrypted by a symmetric key
 *      "key"      : {             // symmetric key to decrypt file data
 *          "algorithm" : "AES",   // "DES", ...
 *          "data"      : "{BASE64_ENCODE}",
 *          ...
 *      },
 *      "snapshot" : "data:image/jpeg;base64,..."
 *  }
 *  </pre></blockquote>
 */
type VideoFileContent struct {
	//VideoContent
	*BaseFileContent

	snapshot TransportableFile
}

func NewVideoFileContent(dict StringKeyMap,
	data TransportableData, filename string, url URL, password DecryptKey,
	snapshot TransportableFile,
) *VideoFileContent {
	return &VideoFileContent{
		BaseFileContent: NewBaseFileContent(dict, ContentType.VIDEO, data, filename, url, password),
		snapshot:        snapshot,
	}
}

// Override
func (content *VideoFileContent) Map() StringKeyMap {
	// serialize 'snapshot'
	img := content.snapshot
	if img != nil && !content.Contains("snapshot") {
		content.Set("snapshot", img.Serialize())
	}
	// OK
	return content.BaseFileContent.Map()
}

// Override
func (content *VideoFileContent) Snapshot() TransportableFile {
	img := content.snapshot
	if img == nil {
		uri := content.Get("snapshot")
		img = ParseTransportableFile(uri)
		content.snapshot = img
	}
	return img
}

// Override
func (content *VideoFileContent) SetSnapshot(snapshot TransportableFile) {
	content.Remove("snapshot")
	//content.SetMapper("snapshot", snapshot)
	content.snapshot = snapshot
}
