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

//goland:noinspection GoSnakeCaseUsage
const (

	//
	//  type
	//
	MIME_TEXT  = "text"
	MIME_IMAGE = "image"
	MIME_AUDIO = "audio"
	MIME_VIDEO = "video"
	MIME_APP   = "application"

	//
	//  subtype
	//
	MIME_PLAIN = "plain"
	MIME_HTML  = "html"
	MIME_XML   = "xml"
	MIME_CSS   = "css"
	MIME_JS    = "javascript"

	MIME_BMP   = "bmp"
	MIME_GIF   = "gif"
	MIME_PNG   = "png"
	MIME_JPG   = "jpeg"
	MIME_ICON  = "x-icon"
	MIME_SVG   = "svg+xml"
	MIME_WEB_P = "webp"

	MIME_WAV   = "wav"
	MIME_OGG   = "ogg"
	MIME_MP3   = "mp3"
	MIME_MP4   = "mp4"
	MIME_MPG   = "mpeg"
	MIME_WEB_M = "webm"

	MIME_PDF          = "pdf"
	MIME_WORD         = "msword"
	MIME_EXCEL        = "vnd.ms-excel"
	MIME_PPT          = "vnd.ms-powerpoint"
	MIME_ZIP          = "zip"
	MIME_JSON         = "json"
	MIME_OCTET_STREAM = "octet-stream"
)

//goland:noinspection GoSnakeCaseUsage
var MIMEType = struct {

	//
	//  text/*
	//
	TEXT_PLAIN string
	TEXT_HTML  string
	TEXT_XML   string
	TEXT_CSS   string
	TEXT_JS    string

	//
	//  image/*
	//
	IMAGE_BMP   string
	IMAGE_GIF   string
	IMAGE_PNG   string
	IMAGE_JPG   string
	IMAGE_ICON  string
	IMAGE_SVG   string
	IMAGE_WEB_P string

	//
	//  audio/*
	//
	AUDIO_WAV string
	AUDIO_OGG string
	AUDIO_MP3 string
	AUDIO_MP4 string
	AUDIO_MPG string

	//
	//  video/*
	//
	VIDEO_MP4   string
	VIDEO_MPG   string
	VIDEO_OGG   string
	VIDEO_WEB_M string

	//
	//  application/*
	//
	APP_PDF          string
	APP_WORD         string
	APP_EXCEL        string
	APP_PPT          string
	APP_ZIP          string
	APP_XML          string
	APP_JSON         string
	APP_OCTET_STREAM string
}{
	//
	//  text/*
	//
	TEXT_PLAIN: MIME_TEXT + MIME_PLAIN, //  text/plain
	TEXT_HTML:  MIME_TEXT + MIME_HTML,  //  text/html
	TEXT_XML:   MIME_TEXT + MIME_XML,   //  text/xml
	TEXT_CSS:   MIME_TEXT + MIME_CSS,   //  text/css
	TEXT_JS:    MIME_TEXT + MIME_JS,    //  text/javascript

	//
	//  image/*
	//
	IMAGE_BMP:   MIME_IMAGE + MIME_BMP,   //  image/bmp
	IMAGE_GIF:   MIME_IMAGE + MIME_GIF,   //  image/gif
	IMAGE_PNG:   MIME_IMAGE + MIME_PNG,   //  image/png
	IMAGE_JPG:   MIME_IMAGE + MIME_JPG,   //  image/jpeg
	IMAGE_ICON:  MIME_IMAGE + MIME_ICON,  //  image/x-icon
	IMAGE_SVG:   MIME_IMAGE + MIME_SVG,   //  image/svg+xml
	IMAGE_WEB_P: MIME_IMAGE + MIME_WEB_P, //  image/webp

	//
	//  audio/*
	//
	AUDIO_WAV: MIME_AUDIO + MIME_WAV, //  audio/wav
	AUDIO_OGG: MIME_AUDIO + MIME_OGG, //  audio/ogg
	AUDIO_MP3: MIME_AUDIO + MIME_MP3, //  audio/mp3
	AUDIO_MP4: MIME_AUDIO + MIME_MP4, //  audio/mp4
	AUDIO_MPG: MIME_AUDIO + MIME_MPG, //  audio/mpeg

	//
	//  video/*
	//
	VIDEO_MP4:   MIME_VIDEO + MIME_MP4,   //  video/mp4
	VIDEO_MPG:   MIME_VIDEO + MIME_MPG,   //  video/mpeg
	VIDEO_OGG:   MIME_VIDEO + MIME_OGG,   //  video/ogg
	VIDEO_WEB_M: MIME_VIDEO + MIME_WEB_M, //  video/webm

	//
	//  application/*
	//
	APP_PDF:          MIME_APP + MIME_PDF,          //  application/pdf
	APP_WORD:         MIME_APP + MIME_WORD,         //  application/msword
	APP_EXCEL:        MIME_APP + MIME_EXCEL,        //  application/vnd.ms-excel
	APP_PPT:          MIME_APP + MIME_PPT,          //  application/vnd.ms-powerpoint
	APP_ZIP:          MIME_APP + MIME_ZIP,          //  application/zip
	APP_XML:          MIME_APP + MIME_XML,          //  application/xml
	APP_JSON:         MIME_APP + MIME_JSON,         //  application/json
	APP_OCTET_STREAM: MIME_APP + MIME_OCTET_STREAM, //  application/octet-stream
}
