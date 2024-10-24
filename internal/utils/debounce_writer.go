package utils

import (
	"context"
	"sync"
	"time"
)

type DebouncedWriter struct {
	ctx             context.Context
	interval        time.Duration
	maxWaitInterval time.Duration
	onWrite         func(bytes []byte) error
	onError         func(error)
	buffer          [][]byte
	stop            chan struct{}
	mutex           sync.Mutex
	forceWrite      <-chan time.Time
}

func NewDebouncedWriter(
	ctx context.Context,
	interval time.Duration,
	maxWaitInterval time.Duration,
	onWrite func(bytes []byte) error,
	onError func(error),
) *DebouncedWriter {
	return &DebouncedWriter{
		ctx:             ctx,
		interval:        interval,
		maxWaitInterval: maxWaitInterval,
		onWrite:         onWrite,
		onError:         onError,
		buffer:          make([][]byte, 0),
	}
}

func (b *DebouncedWriter) Flush() {
	b.doWrite()
}

func (b *DebouncedWriter) Write(bytes []byte) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.stop != nil {
		close(b.stop)
	}
	b.stop = make(chan struct{})
	b.buffer = append(b.buffer, bytes)

	if b.forceWrite == nil {
		b.forceWrite = time.After(b.maxWaitInterval)
	}

	go b.writeSelector(b.stop, b.forceWrite)
	return nil
}

func (b *DebouncedWriter) writeSelector(stop chan struct{}, forceWrite <-chan time.Time) {
	select {
	case <-stop:
	case <-b.ctx.Done():
	case <-forceWrite:
		b.doWrite()
	case <-time.After(b.interval):
		b.doWrite()
	}
}

func (b *DebouncedWriter) doWrite() {
	b.mutex.Lock()
	flatBuffer := Flat(b.buffer)
	clear(b.buffer)
	b.forceWrite = nil
	b.mutex.Unlock()

	err := b.onWrite(flatBuffer)
	if err != nil {
		b.onError(err)
	}
}
