package main

import (
	"flag"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
//	"github.com/scionproto/scion/go/lib/spath"
//	"github.com/scionproto/scion/go/lib/spath/spathmeta"
)

const (
	TOTAL_PACKET_SIZE int = 1500     //setting max packet size
	TOTAL_PACKET_NUM int = 5         //setting packet number
)



func logerror(ef error) {   //func which will be caused frequently to log out errors

if ef!=nil{

log.Println(ef)
}
}


func printUsage() {
	fmt.Println("\nbw_est_client -s SourceSCIONAddress -d DestinationSCIONAddress")
}


func main() {

  var (
  		clientAddress string //local variables
  		serverAddress string
   	  scionconnection *snet.Conn
  		ef    error
  		client  *snet.Addr
  		server  *snet.Addr

  	)


      flag.StringVar(&sourceAddress, "s", "", "SCION source ad") // flag used to fetch value from command line
    	flag.StringVar(&destinationAddress, "d", "", "DESTINATION source ad")
    	flag.Parse()

client_local, ef = snet.AddrFromString(clientAddress)
server_remote, ef = snet.AddrFromString(serverAddress)

  dpath := "/run/shm/dispatcher/default.sock"
	snet.Init(client_local.IA, sciond.GetDefaultSCIONDPath(nil), dpath)

  scionconnection, ef = snet.DialSCION("udp4", client_local, server_remote)
	check(ef)



  var err error                                             //sending packet
  	sendPacketBuffer := make([]byte, 40000)

  n := binary.PutUvarint(sendPacketBuffer, id)  //encoding to compute id
  sendPacketBuffer[n] = 'packet'

    _, ef = scionconnection.Write(sendPacketBuffer)
    logerror(ef)






// The code is not complete as I am not  sure how to configure (recieved packets)  and calculating avergae band width

}
