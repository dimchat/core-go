# Decentralized Instant Messaging Protocol (Go)

[![License](https://img.shields.io/github/license/dimchat/core-go)](https://github.com/dimchat/core-go/blob/main/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/dimchat/core-go/pulls)
[![Platform](https://img.shields.io/github/go-mod/go-version/dimchat/core-go)](https://github.com/dimchat/core-go/wiki)
[![Issues](https://img.shields.io/github/issues/dimchat/core-go)](https://github.com/dimchat/core-go/issues)
[![Repo Size](https://img.shields.io/github/repo-size/dimchat/core-go)](https://github.com/dimchat/core-go/archive/refs/heads/main.zip)
[![Tags](https://img.shields.io/github/tag/dimchat/core-go)](https://github.com/dimchat/core-go/tags)

[![Watchers](https://img.shields.io/github/watchers/dimchat/core-go)](https://github.com/dimchat/core-go/watchers)
[![Forks](https://img.shields.io/github/forks/dimchat/core-go)](https://github.com/dimchat/core-go/forks)
[![Stars](https://img.shields.io/github/stars/dimchat/core-go)](https://github.com/dimchat/core-go/stargazers)
[![Followers](https://img.shields.io/github/followers/dimchat)](https://github.com/orgs/dimchat/followers)

## Dependencies

* Latest Versions

| Name | Version | Description |
|------|---------|-------------|
| [Ming Ke Ming (名可名)](https://github.com/dimchat/mkm-go) | [![Tags](https://img.shields.io/github/tag/dimchat/mkm-go)](https://github.com/dimchat/mkm-go/tags) | Decentralized User Identity Authentication |
| [Dao Ke Dao (道可道)](https://github.com/dimchat/dkd-go) | [![Tags](https://img.shields.io/github/tag/dimchat/dkd-go)](https://github.com/dimchat/dkd-go/tags) | Universal Message Module |

## Examples

### Extends Command

* _Handshake Command Protocol_
  0. (C-S) handshake start
  1. (S-C) handshake again with new session
  2. (C-S) handshake restart with new session
  3. (S-C) handshake success

```go
type HandshakeState uint8

const (
	HandshakeInit    HandshakeState = iota
	HandshakeStart   // C -> S, without session key(or session expired)
	HandshakeAgain   // S -> C, with new session key
	HandshakeRestart // C -> S, with new session key
	HandshakeSuccess // S -> C, handshake accepted
)

func (state HandshakeState) String() string {
	switch state {
	case HandshakeInit:
		return "HandshakeInit"
	case HandshakeStart:
		return "HandshakeStart"
	case HandshakeAgain:
		return "HandshakeAgain"
	case HandshakeRestart:
		return "HandshakeRestart"
	case HandshakeSuccess:
		return "HandshakeSuccess"
	default:
		return fmt.Sprintf("HandshakeState(%d)", state)
	}
}

const HANDSHAKE = "handshake"

// HandshakeCommand defines the interface for handshake commands (session initialization)
//
// # Implements the Command interface for DIM network session establishment
//
//	Data Format: {
//	    "type": 0x88,
//	    "sn": 123,
//
//	    "command": "handshake",
//	    "title": "Hello world!",   // Handshake state indicator ("DIM?", "DIM!")
//	    "session": "{SESSION_KEY}" // Session key for authenticated communication
//	}
type HandshakeCommand interface {
	Command

	// Title returns the handshake state indicator (e.g., "DIM?", "DIM!")
	//
	// Returns: String representing the current handshake state
	Title() string

	// SessionKey returns the session key for authenticated communication
	//
	// Returns: Session key string (empty string if not established)
	SessionKey() string

	// State returns the structured HandshakeState derived from the title
	//
	// Returns: Enumerated HandshakeState value
	State() HandshakeState
}
```

```go
func getState(title string, session string) HandshakeState {
	// check message text
	if title == "" {
		return HandshakeInit
	}
	if title == "DIM!" /*|| message == "OK!"*/ {
		return HandshakeSuccess
	}
	if title == "DIM?" {
		return HandshakeAgain
	}
	// check session key
	if session == "" {
		return HandshakeStart
	}
	return HandshakeRestart
}

type BaseHandshakeCommand struct {
	//HandshakeCommand
	*BaseCommand
}

func NewBaseHandshakeCommand(dict StringKeyMap, title, sessionKey string) *BaseHandshakeCommand {
	if dict != nil {
		// init handshake command with map
		return &BaseHandshakeCommand{
			BaseCommand: NewBaseCommand(dict, "", ""),
		}
	}
	// new handshake command
	content := &BaseHandshakeCommand{
		BaseCommand: NewBaseCommand(nil, "", HANDSHAKE),
	}
	// text message
	content.Set("title", title)
	// session key
	content.Set("session", sessionKey)
	// OK
	return content
}

// Override
func (content *BaseHandshakeCommand) Title() string {
	return content.GetString("title", "")
}

// Override
func (content *BaseHandshakeCommand) SessionKey() string {
	return content.GetString("session", "")
}

// Override
func (content *BaseHandshakeCommand) State() HandshakeState {
	title := content.Title()
	session := content.SessionKey()
	return getState(title, session)
}
```

### Extends Content

```go
import (
	. "github.com/dimchat/core-go/dkd"
	. "github.com/dimchat/core-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

//	Data structure: {
//	    "type" : i2s(0xA0),
//	    "sn"   : 123,
//
//	    "app"  : "{APP_ID}",  // application (e.g.: "chat.dim.sechat")
//	    "extra": info         // others
//	}

type ApplicationContent struct {
	//AppContent
	*BaseContent
}

func NewApplicationContent(dict StringKeyMap, app string) AppContent {
	if dict != nil {
		// init application content with map
		return &ApplicationContent{
			BaseContent: NewBaseContent(dict, ""),
		}
	}
	// new application content
	content := &ApplicationContent{
		BaseContent: NewBaseContent(nil, ContentType.APPLICATION),
	}
	content.Set("app", app)
	return content
}

// Override
func (content *ApplicationContent) Application() string {
	return content.GetString("app", "")
}
```

### Extends ID Address

* Examples in [plugins-go](https://github.com/dimchat/plugins-go)

----

Copyright &copy; 2020-2026 Albert Moky
[![Followers](https://img.shields.io/github/followers/moky)](https://github.com/moky?tab=followers)
