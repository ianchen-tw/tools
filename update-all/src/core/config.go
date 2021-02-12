package core

import (
	"github.com/kyokomi/emoji/v2"
)

// IsDebug is in debug mode
var IsDebug bool

// UseWorkdirToFetch read config/ write cache to work dir
var UseWorkdirToFetch bool

var runRecordFilename = "runRecord.json"
var routineFileName = "config.yaml"

var emoji_skip = emoji.Sprint(":fast-forward_button:")
var emoji_execute = emoji.Sprint(":rocket:")
