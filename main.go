package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	insecurecreds "google.golang.org/grpc/credentials/insecure"
)

func main() {
	// testnet
	// insecure := true
	// target := "grpc-1.elgafar-1.stargaze-apis.com:26660"

	insecure := false
	target := "grpc.stargaze-apis.com:443"
	address := "stars1...."

	client, err := createQueryClient(target, insecure)
	if err != nil {
		panic(err)
	}
	resp, err := client.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   "ustars",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func createQueryClient(target string, insecure bool) (banktypes.QueryClient, error) {
	var conn *grpc.ClientConn
	var err error
	if insecure {
		conn, err = grpc.Dial(target, grpc.WithTransportCredentials(insecurecreds.NewCredentials()), grpc.WithContextDialer(dialerFunc))
		if err != nil {
			return nil, err
		}

		return banktypes.NewQueryClient(conn), nil
	}

	conn, err = grpc.Dial(
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		return nil, err
	}
	return banktypes.NewQueryClient(conn), nil
}

func dialerFunc(_ context.Context, addr string) (net.Conn, error) {
	return Connect(addr)
}

func Connect(protoAddr string) (net.Conn, error) {
	proto, address := ProtocolAndAddress(protoAddr)
	conn, err := net.Dial(proto, address)
	return conn, err
}

func ProtocolAndAddress(listenAddr string) (string, string) {
	protocol, address := "tcp", listenAddr

	parts := strings.SplitN(address, "://", 2)
	if len(parts) == 2 {
		protocol, address = parts[0], parts[1]
	}

	return protocol, address
}
