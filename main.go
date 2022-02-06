package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"

	"github.com/valcanobacon/StreamWall/api"
	"github.com/valcanobacon/StreamWall/money"
	"github.com/valcanobacon/StreamWall/music"
)

var (
	macFile = flag.String("macaroon", "readonly.macaroon", "The file container the macaroon")
	tlsFile = flag.String("tls-cert-file", "tls.cert", "The file container the cert file")
	lndAddr = flag.String("lndAddr", "localhost:10009", "The server address in the format of host:port")

	apiPort = flag.Int("port", 8080, "Port to run the api server on")

	durationFiles      = "songs/*/durations.txt"
	durationFilePrefix = "songs/"
)

func main() {
	flag.Parse()

	durations, err := music.LoadDurations(durationFiles, durationFilePrefix)
	if err != nil {
		log.Fatal(err)
	}

	bank := money.NewBank()

	ptCtx, ptCancel := context.WithCancel(context.Background())
	defer ptCancel()

	log.Printf("Starting Transaction Processor")
	go bank.ProcessTransactions(ptCtx)

	conn, err := newGrpcConn(*lndAddr, *tlsFile, *macFile)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	pp := &money.PaymentProcessor{
		Client: lnrpc.NewLightningClient(conn),
		Bank:   bank,
		Index:  0,
	}

	ppCtx, ppCancel := context.WithCancel(context.Background())
	defer ppCancel()

	log.Printf("Starting Lightning Payment Processor")
	go pp.Run(ppCtx)

	router := api.NewRouter(bank, durations)

	log.Printf("Starting web server on %v\n", *apiPort)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *apiPort), router)
	if err != nil {
		log.Fatal((err))
	}
}

func newGrpcConn(addr, tlsFile, macFile string) (*grpc.ClientConn, error) {
	tlsCert, err := credentials.NewClientTLSFromFile(tlsFile, "")
	if err != nil {
		return nil, err
	}

	macCred, err := macCredFromFile(macFile)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCert),
		grpc.WithPerRPCCredentials(macCred),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// macCredFromFile loads the macaroon from a file
func macCredFromFile(macFile string) (*macaroons.MacaroonCredential, error) {
	bs, err := ioutil.ReadFile(macFile)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(bs); err != nil {
		return nil, err
	}

	cred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		return nil, err
	}

	return &cred, nil
}
