package controllers

import (
	"context"
	"fmt"
	"grpc-json-server/pb"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Sample(c *gin.Context) {
	fmt.Println("I am Grpc Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("do not dial %v", err)
	}

	defer cc.Close()

	C := pb.NewSampleServiceClient(cc)
	fmt.Println("Created client is ", C)

	type Message struct {
		Msg string `json:"msg"`
	}

	var reqBody Message

	err = c.Bind(&reqBody)

	if err != nil {
		return
	}
	fmt.Println(reqBody)
	msg := DoSample(C, reqBody.Msg)
	c.JSON(http.StatusOK, gin.H{
		"result": msg,
	})
}

func DoSample(c pb.SampleServiceClient, msg string) string {
	req := &pb.SimpleRequest{
		Msg: msg,
	}
	res, err := c.Sample(context.Background(), req)

	if err != nil {
		log.Fatalf("error occured when we call the sample in client %v", err)
		return ""
	}
	return res.GetMsg()
}
