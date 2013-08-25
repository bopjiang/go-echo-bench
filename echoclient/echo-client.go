/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 13-8-25
 * Time: 下午12:06
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	tcpAddress = flag.String("tcp-addr", "127.0.0.1:8080", "<addr>:<port> to connect")
	numclient  = flag.Int64("n", 10, "number of concurrent client")
)

var waitgroup sync.WaitGroup

func DoClient(num int64, conn net.Conn) {

	buf := make([]byte, 1024)

	t1 := time.Now().UnixNano()
	//strTime1 := strconv.Itoa(t1)
	strTime1 := fmt.Sprintf("t=%d\n", t1)
	conn.Write([]byte(strTime1))

	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	var t = int64(0)
	_, err = fmt.Sscanf(string(buf[:n]), "t=%d", &t)
	if err != nil {
		log.Fatal(err)
	}

	t1Replay := int64(t)
	t2 := time.Now().UnixNano()

	log.Printf("[%d], %d, %d, %dms", num, t2, (t1 - t1Replay), (t2-t1Replay)/1000/1000)
}

func ClientProc(num int64) {
	conn, err := net.Dial("tcp", *tcpAddress)
	if err != nil {
		log.Fatalf("[%d]failed to connect to %s : ", num, *tcpAddress, err)
	}

	defer conn.Close()

	for {
		DoClient(num, conn)
		time.Sleep(10 * time.Second)
	}
	waitgroup.Done()
}

func main() {
	flag.Parse()
	for i := int64(1); i <= *numclient; i = i + 1 {
		waitgroup.Add(1)
		go ClientProc(i)
		time.Sleep(100 * time.Microsecond)
	}

	waitgroup.Wait()
}
