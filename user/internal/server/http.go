package server

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	stdhttp "net/http"
	pb "user/api/user/v1"
	"user/internal/conf"
	"user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func jsonMarshal(res *pb.CommonReply) ([]byte, error) {
	newProto := protojson.MarshalOptions{EmitUnpopulated: true}
	output, err := newProto.Marshal(res)
	if err != nil {
		return nil, err
	}

	var stuff map[string]any
	if err := json.Unmarshal(output, &stuff); err != nil {
		return nil, err
	}

	if stuff["data"] != nil {
		delete(stuff["data"].(map[string]any), "@type")
	}
	return json.MarshalIndent(stuff, "", "  ")
}

// EncoderError 错误响应封装
func EncoderError() http.EncodeErrorFunc {
	return func(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
		if err == nil {
			return
		}
		se := &pb.CommonReply{}
		gs, ok := status.FromError(err)
		if !ok {
			se.Code = stdhttp.StatusInternalServerError
			se.Message = "Internal Server Error"
		} else {
			//se.Code = pb.sta.FromGRPCCode(gs.Code())
			se.Message = gs.Message()
		}
		codec, _ := http.CodecForRequest(r, "Accept")
		body, err := codec.Marshal(se)
		if err != nil {
			w.WriteHeader(stdhttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/"+codec.Name())
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write(body)
	}
}

func EncoderResponse() http.EncodeResponseFunc {
	return func(w stdhttp.ResponseWriter, request *stdhttp.Request, i interface{}) error {
		resp := &pb.CommonReply{
			Code:    200,
			Message: "success",
		}
		var data []byte
		var err error
		if m, ok := i.(proto.Message); ok {
			payload, err := anypb.New(m)
			if err != nil {
				return err
			}
			resp.Data = payload
			data, err = jsonMarshal(resp)
			if err != nil {
				return err
			}
		} else {
			dataMap := map[string]interface{}{
				"code":    200,
				"message": "success",
				"data":    i,
			}
			data, err = json.Marshal(dataMap)
			if err != nil {
				return err
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			return err
		}
		return nil
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, user *service.UserService, ac *conf.Auth, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(ac.JwtKey), nil
				}, jwt.WithSigningMethod(jwtv5.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}),
			handlers.AllowCredentials(),
		)),
		http.ResponseEncoder(EncoderResponse()),
		//http.ErrorEncoder(EncoderError()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	pb.RegisterUserHTTPServer(srv, user)
	return srv
}

// NewWhiteListMatcher 白名单不需要token验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/user.v1.User/Verify"] = struct{}{}
	whiteList["/user.v1.User/Login"] = struct{}{}
	whiteList["/user.v1.User/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
