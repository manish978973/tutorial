//References: https://github.com/netsec-ethz/scion-homeworks/blob/master/latency/dataplane_client.go
// https://github.com/netsec-ethz/scion-apps/blob/master/sensorapp/sensorfetcher/sensorfetcher.go




package main

import (
	"flag"
	"fmt"  //import fmt for printing
	"log"  //importing log for logging out errors
	"math/rand" //importing for computing mathematical computations
	"time" // used to estimate time
  "encoding/binary"  
	"github.com/scionproto/scion/go/lib/snet" //importing snet packages
	"github.com/scionproto/scion/go/lib/sciond" //importing sciond packages
)




func logerror(ef error){    //func which will be caused frequently to log out errors

if ef!=nil{

log.Println(ef)
}
}



func main() {

  var (
	clientAddress string
    serverAddress string
     ef error
    client_local *snet.Addr
   server_destination *snet.Addr
   scionconnection *snet.Conn
  )

  flag.StringVar(&clientAddress, "c", "", "Client SCION AS address")
  flag.StringVar(&serverAddress, "s", "", "Server SCION AS Address")
  flag.Parse()


  client_local, ef = snet.AddrFromString(clientAddress)
  logerror(ef)
  server_destination, ef = snet.AddrFromString(serverAddress)
  logerror(ef)

  dpath := "/run/shm/dispatcher/default.sock"
  snet.Init(client_local.IA, sciond.GetDefaultSCIONDPath(nil), dpath)

  scionconnection, ef = snet.DialSCION("udp4", client_local, server_destination)
  logerror(ef)

  receivePacketBuffer := make([]byte, 2500) //Intiating a dynamic array of respective size
  sendPacketBuffer := make([]byte, 16)   //Intiating a dynamic array of respective size


  seed := rand.NewSource(time.Now().UnixNano())

  	var total_number int64 = 0




id := rand.New(seed).Uint64() //generating random value
n := binary.PutUvarint(sendPacketBuffer, id)  //encoding  id to buffer
sendPacketBuffer[n] = 0

 time_sent := time.Now()

 _, ef = scionconnection.Write(sendPacketBuffer) //sending response to server
 logerror(ef)

 _, _, ef = scionconnection.ReadFrom(receivePacketBuffer) //reading response from server
 time_received := time.Now()
    logerror(ef)

    ret_id, n := binary.Uvarint(receivePacketBuffer)  //decoding the id from buffer  and verifying if the packet was returned via same id
      if ret_id == id {

        diff := (time_received.UnixNano() - time_sent.UnixNano()) //change in time as per dataplane method
        total_number = diff
    

           }
  


    var difference float64 = float64(total_number) // / float64(iters)  //Taking average of latencies for precision

    fmt.Printf("\nClient: %s\nServer: %s\n", clientAddress, serverAddress);
    fmt.Println("LATENCY_DATAPLANE_METHOD:")
    // Print in ms, so divide by 1e6 from nano
    fmt.Printf("\tRTT - %.3fns\n", difference)
    fmt.Printf("\tLatency - %.3fns\n", difference/2)














}
