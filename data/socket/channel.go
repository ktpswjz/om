package socket

import (
	"container/list"
	"sync"
)

type Channel interface {
	Write(message *Message)
	Read() <-chan *Message

	getElement() *list.Element
	close()
}

type innerChannel struct {
	channel chan *Message
	element *list.Element
}

func (s *innerChannel) Write(message *Message) {
	select {
	case s.channel <- message:
	default:
	}
}

func (s *innerChannel) Read() <-chan *Message {
	return s.channel
}

func (s *innerChannel) getElement() *list.Element {
	return s.element
}

func (s *innerChannel) close() {
	close(s.channel)
}

type ChannelCollection interface {
	NewChannel() Channel
	Remove(channel Channel)
	Write(message *Message)
}

func NewChannelCollection() ChannelCollection {
	instance := &innerChannelCollection{}
	instance.channels = list.New()

	return instance
}

type innerChannelCollection struct {
	sync.RWMutex

	channels *list.List
}

func (s *innerChannelCollection) NewChannel() Channel {
	s.Lock()
	defer s.Unlock()

	instance := &innerChannel{}
	instance.channel = make(chan *Message, 1024)
	instance.element = s.channels.PushBack(instance)

	return instance
}

func (s *innerChannelCollection) Remove(channel Channel) {
	if channel == nil {
		return
	}

	s.Lock()
	defer s.Unlock()

	s.channels.Remove(channel.getElement())
	channel.close()
}

func (s *innerChannelCollection) Write(message *Message) {
	s.Lock()
	defer s.Unlock()

	for e := s.channels.Front(); e != nil; {
		ev, ok := e.Value.(Channel)
		if !ok {
			return
		}

		ev.Write(message)
		e = e.Next()
	}
}
