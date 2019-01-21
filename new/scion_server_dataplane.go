// Code references : https://github.com/netsec-ethz/scion-homeworks/blob/master/latency/timestamp_server.go and reference https://github.com/perrig/scionlab/blob/master/sensorapp/sensorfetcher/sensorserver.go
//Code references : https://github.com/netsec-ethz/scion-homeworks/blob/master/latency/dataplane_server.go

package  main

import (
 "fmt"
 "flag"         //import fmt for printing
 "log"                                                                      //importing log for logging out errors
 "github.com/scionproto/scion/go/lib/snet"                                  //importing snet packages
 "github.com/scionproto/scion/go/lib/sciond"                                //importing packages
 )





func logerror(ef error) {                                              //func which will be caused frequently to log out errors

	if ef!=nil {

log.Println(ef)
}
}



func main(){

var (
saddr string                                //Intialising local variables
ef error
ser *snet.Addr
scionconnection *snet.Conn
)




 flag.StringVar(&saddr, "saddr", "", "server addr")  // flag used to fetch value from command line
 flag.Parse()



 ser, ef = snet.AddrFromString(saddr)      // AddrFromString converts an address string of format isd-as,[ipaddr]:port
 logerror(ef)


 dpath := "/run/shm/dispatcher/default.sock"
	snet.Init(ser.IA, sciond.GetDefaultSCIONDPath(nil), dpath)  //Init initializes the default SCION networking context.



  scionconnection, ef = snet.ListenSCION("udp4", ser) //  ListenSCION registers laddr with the dispatcher. Nil values for laddr are
                                                          // not supported yet. The returned connection's ReadFrom and WriteTo methods
                                                          // can be used to receive and send SCION packets with per-packet addressing.
                                                          // Parameter network must be "udp4".
  logerror(ef)
receivePacketBuffer := make([]byte, 2500)  //Intiating a dynamic array of respective size

  for {
  		n, clientAddr, ef := scionconnection.ReadFrom(receivePacketBuffer) //decoding value from buffer
  	 logerror(ef)

  		// Packet received, send back response to same client

  		_, ef = scionconnection.WriteTo(receivePacketBuffer[:n], clientAddr) //Hint given in tutorial (sending back the response to client)
  		 logerror(ef)
  		fmt.Println("Scion connection from", clientAddr)
  	}



}
