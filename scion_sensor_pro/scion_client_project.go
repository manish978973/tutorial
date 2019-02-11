
package main

import (
	"flag"
	"fmt"  //import fmt for printing
	"log"  //importing log for logging out errors
	"math/rand" //importing for computing mathematical computations
	"time"
  "encoding/binary"
	"github.com/scionproto/scion/go/lib/snet" //importing snet packages
	"github.com/scionproto/scion/go/lib/sciond" //importing sciond packages
)

func logerror(ef error) {   //func which will be caused frequently to log out errors

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
sendPacketBuffer := make([]byte, 2500)




//seed := rand.NewSource(time.Now().UnixNano())
var total_number int64 = 0

//id := rand.New(seed).Uint64() //generating random value
//n := binary.PutUvarint(sendPacketBuffer, id)  //encoding id to buffer
//sendPacketBuffer[n] = 0

//time_sent := time.Now()
	 n, ef = scionconnection.Write(sendPacketBuffer) //send response to serber
	 logerror(ef)

	 n, _, ef = scionconnection.ReadFrom(receivePacketBuffer)  //reading response from server
			 logerror(ef)

		//	 ret_id, n := binary.Uvarint(receivePacketBuffer)  //decoding the id anc verifying if the packet was returned via same id
		//		 if ret_id == id {
		//			 sensorvalue_received, _ := binary.Varint(receivePacketBuffer[n:]) //estimating the time received so as to compute Latency
					// diff := (time_received - time_sent.UnixNano())
		//			sensor_value := (sensorvalue_received)
		//    total_number = sensor_value


		//		 }




			// var difference float64 = float64(total_number)
		//	var sensorvalue float64 = float64(total_number)

			 fmt.Println("WEIGHT OF THE OBJECT IS:")
			 // Print in ms, so divide by 1e6 from nano
			 fmt.Print(string(receivePacketBuffer[:n]))
		//	 fmt.Printf("\tWeight is - %.3fgrams\n", sensorvalue)      //estimating RTT
		//	 fmt.Printf("\tLatency is - %.3fns\n", difference/2) //estimating Latency





}
