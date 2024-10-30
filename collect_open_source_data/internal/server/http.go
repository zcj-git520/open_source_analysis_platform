package server

import (
	v1 "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/conf"
	"collect_open_source_data/internal/service"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/handlers"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	stdhttp "net/http"
)

func jsonMarshal(res *v1.CommonReply) ([]byte, error) {
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

func EncoderResponse() http.EncodeResponseFunc {
	return func(w stdhttp.ResponseWriter, request *stdhttp.Request, i interface{}) error {
		resp := &v1.CommonReply{
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
func NewHTTPServer(c *conf.Server, greeter *service.OpenSourceService, ac *conf.Auth, logger log.Logger) *http.Server {
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
	v1.RegisterOpenSourceHTTPServer(srv, greeter)
	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/open_source.v1.OpenSource/RepoFav"] = struct{}{}
	whiteList["/open_source.v1.OpenSource/GetRepoFav"] = struct{}{}
	whiteList["/open_source.v1.OpenSource/GetMessage"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return true
		}
		return false
	}
}
