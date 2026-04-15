// Package client ...
package client

import (
	"context"
	"log"
	rumrpc "rum/app/misc/rum"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func POST(addr string, call []*rumrpc.IPost) {
	log.Println("in client")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
	}

	client := rumrpc.NewOnRumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	x := rumrpc.IPostRequest{
		Post: call,
	}
	if res, err := client.POST(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}

func DELETEPROFILE(addr string, call []*rumrpc.IDelete) {
	log.Println("in client")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
	}

	client := rumrpc.NewOnRumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	x := rumrpc.IDeleteRequest{
		Delete: call,
	}
	if res, err := client.DELETE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}

func ACTIVATEPROFILE(addr string, call []*rumrpc.IActivate) {
	log.Println("in client")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
	}

	client := rumrpc.NewOnRumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	x := rumrpc.IActivateRequest{
		Activate: call,
	}
	if res, err := client.ACTIVATE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}
func DEACTIVATEPROFILE(addr string, call []*rumrpc.IDeactivate) {
	log.Println("in client")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
	}

	client := rumrpc.NewOnRumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	x := rumrpc.IDeactivateRequest{
		Deactivate: call,
	}
	if res, err := client.DEACTIVATE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}

func REMOVESERVICE(addr string, call []*rumrpc.IDelete) {
	log.Println("in client")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := rumrpc.NewOnRumServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	x := rumrpc.IRemoveServiceRequest{Delete: call}
	if res, err := client.REMOVESERVICE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}

func DEACTIVATESERVICE(addr string, call []*rumrpc.IDelete) {
	log.Println("in client")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := rumrpc.NewOnRumServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	x := rumrpc.IDeactivateServiceRequest{Delete: call}
	if res, err := client.DEACTIVATESERVICE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}

func ACTIVATESERVICE(addr string, call []*rumrpc.IDelete) {
	log.Println("in client")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := rumrpc.NewOnRumServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	x := rumrpc.IActivateServiceRequest{Delete: call}
	if res, err := client.ACTIVATESERVICE(ctx, &x); err != nil {
		panic(err)
	} else {
		log.Println("req succeed: ", res.Succeed)
	}
}
