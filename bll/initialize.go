package bll

import (
	"DBMS/SwcDbmsCommon/Generated/go/proto/service"
	"DBMS/config"
	"DBMS/dal"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

func Initialize() {
	dal.InitializeNewDataBaseIfNotExist(dal.DefaultMetaInfoDataBaseName, dal.DefaultSwcDataBaseName, dal.DefaultSwcSnapshotDataBaseName, dal.DefaultSwcIncrementOperationDataBaseName, dal.DefaultSwcAttachmentDataBaseName)

	var createInfo dal.MongoDbConnectionCreateInfo
	createInfo.Host = config.AppConfig.MongodbIP
	createInfo.Port = config.AppConfig.MongodbPort
	createInfo.User = config.AppConfig.MongodbUser
	createInfo.Password = config.AppConfig.MongodbPassword
	connectionInfo := dal.ConnectToMongoDb(createInfo)

	if connectionInfo.Err != nil {
		log.Fatal(connectionInfo.Err)
	}

	databaseInstance := dal.ConnectToDataBase(connectionInfo, dal.DefaultMetaInfoDataBaseName, dal.DefaultSwcDataBaseName, dal.DefaultSwcSnapshotDataBaseName, dal.DefaultSwcIncrementOperationDataBaseName, dal.DefaultSwcAttachmentDataBaseName)

	dal.SetDbInstance(databaseInstance)
}

func NewGrpcServer() {
	address := config.AppConfig.GrpcIP + ":" + strconv.Itoa(int(config.AppConfig.GrpcPort))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*256), grpc.MaxSendMsgSize(1024*1024*256)) // 256mb, 256mb

	var instanceDBMSServerController DBMSServerController
	service.RegisterDBMSServer(s, instanceDBMSServerController)
	reflection.Register(s)

	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

}
