package main

import (
	"TechSchoolGRPC/client"
	"TechSchoolGRPC/pb"
	"TechSchoolGRPC/sample"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)


//func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {
//	req := &pb.CreateLaptopRequest{Laptop:laptop}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	defer cancel()
//
//	res, err := laptopClient.CreateLaptop(ctx, req)
//	if err !=nil {
//		st, ok := status.FromError(err)
//		if ok && st.Code() == codes.AlreadyExists {
//			log.Print("Laptop already exists")
//		} else {
//			log.Fatal("cannot create laptop ", err)
//		}
//		return
//	}
//	log.Printf("laptop created with id %v", res.Id)
//}

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptop := sample.NewLaptop()
		//createLaptop(laptopClient, laptop)
		laptopClient.CreateLaptop(laptop)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopClient.SearchLaptop(filter)
}

//func rateLaptop(laptopClient pb.LaptopServiceClient, laptopIDs []string, scores []float64) error {
//	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	defer cancel()
//
//	stream, err := laptopClient.RateLaptop(ctx)
//	if err != nil {
//		return fmt.Errorf("cannot rate laptop: %v", err)
//	}
//
//	waitResponse := make(chan error)
//	go func() {
//		for ;; {
//			res, err := stream.Recv()
//			if err == io.EOF {
//				log.Printf("no more responses")
//				waitResponse <- nil
//				return
//			}
//
//			if err != nil {
//				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
//				return
//			}
//
//			log.Println("received response: ", res)
//		}
//	}()
//
//	for i, laptopID := range laptopIDs {
//		req := &pb.RateLaptopRequest{
//			LaptopId: laptopID,
//			Score:    scores[i],
//		}
//		err := stream.Send(req)
//		if err !=nil {
//			return fmt.Errorf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
//		}
//		log.Print("sent requests: ", req)
//	}
//
//	err = stream.CloseSend()
//	if err != nil {
//		return fmt.Errorf("cannot close send: %v", err)
//	}
//	fmt.Println("Waiting for rate finish")
//	err = <-waitResponse
//	return err
//}

//func uploadImage(laptopClient pb.LaptopServiceClient, laptopID string, imagePath string) {
//	file, err := os.Open(imagePath)
//	if err != nil {
//		log.Fatal("cannot open image file: ", err)
//	}
//	defer file.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	defer cancel()
//
//	stream, err := laptopClient.UploadImage(ctx)
//	if err != nil {
//		log.Fatal("cannot upload image: ", err)
//	}
//
//	req := &pb.UploadImageRequest{
//		Data: &pb.UploadImageRequest_Info{
//			Info: &pb.ImageInfo{
//				LaptopId:  laptopID,
//				ImageType: filepath.Ext(imagePath),
//			},
//		},
//	}
//
//	err = stream.Send(req)
//	if err != nil {
//		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
//	}
//
//	reader := bufio.NewReader(file)
//	buffer := make([]byte, 1024)
//
//	for {
//		n, err := reader.Read(buffer)
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatal("cannot read chunk to buffer: ", err)
//		}
//
//		req := &pb.UploadImageRequest{
//			Data: &pb.UploadImageRequest_ChunkData{
//				ChunkData: buffer[:n],
//			},
//		}
//
//		err = stream.Send(req)
//		if err != nil {
//			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
//		}
//	}
//
//	res, err := stream.CloseAndRecv()
//	if err != nil {
//		log.Fatal("cannot receive response: ", err)
//	}
//
//	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
//}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(),  "tmp/laptop.jpeg")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i:=0; i<n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n) ? ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i:= 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}
		err :=laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
		//err := rateLaptop(laptopClient, laptopIDs, scores)
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
}

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const laptopServicePath = "/techschool.pcbook.LaptopService/"

	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(*serverAddress, grpc.WithInsecure(),
						grpc.WithUnaryInterceptor(interceptor.Unary()),
						grpc.WithStreamInterceptor(interceptor.Stream()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}
	laptopClient := client.NewLaptopClient(cc2)
	testCreateLaptop(laptopClient)
	testSearchLaptop(laptopClient)
	testUploadImage(laptopClient)
	testRateLaptop(laptopClient)
}
