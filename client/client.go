package client

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	grpclib "github.com/go-kratos/kratos/v2/transport/grpc"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/nartvt/go-core/middleware/jwt"
	profile "github.com/nartvt/profile-service/api/profile/v1"
	"google.golang.org/grpc"
)

var (
	Url           = ""
	TimeOut       = 120 * time.Second
	ProfileClient *profileClient
	conn          *grpc.ClientConn
)

type profileClient struct {
	profileClient profile.ProfileServiceClient
}

func InitEnvironment(profileUrl string) {
	Url = profileUrl
	initWithAddressUrl(profileUrl)
}

func initWithAddressUrl(url string) {
	var err error
	conn, err = connectGrpc(url)
	if err != nil {
		fmt.Println("Connection profile connection error : ", err.Error())
		panic(err)
	}

	ProfileClient = &profileClient{
		profileClient: profile.NewProfileServiceClient(conn),
	}
}

func CloseConnect() {
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		log.Errorf("Close connection has an error: %s", err.Error())
		return
	}
	ProfileClient = nil
	log.Info("Close connection successfully !")
}

func registerClaim(subject string) jwtlib.RegisteredClaims {
	return jwtlib.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(TimeOut)),
	}
}

func ContextWithTimeOut() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	go cancelContext(cancel)
	return ctx
}

func cancelContext(cancel context.CancelFunc) {
	time.Sleep(TimeOut)
	cancel()
	log.Info("Cancel context successfully !")
}

func optionRegisterClaim(userId string) jwt.Option {
	return jwt.WithClaims(func() jwtlib.Claims {
		return registerClaim(userId)
	})
}

func middlewareWithJwt(userId string) middleware.Middleware {
	return jwt.Client(optionRegisterClaim(userId))
}

func optionMiddlewareWithJwt(userId string) grpclib.ClientOption {
	return grpclib.WithMiddleware(
		middlewareWithJwt(userId),
	)
}

func userIntoContext(userId string) context.Context {
	ctx, err := jwt.ClientGrpcAuth(ContextWithTimeOut(), optionRegisterClaim(userId))
	if err != nil {
		log.Errorf("middlewareWithJwt: %s", err.Error())
		panic(err)
	}
	return ctx
}

func connectGrpc(url string) (*grpc.ClientConn, error) {
	return grpclib.DialInsecure(
		context.TODO(),
		grpclib.WithEndpoint(url),
		grpclib.WithTimeout(TimeOut),
	)
}

func GrpcProfileClient() *profileClient {
	if ProfileClient == nil {
		initWithAddressUrl(Url)
	}
	return ProfileClient
}

func (w profileClient) Profile() profile.ProfileServiceClient {
	return w.profileClient
}

func (w *profileClient) GetUserProfileById(userId string) (*profile.GetUserProfileResponse, error) {
	resp, err := w.profileClient.GetUserProfileInternal(userIntoContext(userId), &profile.GetUserProfileInternalRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *profileClient) GetListProfile(userIds []string) (*profile.GetListUserProfileResponse, error) {
	getListUserResp, err := w.profileClient.GetListUserProfileInternal(userIntoContext(""), &profile.GetListUserProfileInternalRequest{UserIds: userIds})
	if err != nil {
		log.Errorf("[GET_CHILD_REFERRAL] get list user ids error : %s", err.Error())
		return nil, err
	}
	return getListUserResp, nil
}
