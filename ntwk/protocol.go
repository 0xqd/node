package net

// protocol layer
//protocol operates on messages and channels not bytestreams

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	block "github.com/polygonledger/node/block"
	cryptoutil "github.com/polygonledger/node/crypto"
)

const (
	// Port is the port number that the server listens to.
	Server_address string = "127.0.0.1"
	Port                  = ":8888"
	//TODO move to message type
	Genesis_Address string = "P0614579c42f2"
	DELIM           byte   = '|'
	DELIM_HEAD      byte   = '#'
	EMPTY_MSG              = "EMPTY"
	ERROR_READ             = "error_read"
)

//given a sream read from it
//TODO proper error handling
func NReadStream(rw *bufio.ReadWriter) string {
	msg, err := rw.ReadString(DELIM)
	if err != nil {
		//issue
		//special case is empty message if client disconnects?
		if len(msg) == 0 {
			log.Println("empty message")
			return EMPTY_MSG
		} else {
			log.Println("Failed ", err)
			//log.Println(err.)
			return ERROR_READ
		}
	}
	return msg
}

func NReadMessage(rw *bufio.ReadWriter) Message {
	var msg Message
	msgString := NReadStream(rw)
	if msgString == EMPTY_MSG {
		return EmptyMsg()
	}
	msg = ParseMessage(msgString)
	return msg
}

//we have a message string and will parse it to a message struct
//delimiters of two kind:
//* DELIM for delimiting the entire message
//* DELIM_HEAD
//currently we employ delimiters instead of byte encoding, so the size of messages is unlimited
//can however easily fix by adding size to header and reject messages larger than maximum size
func ParseMessage(msgString string) Message {
	msgString = strings.Trim(msgString, string(DELIM))
	s := strings.Split(msgString, string(DELIM_HEAD))

	//ERROR handling of malformed messages

	var msg Message
	msg.MessageType = s[0]
	msg.Command = s[1]
	dataJson := s[2] //data can empty but still we expect the delim to be there

	msg.Data = []byte(dataJson)
	fmt.Println(msg)
	return msg
}

func NReadMsg(rw *bufio.ReadWriter) string {
	//TODO handle err
	msg, _ := rw.ReadString(DELIM)
	msg = strings.Trim(msg, string(DELIM))
	return msg
}

func EncodeReply(resp string) string {
	//TODO! header missing
	msg := EncodeMsg(REP, resp, "")
	return msg
}

func EncodeMessageTx(txJson []byte) string {
	//emptyData := ""
	msg := EncodeMsg(REQ, CMD_TX, string(txJson))
	return msg
}

//TODO
func NRequestReply(rw *bufio.ReadWriter, req_msg string) string {
	//REQUEST
	//TODO msg type
	NWrite(rw, req_msg)

	//REPLY
	resp_msg := NReadMsg(rw)

	return resp_msg
}

func NWrite(rw *bufio.ReadWriter, message string) error {
	n, err := rw.WriteString(message)
	if err != nil {
		return errors.Wrap(err, "Could not write data ("+strconv.Itoa(n)+" bytes written)")
	} else {
		log.Println(strconv.Itoa(n) + " bytes written")
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil
}

func ConstructMessage(cmd string) string {
	msg := cmd + string(DELIM)
	return msg
}

//TODO
//request account address
// func NReceiveAccount(rw *bufio.ReadWriter) error {
// 	log.Println("RequestAccount ", CMD_RANDOM_ACCOUNT)
// 	return nil
// }

//generic request<->reply
func RequestReplyChan(request string, msg_in_chan chan string, msg_out_chan chan string) string {
	msg_in_chan <- request
	resp := <-msg_out_chan
	return resp
}

//TODO this is higher level and should be somewhere else
func RandomTx(account_s block.Account) block.Tx {
	// s := cryptoutil.RandomPublicKey()
	// address_s := cryptoutil.Address(s)
	// account_s := block.AccountFromString(address_s)
	// log.Printf("%s", s)

	//FIX
	//doesn't work on client side
	//account_r := chain.RandomAccount()

	rand.Seed(time.Now().UnixNano())
	randNonce := rand.Intn(100)

	kp := cryptoutil.PairFromSecret("test111??")
	log.Println("PUBKEY ", kp.PubKey)

	r := cryptoutil.RandomPublicKey()
	address_r := cryptoutil.Address(r)
	account_r := block.AccountFromString(address_r)

	//TODO make sure the amount is covered by sender
	rand.Seed(time.Now().UnixNano())
	randomAmount := rand.Intn(20)

	log.Printf("randomAmount ", randomAmount)
	log.Printf("randNonce ", randNonce)
	testTx := block.Tx{Nonce: randNonce, Sender: account_s, Receiver: account_r, Amount: randomAmount}
	sig := cryptoutil.SignTx(testTx, kp)
	sighex := hex.EncodeToString(sig.Serialize())
	testTx.Signature = sighex
	log.Println(">> ran tx", testTx.Signature)
	return testTx
}

//request random account address
// func RequestRandomAccount(rw *bufio.ReadWriter) error {
// 	msg := ConstructMessage(CMD_RANDOM_ACCOUNT)
// 	error := NWrite(rw, msg)
// 	return error
// }
