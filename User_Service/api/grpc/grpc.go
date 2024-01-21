package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/hkm15022001/Supply-Chain-Event-Management/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Implement your gRPC service here
type grpcServerStruct struct {
	pb.UnimplementedLongShipServer
}

// mustEmbedUnimplementedCalculatorServer implements pb.CalculatorServer.
func (*grpcServerStruct) mustEmbedUnimplementedCalculatorServer() {

}
func (*grpcServerStruct) mustEmbedUnimplementedLongShipServer() {

}

func extractIDs(packageListResult map[int32]*pb.PackageItems) [][]string {
	var ids [][]string

	for _, item := range packageListResult {
		ids = append(ids, item.Id)
	}

	return ids
}
func (s *grpcServerStruct) CreateLongShipFromVRP(ctx context.Context, req *pb.PackageListResult) (*pb.LongShipResponse, error) {
	// log.Print(req)
	// var data []byte // Thay thế dòng này với dữ liệu protobuf nhận được từ request
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatal("Error marshalling protobuf data: ", err)
	}
	// Giải mã dữ liệu protobuf vào một biến kiểu PackageListResult
	receivedResult := &pb.PackageListResult{}
	if err := proto.Unmarshal(data, receivedResult); err != nil {
		log.Fatal("Error unmarshalling protobuf data: ", err)
	}

	// Truy cập các giá trị trong biến receivedResult
	for _, value := range receivedResult.PackageListResult {
		itemsLength := len(value.Id)
		if itemsLength == 0 {
			return &pb.LongShipResponse{Ok: true, Response: "Nothing to do"}, nil

		}

	}
	return &pb.LongShipResponse{Ok: true, Response: "Done!"}, nil
}

// grpc server
func RunServer(grpcURL string) {
	log.Println("Starting gRPC server...")

	grpcServer := grpc.NewServer()
	pb.RegisterLongShipServer(grpcServer, &grpcServerStruct{})
	// Listen on port 50052
	lis, err := net.Listen("tcp", grpcURL)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server is starting...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

// func main() {
// 	RunServer()
// }
