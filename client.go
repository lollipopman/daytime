/* UDPDaytimeClient
 */
package main

import (
  "net"
  "os"
  "log"
  "fmt"
  "time"
)

func receiveDaytimes(conn net.Conn) {
  var buf [1024]byte

  for {
    n, err := conn.Read(buf[0:])
    checkError(err)
    fmt.Println(string(buf[0:n]))
  }
}

func sendHeartbeats(conn net.Conn) {
  for {
    _, err := conn.Write([]byte("HEARTBEAT"))
    checkError(err)
    time.Sleep(2 * time.Second)
  }
}

func checkError(err error) {
  if err != nil {
    log.Fatalf("Fatal error ", err.Error())
  }
}

func main() {
  if len(os.Args) != 2 {
    log.Fatalf("Usage: %s host:port", os.Args[0])
  }
  service := os.Args[1]

  udpAddr, err := net.ResolveUDPAddr("udp4", service)
  checkError(err)

  conn, err := net.DialUDP("udp", nil, udpAddr)
  checkError(err)

  go sendHeartbeats(conn)

  receiveDaytimes(conn)
}

