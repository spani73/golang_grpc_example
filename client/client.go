package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/spani73/proto_example/coffeeshop_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("Connected to server")

	defer conn.Close()

	c := pb.NewCoffeeShopClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("could not get menu: %v", err)
	}

	done := make(chan bool)

	var items []*pb.Item

	go func ()  {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			
			if err != nil {
				log.Fatalf("could not receive menu: %v", err)
			}

			items = resp.Items
			log.Printf("Items: %v", items)
		}
	}()

	<- done
	
	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	log.Printf("Receipt: %v", receipt)

	status,err := c.GetOrderStatus(ctx,receipt)
	log.Printf("Status: %v", status)
}