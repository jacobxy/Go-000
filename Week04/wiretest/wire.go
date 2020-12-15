//+build wireinject

package main

import (
	"github.com/google/wire"
)

func InitializeSpeaker(text string) NewSpeaker {
	wire.Build(NewSpeaker, NewMessage)
	return Speaker{}
}
