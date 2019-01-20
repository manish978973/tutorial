// Code references : https://github.com/netsec-ethz/scion-homeworks/blob/master/latency/timestamp_client.go and reference https://github.com/perrig/scionlab/blob/master/sensorapp/sensorfetcher/sensorfetcher.go


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



func printUsage() {
	fmt.Println("\ntimestamp_client -c SourceSCIONAddress -s DestinationSCIONAddress")
}



func logerror(ef error) {   //func which will be caused frequently to log out errors

if ef!=nil{

log.Println(ef)
}
}

const (
	TOTAL_NUM_ITERS = 30   //number of iterations
	TOTAL_MAX_NUM_TRIES = 50  //max limit of iterations
)


func main() {

var (

	clientAddress string
 serverAddress string
 ef error
 client_local *snet.Addr
 server_destination *snet.Addr
 scionconnection *snet.Conn
)

// We use flag for fetching values from command prompt
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
  iters := 0  // number of iterations
	num_tries := 0 //no of attempts


  for iters < TOTAL_NUM_ITERS && num_tries < TOTAL_MAX_NUM_TRIES{

 num_tries= num_tries+1
 id := rand.New(seed).Uint64() //generating random value
 n := binary.PutUvarint(sendPacketBuffer, id)  //encoding to compute id
 sendPacketBuffer[n] = 0

 time_sent := time.Now()
		_, ef = scionconnection.Write(sendPacketBuffer)
		logerror(ef)

    _, _, ef = scionconnection.ReadFrom(receivePacketBuffer)
    		logerror(ef)

        ret_id, n := binary.Uvarint(receivePacketBuffer)  //decoding the id anc verifying if the packet was returned via same id
       		if ret_id == id {
       			time_received, _ := binary.Varint(receivePacketBuffer[n:]) //estimating the time received so as to compute Latency
       			diff := (time_received - time_sent.UnixNano())
       			total_number += diff
       			iters += 1
       		}
       	}

       	if iters != TOTAL_NUM_ITERS {
       	logerror(fmt.Errorf("Error, exceeded attempts max"))
       	}

       	var difference float64 = float64(total_number) / float64(iters)  //Taking average of latencies for precision

       	fmt.Printf("\nClient: %s\nServer: %s\n", clientAddress, serverAddress);
       	fmt.Println("LATENCY_TIMESTAMP_METHOD:")
       	// Print in ms, so divide by 1e6 from nano
       	fmt.Printf("\tRTT - %.3fms\n", difference/1e9)
       	fmt.Printf("\tLatency - %.3fms\n", difference/2e9)





}
