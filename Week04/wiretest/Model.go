package main

import "fmt"

type Message string

func NewMessage(text string) Message {
	return Message(text)
}

type Speaker struct {
	Message Message
}

func (s Speaker) Say() {
	fmt.Println(s.Message)
}

func NewSpeaker(m Message) Speaker {
	return Speaker{
		Message: m,
	}
}

type MyTest struct {
	Message Message
	Spk     Speaker
}

func NewMyTest(m Message, s Speaker) MyTest {
	return MyTest{
		Message: m,
		Spk:     s,
	}
}
