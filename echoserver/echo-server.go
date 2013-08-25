/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 13-8-25
 * Time: 上午10:50
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"bufio"
	"log"
	"net"
)

var connNum = 0

type ClientConn struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	writeChan chan []byte

	exitChan chan struct{}
}

func NewConn(conn net.Conn) *ClientConn {
	connNum = connNum + 1
	log.Print("Current connection ", connNum, " from ", conn.RemoteAddr())
	return &ClientConn{
		conn:      conn,
		reader:    bufio.NewReaderSize(conn, 4*1024),
		writer:    bufio.NewWriterSize(conn, 1*1024),
		writeChan: make(chan []byte, 1),
		exitChan:  make(chan struct{}),
	}
}

func handle(conn net.Conn) {
	var line []byte
	var err error

	clientConn := NewConn(conn)
	defer func() {
		connNum = connNum - 1
		close(clientConn.exitChan)
	}()

	go clientWriteProc(clientConn)
	for {
		line, err = clientConn.reader.ReadSlice('\n')
		if err != nil {
			log.Print("Failed to read:", clientConn.conn.RemoteAddr(), err)
			break
		}

		clientConn.writeChan <- line
		//TODO: have problem, buffer will be overwritten by next read
	}
}

func clientWriteProc(clientConn *ClientConn) {
	for {
		select {
		case data := <-clientConn.writeChan:
			clientConn.writer.Write(data)
			clientConn.writer.Flush()
		case <-clientConn.exitChan:
			goto exit
		}
	}

exit:
	clientConn.conn.Close()
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	listerner, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("Failed to listen:", err.Error())
		return
	}

	for {
		conn, err := listerner.Accept()
		if err != nil {
			log.Fatal("Failed to accept:", err.Error())
		}

		go handle(conn)
	}
}
