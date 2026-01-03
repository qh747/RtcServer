package testrpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// proto生成go文件命令：protoc --go_out=. --go-grpc_out=. rpc.proto

type PushServer struct {
	UnimplementedPushRequestServer
}

func (serv *PushServer) PushRequest(ctx context.Context, req *PushReqArgs) (*PushRespArgs, error) {
	fmt.Printf("Request. room: %s. user: %s. type: %s. msg: %s\n", req.GetRoom(), req.GetUser(), req.GetType(), req.GetMsg())

	return &PushRespArgs{
		Code: 0,
		Msg:  "Success response",
	}, nil
}

func TestRpcFst(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(2)

	done := make(chan bool)

	go func() {
		defer wait.Done()

		// 注册rpc函数
		rpcServ := grpc.NewServer()
		RegisterPushRequestServer(rpcServ, &PushServer{})

		//启动服务监听
		serv, _ := net.Listen("tcp", "127.0.0.1:8087")

		go func() {
			rpcServ.Serve(serv)
		}()

		<-done

		rpcServ.GracefulStop()
		serv.Close()
	}()

	wait.Go(func() {
		defer wait.Done()
		time.Sleep(2 * time.Second)

		conn, _ := grpc.NewClient("127.0.0.1:8087", grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()

		cli := NewPushRequestClient(conn)
		resp, _ := cli.PushRequest(context.Background(), &PushReqArgs{
			Room: "room_1",
			User: "user_1",
			Type: "push",
			Msg:  "this is a push request",
		})

		fmt.Printf("Response code: %d, msg: %s\n", resp.GetCode(), resp.GetMsg())
		done <- true
	})

	wait.Wait()
}

type PushSecureServer struct {
	UnimplementedPushRequestServer
}

func (serv *PushSecureServer) PushRequest(ctx context.Context, req *PushReqArgs) (*PushRespArgs, error) {
	fmt.Printf("Request. room: %s. user: %s. type: %s. msg: %s\n", req.GetRoom(), req.GetUser(), req.GetType(), req.GetMsg())

	// token验证
	token, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &PushRespArgs{
			Code: -1,
			Msg:  "Token empty",
		}, nil
	}

	if val, ok := token["token"]; !ok || "token_rpc_secure_1" != val[0] {
		return &PushRespArgs{
			Code: -1,
			Msg:  "Token invalid",
		}, nil
	}

	return &PushRespArgs{
		Code: 0,
		Msg:  "Success response",
	}, nil
}

type AuthToken struct{}

func (auth AuthToken) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"token": "token_rpc_secure_1",
	}, nil
}

func (auth AuthToken) RequireTransportSecurity() bool {
	return false
}

func TestRpcToken(t *testing.T) {
	var wait sync.WaitGroup
	wait.Add(2)

	done := make(chan bool)

	go func() {
		defer wait.Done()

		// 注册rpc函数
		rpcServ := grpc.NewServer()
		RegisterPushRequestServer(rpcServ, &PushSecureServer{})

		//启动服务监听
		serv, _ := net.Listen("tcp", "127.0.0.1:8087")

		go func() {
			rpcServ.Serve(serv)
		}()

		<-done

		rpcServ.GracefulStop()
		serv.Close()
	}()

	wait.Go(func() {
		defer wait.Done()
		time.Sleep(2 * time.Second)

		var opts []grpc.DialOption
		// 关闭密钥认证
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// 使用自定义token
		opts = append(opts, grpc.WithPerRPCCredentials(new(AuthToken)))

		conn, _ := grpc.NewClient("127.0.0.1:8087", opts...)
		defer conn.Close()

		cli := NewPushRequestClient(conn)
		resp, _ := cli.PushRequest(context.Background(), &PushReqArgs{
			Room: "room_1",
			User: "user_1",
			Type: "push",
			Msg:  "this is a push request",
		})

		fmt.Printf("Response code: %d, msg: %s\n", resp.GetCode(), resp.GetMsg())
		done <- true
	})

	wait.Wait()
}
