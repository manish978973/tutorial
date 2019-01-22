 // reference : https://github.com/netsec-ethz/scion-homeworks/blob/master/bottleneck_bw_est/v1_bw_est_client.go and sensor reciever.go


package main

import(

  "flag"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/spath"
	"github.com/scionproto/scion/go/lib/spath/spathmeta"
)




type Checkpoint struct {
	Sent_struc, Rec_struc int64
}



func logerror(ef error) { //fn to log errors
	if ef != nil {
		log.Println(ef)
	}
}

var (
	crypthash map[uint64]*Checkpoint         //initialising map datatype
	scionconnection *snet.Conn              //configuring scion connection

)


func send() {

	var ef error
	sendPacketBuffer := make([]byte, 2500 ) //intialising a buffer array

	seed := rand.NewSource(time.Now().UnixNano()) //Fetching current time
	iters := 0
	for iters < (5) {  //Iterating
		iters += 1

		id := rand.New(seed).Uint64()
		_ = binary.PutUvarint(sendPacketBuffer, id) //wncoding the value of id in sendpacketbuffer

		crypthash[id] = &Checkpoint{time.Now().UnixNano(), 0} //configuring 	Sent_struc variable in Checkpoint with time.now and Rec_struc as 0
		_, ef = scionconnection.Write(sendPacketBuffer) //sending packet via scion connection
		check(ef)

	}
}




func AvgBW() (float64, float64) {

	sorted := make([]*Checkpoint, 5)
	i := 0
	for _, a := range crypthash { //sorting values in crypthash
		if a.Rec_struc != 0 {
			sorted[i] = a
			i += 1
		}
	}


	var sent_int, recvd_int int64

	for i := 1; i < 5; i+=1 {
		sent += (sorted[i].Sent_struc - sorted[i-1].Sent_struc) //caclulating the difference of sorted sent packets
		recvd += (sorted[i].Rec_struc - sorted[i-1].Rec_struc) //calcuating the difference of sorted receievd packets
	}
	// Calculate BW = (#Bytes*8 / #nanoseconds) / 1e6
	bw_S := float64(2500*8*1e3) / (float64(sent) / float64(5-1))  //calculating the average of sent packets
	bw_R := float64(2500*8*1e3) / (float64(recvd) / float64(5-1)) //calcuating the average of received packets

	return bw_S, bw_R   //returning
}




func receive() int {

	var ef error
	receivePacketBuffer := make([]byte, 2500 )  //intialising the array


	num := 0
	for num < 5 {
		_, _, ef = scionconnection.ReadFrom(receivePacketBuffer) //reading packet
		if (ef != nil) {
			break
		}
		ret_id, n := binary.Uvarint(receivePacketBuffer)
	                                                                      //not sure how to check the equality of ids
			time_recvd, _ := binary.Varint(receivePacketBuffer[n:])
			val.recvd = time_recvd
			num += 1

	}
	return num
}


func main() {
	var (
		sourceAddress string
		destinationAddress string

		ef    error
		local_client  *snet.Addr
		remote_server *snet.Addr

	)


	flag.StringVar(&clientAddress, "c", "", "CLient SCION Addr") //reading from command prompt
	flag.StringVar(&serverAddress, "s", "", "Server SCION Addr")
	flag.Parse()



		local_client, ef = snet.AddrFromString(clientAddress)  // // AddrFromString converts an address string of format isd-as,[ipaddr]:port
		logerror(ef)


		remote_server, ef = snet.AddrFromString(destinationAddress)
		logerror(ef)


	dpath := "/run/shm/dispatcher/default.sock"
	snet.Init(local_client.IA, sciond.GetDefaultSCIONDPath(nil), dpath)

	sciovonnection, ef = snet.DialSCION("udp4", local_client, remote_server) //Establishing scion connection netwen client and server
	check(ef)

	crypthash = make(map[uint64]*Checkpoint)

	send()
	num := receive()

	fmt.Println("Total packets are:", num)



	bw_sent, bw_recvd := getAverageBottleneckBW()

	fmt.Println("BW BANDWIDTH:")
	fmt.Printf("\tBW - %.3fMbps\n", bw_recvd) //printing bandiwdth
}
