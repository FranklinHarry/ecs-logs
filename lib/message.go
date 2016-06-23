package ecslogs

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type MessageReader interface {
	ReadMessage() (Message, error)
}

type MessageReadCloser interface {
	io.Closer

	MessageReader
}

type MessageWriter interface {
	WriteMessage(Message) error
}

type MessageWriteCloser interface {
	io.Closer

	MessageWriter
}

type MessageBatchWriter interface {
	MessageWriter

	WriteMessageBatch([]Message) error
}

type MessageBatchWriteCloser interface {
	io.Closer

	MessageBatchWriter
}

type Message struct {
	Level   Level     `json:"level"`
	PID     int       `json:"pid,omitempty"`
	UID     int       `json:"uid,omitempty"`
	GID     int       `json:"gid,omitempty"`
	Errno   int       `json:"errno,omitempty"`
	Line    int       `json:"line,omitempty"`
	Func    string    `json:"func,omitempty"`
	File    string    `json:"file,omitempty"`
	ID      string    `json:"id,omitempty"`
	Host    string    `json:"host,omitempty"`
	Group   string    `json:"group,omitempty"`
	Stream  string    `json:"stream,omitempty"`
	Content string    `json:"content,omitempty"`
	Time    time.Time `json:"time,omitempty"`
}

func (m Message) String() string {
	return fmt.Sprint(m)
}

func (m Message) Format(f fmt.State, _ rune) {
	b, _ := json.Marshal(m)
	f.Write(b)
}

func NewMessageEncoder(w io.Writer) MessageBatchWriteCloser {
	return messageEncoder{
		j: json.NewEncoder(w),
		w: w,
	}
}

type messageEncoder struct {
	j *json.Encoder
	w io.Writer
}

func (e messageEncoder) Close() (err error) {
	if c, ok := e.w.(io.Closer); ok {
		err = c.Close()
	}
	return
}

func (e messageEncoder) WriteMessage(msg Message) (err error) {
	err = e.j.Encode(msg)
	return
}

func (e messageEncoder) WriteMessageBatch(batch []Message) (err error) {
	for _, msg := range batch {
		if err = e.WriteMessage(msg); err != nil {
			return
		}
	}
	return
}

func NewMessageDecoder(r io.Reader) MessageReadCloser {
	return messageDecoder{
		j: json.NewDecoder(r),
		r: r,
	}
}

type messageDecoder struct {
	j *json.Decoder
	r io.Reader
}

func (d messageDecoder) Close() (err error) {
	if c, ok := d.r.(io.Closer); ok {
		err = c.Close()
	}
	return
}

func (d messageDecoder) ReadMessage() (msg Message, err error) {
	err = d.j.Decode(&msg)
	return
}
