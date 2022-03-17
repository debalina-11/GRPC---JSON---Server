package controllers

import (
	"context"
	"fmt"
	"grpc-json-server/models"
	"grpc-json-server/pb"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Register(c *gin.Context) {
	fmt.Println("I am grpc client")

	//creating the client
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("do not dial %v", err)
	}
	defer cc.Close()
	C := pb.NewSampleServiceClient(cc)

	//retriving data from json
	user := models.User{}
	err = c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "can't bind",
		})
	}
	//creting the request for grpc server
	req := &pb.RegisterUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Gender:   user.Gender,
	}

	res, err := C.RegisterUser(context.Background(), req)

	if err != nil {
		log.Fatalf("error occured when we call the sample in client %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "can't call the server",
		})
	}

	c.JSON(http.StatusOK, res)
}

func Login(c *gin.Context) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("do not dial %v", err)
	}
	defer cc.Close()
	C := pb.NewSampleServiceClient(cc)

	type body struct {
		Email    string `josn:"email"`
		Password string `json:"password"`
	}
	var reqBody body
	c.Bind(&reqBody)

	// creating the request for grpc server

	req := &pb.LoginUserRequest{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}
	res, err := C.LoginUser(context.Background(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": "can't call the server",
		})
	}
	token := res.GetToken()
	c.Writer.Header().Add("x-auth-token", "Bearer "+token)

	c.JSON(http.StatusOK, gin.H{
		"token":  token,
		"result": "login successfull",
	})
}

func About(c *gin.Context) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("do not dial %v", err)
	}
	defer cc.Close()
	C := pb.NewSampleServiceClient(cc)

	token := c.Request.Header.Get("x-auth-token")

	req := &pb.AboutUserRequest{
		Token: token,
	}

	res, err := C.AboutUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "can't call the server")
	}

	c.JSON(http.StatusOK, gin.H{
		"name":   res.Name,
		"email":  res.Email,
		"gender": res.Gender,
	})
}

func Update(c *gin.Context) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("do not dial %v", err)
	}
	defer cc.Close()
	C := pb.NewSampleServiceClient(cc)

	//retriving data from json
	user := models.User{}
	err = c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"res": "can't bind",
		})
	}

	token := c.Request.Header.Get("x-auth-token")

	req := &pb.UpdateUserRequest{
		Name:   user.Name,
		Email:  user.Email,
		Gender: user.Gender,
		Token:  token,
	}

	res, err := C.UpdateUser(context.Background(), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, "can't call the server")
	}
	c.JSON(http.StatusOK, res)
}

func Delete(c *gin.Context) {
	//creating the client
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("do not dial %v", err)
	}

	defer cc.Close()
	C := pb.NewSampleServiceClient(cc)

	token := c.Request.Header.Get("x-auth-token")

	req := &pb.DeleteUserRequest{
		Token: token,
	}

	res, err := C.DeleteUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "can't call the server")
	}

	c.JSON(http.StatusOK, res)
}
