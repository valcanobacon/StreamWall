package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"

	"github.com/valcanobacon/StreamWall/money"
	"github.com/valcanobacon/StreamWall/music"
)

var durations map[string]float64
var bank *money.Bank

var (
	macFile    = flag.String("macaroon", "readonly.macaroon", "The file container the macaroon")
	tlsFile    = flag.String("tls-cert-file", "tls.cert", "The file container the cert file")
	serverAddr = flag.String("addr", "localhost:10009", "The server address in the format of host:port")
	//serverPort = flag.Int("port", 10009)
)

func main() {
	flag.Parse()

	// configure the songs directory name and port
	const songsDir = "songs"
	const port = 8080

	bank = money.NewBank()
	bank.SetSession(uuid.Nil, 50)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go bank.ProcessTransactions(ctx)

	durations = music.LoadDurations("songs/*/durations.txt")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/sessions", func(r chi.Router) {
		r.Post("/", sessionCreateHandler)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Get("/", sessionGetHandler)
			r.Route("/streams", func(r chi.Router) {
				r.Get("/*", streamHandler)
			})
		})
	})

	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	//_, err := credentials.NewClientTLSFromFile(data.Path(tlsFile))
	tlsCert, err := credentials.NewClientTLSFromFile(*tlsFile, "")
	if err != nil {
		log.Fatal(err)
	}

	macCred, err := macCredFromFile((*macFile))
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCert),
		grpc.WithPerRPCCredentials(macCred),
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	pp := &money.PaymentProcessor{
		Client: lnrpc.NewLightningClient(conn),
		Bank:   bank,
		Index:  0,
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	go pp.Run(ctx)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func macCredFromFile(macFile string) (*macaroons.MacaroonCredential, error) {
	macBytes, err := ioutil.ReadFile(macFile)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		return nil, err
	}

	macCred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		return nil, err
	}

	return &macCred, nil
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	prefix := "/sessions/" + sessionID + "/streams"
	stream := strings.Replace(r.URL.String(), prefix, "", 1)

	if strings.HasSuffix(r.URL.String(), ".ts") {

		sid, err := uuid.Parse(chi.URLParam(r, "sessionID"))
		if err != nil {
			return
		}

		s := bank.GetSession(sid)
		if s == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		duration := durations[stream]
		satsPerSecond := 1.0
		cost := int64(duration * satsPerSecond)

		if s.Credits < cost {
			http.Error(w, http.StatusText(403), 403)
			return
		}

		log.Printf("%v %g seconds at %g costs %g", r.URL, duration, satsPerSecond, cost)

		bank.NewTransaction(sid, (-1 * cost))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fileServer := http.FileServer(http.Dir("songs"))
	h := http.StripPrefix(prefix, fileServer)
	h.ServeHTTP(w, r)
}

func sessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	s := bank.NewSession()
	render.Render(w, r, s)
}

func sessionGetHandler(w http.ResponseWriter, r *http.Request) {
	sid, err := uuid.Parse(chi.URLParam(r, "sessionID"))
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
	}
	s := bank.GetSession(sid)
	if s == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	render.Render(w, r, s)
}
