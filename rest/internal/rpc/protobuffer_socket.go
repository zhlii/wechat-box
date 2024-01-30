package rpc

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/zhlii/wechat-box/rest/internal/logs"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/pair1"
	"go.nanomsg.org/mangos/v3/transport/all"
	_ "go.nanomsg.org/mangos/v3/transport/all"
	"google.golang.org/protobuf/proto"
)

type protobufferSocket struct {
	addr   string
	socket protocol.Socket
}

func newProtobufferSocker(ip string, port int) *protobufferSocket {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	return &protobufferSocket{addr: "tcp://" + addr}
}

// connect to rpc server
// timeout second
func (s *protobufferSocket) conn(timeout uint) (err error) {
	all.AddTransports(nil)
	if s.socket, err = pair1.NewSocket(); err != nil {
		return err
	}

	if timeout > 0 {
		t := time.Duration(timeout) * time.Second

		s.socket.SetOption(mangos.OptionRecvDeadline, t)
		s.socket.SetOption(mangos.OptionSendDeadline, t)
	}
	s.socket.SetOption(mangos.OptionMaxRecvSize, 16*1024*1024)

	logs.Debug(fmt.Sprintf("protobuffer socket dial to %s", s.addr))

	return s.socket.Dial(s.addr)
}

func (s *protobufferSocket) call(req *Request) *Response {
	if err := s.send(req); err != nil {
		logs.Error(err.Error())
		return &Response{}
	}
	if resp, err := s.recv(); err != nil {
		logs.Error(err.Error())
		return &Response{}
	} else {
		return resp
	}
}

func (s *protobufferSocket) send(req *Request) error {
	if s.socket == nil {
		return errors.New("socket is nil")
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return err

	}
	return s.socket.Send(data)
}

func (s *protobufferSocket) recv() (*Response, error) {
	resp := &Response{}
	if s.socket == nil {
		return resp, errors.New("socket is nil")
	}
	data, err := s.socket.Recv()
	if err != nil {
		return resp, err

	}
	err = proto.Unmarshal(data, resp)
	return resp, err
}

func (s *protobufferSocket) close() error {
	if s.socket == nil {
		return errors.New("socket is nil")
	}
	return s.socket.Close()
}
