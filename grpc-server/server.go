package main

import (
	"context"
	"fmt"
	"grpc-json-server/config"
	"grpc-json-server/database"
	"grpc-json-server/models"
	"grpc-json-server/pb"
	"grpc-json-server/routers"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var cfg config.Configuration

type server struct{}

//middleware or interceptor-----------------------------------------------------------------------------------
//unary middleware---------------------------------------------------------------------------

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("couldn't parse incoming context metadata")
	}
	fmt.Println(info.FullMethod)
	if info.FullMethod == "/SampleService/Sample" {
		msg := req.(*pb.SimpleRequest)
		fmt.Println(msg)
		fmt.Println("---> unary interceptor")
	}
	if info.FullMethod == "/SampleService/AboutUser" {
		request := req.(*pb.AboutUserRequest)
		token := request.GetToken()
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
		claims := jwt.MapClaims{}

		_, err1 := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if err1 != nil {
			return nil, err1
		}
		fmt.Println(claims["email"])
		md.Append("email", fmt.Sprint(claims["email"]))
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	if info.FullMethod == "/SampleService/UpdateUser" {
		request := req.(*pb.UpdateUserRequest)
		token := request.GetToken()
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
		claims := jwt.MapClaims{}

		_, err1 := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if err1 != nil {
			return nil, err1
		}
		fmt.Println(claims["email"])
		md.Append("email", fmt.Sprint(claims["email"]))
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	if info.FullMethod == "/SampleService/DeleteUser" {
		request := req.(*pb.DeleteUserRequest)
		token := request.GetToken()

		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}

		claims := jwt.MapClaims{}
		_, err1 := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if err1 != nil {
			return nil, err1
		}
		fmt.Println(claims["email"])
		md.Append("email", fmt.Sprint(claims["email"]))
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	return handler(ctx, req)
}

//stream middleware----------------------------------------------------------------------------
func streamIntercptor(srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	fmt.Println("---> stream interceptor")
	return handler(srv, stream)

}

func (*server) Sample(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	msg := req.GetMsg()

	res := &pb.SimpleResponse{
		Msg: fmt.Sprintf("Hello %s", msg),
	}

	return res, nil
}

func (*server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user := models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Gender:   req.GetGender(),
	}
	//bcrypt the password

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't generting the hashed password")
	}
	user.Password = string(hashedPass)
	result := database.Database.Db.Model(&models.User{}).Create(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "can't create the user")
	}
	res := &pb.RegisterUserResponse{
		Result: "Successfully register the user",
	}
	return res, nil
}

func (*server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	user := models.User{}

	result := database.Database.Db.Model(&models.User{}).Where("email=?", email).Find(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "error occured fetching the database")
	}

	if user.Email == "" {
		return nil, status.Errorf(codes.Internal, "user not found ")
	}

	//compare password------------------------------------------

	err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err1 != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	//creating the token------------------------------------------

	err2 := cleanenv.ReadEnv(&cfg)
	if err2 != nil {
		return nil, status.Errorf(codes.Internal, "can't read the cfg")
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["email"] = user.Email
	atClaims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(cfg.JwtSecret))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't create the token")
	}

	res := &pb.LoginUserResponse{
		Token: token,
	}
	return res, nil
}

func (*server) AboutUser(ctx context.Context, req *pb.AboutUserRequest) (*pb.AboutUserResponse, error) {

	user := models.User{}

	// Get the metadata from the incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("couldn't parse incoming context metadata")
	}
	email := md.Get("email")
	fmt.Println("The Email is: ", email)

	result := database.Database.Db.Model(&models.User{}).Where("email = ?", email).Find(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("database error")
	}

	res := &pb.AboutUserResponse{
		Name:   user.Name,
		Email:  user.Email,
		Gender: user.Gender,
	}

	return res, nil
}

func (*server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := models.User{
		Name:   req.GetName(),
		Email:  req.GetEmail(),
		Gender: req.GetGender(),
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("couldn't parse incoming context metadata")
	}
	email := md.Get("email")
	fmt.Println("The Email is: ", email)

	data := models.User{}

	result := database.Database.Db.Model(&models.User{}).Where("email = ?", email).Find(&data).Updates(map[string]interface{}{"name": user.Name, "email": user.Email, "gender": user.Gender})

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "can't Update the user")
	}

	res := &pb.UpdateUserResponse{
		Result: "Successfully Upadte the user",
	}
	return res, nil
}

func (*server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {

	// Get the metadata from the incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("couldn't parse incoming context metadata")
	}
	email := md.Get("email")
	fmt.Println(email)

	resu := database.Database.Db.Model(models.User{}).Where("email=?", email).Delete(models.User{})
	if resu.Error != nil {
		return nil, fmt.Errorf("database error")
	}

	res := &pb.DeleteUserResponse{
		Result: "delete successfully",
	}

	return res, nil
}

func main() {

	//if code crash then it log the file name with line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//connecting database
	database.Connect()
	sqlDB, err := database.Database.Db.DB()

	if err != nil {
		panic(err.Error())
	}

	//starting the grpc server
	fmt.Println("grpc server started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Not Listened %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamIntercptor),
	)

	pb.RegisterSampleServiceServer(s, &server{})

	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("not served %v", err)
		}
	}()

	//starting the http gin server
	routers.Start()

	//wait for ctrl c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//block untill enter ctrl c
	<-ch
	fmt.Println("closing the database")
	sqlDB.Close()
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("closing the listener")
	lis.Close()
	fmt.Println("End of Program")
}
