

package  main

import (
 //"fmt"             //import fmt for printing
 "log"
 "flag"
 //"encoding/binary"
 //"time"         //importing log for logging out errors
 "github.com/scionproto/scion/go/lib/snet" //importing snet packages
 "github.com/scionproto/scion/go/lib/sciond"//importing packages
 )


 func logerror(ef error){    //func which will be caused frequently to log out errors

 if ef!=nil{
  log.Println(ef)
 }
 }



 func main(){

 var (

  saddr string   //Intialising local variables
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



  receivePacketBuffer := make([]byte, 2500)
  sendPacketBuffer := make([]byte, 2500)   //Intiating a dynamic array of respective size

    for {
    		b, clientAddr, ef := scionconnection.ReadFrom(receivePacketBuffer)  //decoding value to buffer
    	 logerror(ef)

       var sensorValues int = 0.89

    		// Packet received, send back response to same client
    	//	a := binary.PutVarint(receivePacketBuffer[b:], time.Now().UnixNano())  //encoding value to buffer
    //   a := binary.PutVarint(receivePacketBuffer[b:], sensoraverage)// sensor average is the average of weight values
    copy(sendPacketBuffer, sensorValues)
  //  		_, ef = scionconnection.WriteTo(receivePacketBuffer[: b+a], clientAddr)  //sending back the response to client
        _, ef = scionconnection.WriteTo(sendPacketBuffer[:(sensorValues)], clientAddr)
    		 logerror(ef)
    //		fmt.Println("Scion connection from", clientAddr)
    	}



  }
