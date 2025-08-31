package common

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"
)

type Socket struct {
	conn net.Conn
}

func NewSocket() *Socket {
	return &Socket{}
}

func (s *Socket) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *Socket) writeFull(data []byte) error {
	total := 0
	for total < len(data) {
		n, err := s.conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}


func (s *Socket) Send(msg string) error {
	msgBytes := []byte(msg)
	length := uint32(len(msgBytes))

	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)

	if err := s.writeFull(lengthBytes); err != nil {
		return err
	}
	if err := s.writeFull(msgBytes); err != nil {
		return err
	}
	return nil
}

func (s *Socket) ReadResponse(ctx context.Context) (string, error) {
	readCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		msg, err := bufio.NewReader(s.conn).ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		readCh <- msg
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		return "", err
	case msg := <-readCh:
		return msg, nil
	}
}


func (s *Socket) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

