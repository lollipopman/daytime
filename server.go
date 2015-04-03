/* UDPDaytimeServer
*/
package main

import (
  "net"
  "time"
  "log"
)

func handleClientHeartbeats(conn *net.UDPConn, clients chan string) {

  var buf [1024]byte

  for {
    n, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
      log.Print("Error: Received bad UDP packet\n")
    } else if string(buf[0:n]) != "HEARTBEAT" {
      log.Print("Error: Received packet without a heatbeat message", string(buf[0:n]), "\n")
    } else {
      clients <- addr.String()
    }
  }
}

func serveDaytimes(conn *net.UDPConn, clients chan string) {
    var activeClients map[string]time.Time
    activeClients = make(map[string]time.Time)

    for {
      select {
      case clientAddr := <-clients:
        activeClients[clientAddr] = time.Now()
      default:
        pruneDeadClients(activeClients)
        for clientAddr, _ := range activeClients {
          daytime := time.Now().String()
          udpAddr, _ := net.ResolveUDPAddr("udp4", clientAddr)
          conn.WriteToUDP([]byte(daytime), udpAddr)
        }
        time.Sleep(10 * time.Millisecond)
      }
    }
}

func pruneDeadClients(activeClients map[string]time.Time) {
  for clientAddr, lastHeartbeat := range activeClients {
    now := time.Now()
    lastHeartbeatDuration := now.Sub(lastHeartbeat)
    if lastHeartbeatDuration.Seconds() >= 10 {
      delete(activeClients, clientAddr)
    }
  }
}

func checkError(err error) {
  if err != nil {
    log.Fatalf("Fatal error ", err.Error())
  }
}

func main() {
  service := ":1200"
  udpAddr, err := net.ResolveUDPAddr("udp4", service)
  checkError(err)

  conn, err := net.ListenUDP("udp", udpAddr)
  checkError(err)

  clients := make(chan string)
  go handleClientHeartbeats(conn, clients)
  serveDaytimes(conn, clients)
}
