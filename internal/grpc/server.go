package grpc

import (
	"context"
	"fmt"
	"github.com/thxhix/shortener/internal/config"
	pb "github.com/thxhix/shortener/internal/gen"
	"github.com/thxhix/shortener/internal/models"
	"github.com/thxhix/shortener/internal/url"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"time"

	"net"
)

// GRPCServer hosts the gRPC endpoints for the URL shortener service.
// It adapts requests to the underlying use case (business logic) and
// logs server lifecycle events via the provided logger.
type GRPCServer struct {
	pb.UnimplementedShortenerServiceServer

	cfg    config.Config
	uc     url.URLUseCaseInterface
	logger *zap.SugaredLogger
}

// New creates a new GRPCServer bound to the given configuration,
// use case implementation and logger. The returned server is ready
// to be registered into a *grpc.Server and served.
func New(cfg config.Config, uc url.URLUseCaseInterface, logger *zap.SugaredLogger) *GRPCServer {
	return &GRPCServer{
		cfg:    cfg,
		uc:     uc,
		logger: logger,
	}
}

// StartPooling starts the gRPC server listening on the address configured
// in cfg.GRPCConfig.Address. It registers all RPC handlers and serves them
// in a background goroutine.
func (s *GRPCServer) StartPooling(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.GRPCConfig.Address)
	if err != nil {
		return err
	}

	gs := grpc.NewServer()
	pb.RegisterShortenerServiceServer(gs, s)

	s.logger.Infof("gRPC server listening on %s", s.cfg.GRPCConfig.Address)

	go func() {
		if serveErr := gs.Serve(lis); serveErr != nil {
			s.logger.Errorf("gRPC Serve error: %v", serveErr)
		}
	}()

	go func() {
		<-ctx.Done()
		stopped := make(chan struct{})
		go func() {
			gs.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-time.After(5 * time.Second):
			gs.Stop()
		}
	}()

	return nil
}

// StoreLink shortens a single original URL and returns the short URL.
// On success, StatusCode is http.StatusOK and ShortUrl contains the result.
// On failure, Error is filled and StatusCode is http.StatusInternalServerError.
func (s *GRPCServer) StoreLink(ctx context.Context, req *pb.StoreLinkRequest) (*pb.StoreLinkResponse, error) {
	var response pb.StoreLinkResponse

	shorten, err := s.uc.Shorten(ctx, req.Url)
	if err != nil {
		response.Error = fmt.Sprintf("Не удалось сократить: %s", err.Error())
		response.StatusCode = http.StatusInternalServerError
	} else {
		response.ShortUrl = shorten
		response.StatusCode = http.StatusOK
	}

	return &response, nil
}

// GetFullURL resolves a short URL (ID/hash) back to its original URL.
// On success, OriginalUrl is populated. On failure, Error is filled.
func (s *GRPCServer) GetFullURL(ctx context.Context, req *pb.GetFullURLRequest) (*pb.GetFullURLResponse, error) {
	var response pb.GetFullURLResponse

	original, err := s.uc.GetFullURL(ctx, req.ShortenUrl)
	if err != nil {
		response.Error = fmt.Sprintf("Не удалось получить полную ссылку: %s", err.Error())
	} else {
		response.OriginalUrl = original
	}

	return &response, nil
}

// BatchStoreLink shortens a batch of original URLs. Each request item carries
// a correlation_id and an original_url; each response item contains the same
// correlation_id and the produced short URL.
func (s *GRPCServer) BatchStoreLink(ctx context.Context, req *pb.BatchStoreLinkRequest) (*pb.BatchStoreLinkResponse, error) {
	var response pb.BatchStoreLinkResponse
	var list models.BatchShortenRequestList

	for _, item := range req.Items {
		list = append(list, models.BatchShortenRequest{
			ID:  item.CorrelationId,
			URL: item.OriginalUrl,
		})
	}

	shorten, err := s.uc.BatchShorten(ctx, list)
	if err != nil {
		response.Error = fmt.Sprintf("Не удалось сократить ссылки: %s", err.Error())
		response.StatusCode = http.StatusInternalServerError
	} else {
		for _, item := range shorten {
			response.Items = append(response.Items, &pb.BatchResponseItem{
				CorrelationId: item.ID,
				ShortUrl:      item.Hash,
			})
		}
		response.StatusCode = http.StatusOK
	}

	return &response, nil
}

// PingDatabase checks database connectivity. It returns http.StatusOK when
// the database is reachable, otherwise http.StatusInternalServerError with
// a human-readable error message in the response.
func (s *GRPCServer) PingDatabase(ctx context.Context, req *pb.PingDatabaseRequest) (*pb.PingDatabaseResponse, error) {
	var response pb.PingDatabaseResponse
	err := s.uc.PingDB()
	if err != nil {
		response.Error = fmt.Sprintf("Нет ответа от БД: %s", err.Error())
		response.StatusCode = http.StatusInternalServerError
	} else {
		response.StatusCode = http.StatusOK
	}
	return &response, nil
}

// UserList returns all URLs that belong to a given user. On success, the
// response contains a list of {short_url, original_url} pairs and StatusOK.
// On failure, Error is populated and StatusCode is http.StatusInternalServerError.
func (s *GRPCServer) UserList(ctx context.Context, req *pb.UserListRequest) (*pb.UserListResponse, error) {
	var response pb.UserListResponse

	list, err := s.uc.UserList(ctx, req.UserId)
	if err != nil {
		response.Error = fmt.Sprintf("Не удалось получить данные по пользователю: %s", err.Error())
		response.StatusCode = http.StatusInternalServerError
	} else {
		for _, item := range list {
			response.Urls = append(response.Urls, &pb.UserURL{
				ShortUrl:    item.Short,
				OriginalUrl: item.Original,
			})
		}
		response.StatusCode = http.StatusOK
	}

	return &response, nil
}

// UserDeleteRows schedules asynchronous deletion of a batch of user URLs.
func (s *GRPCServer) UserDeleteRows(ctx context.Context, req *pb.UserDeleteRowsRequest) (*pb.UserDeleteRowsResponse, error) {
	var response pb.UserDeleteRowsResponse

	s.uc.UserDeleteRows(req.UserId, req.UrlIds, s.cfg.DeleteWorkersCount, s.cfg.DeleteBatchSize)

	return &response, nil
}

// GetServiceStats returns aggregate service statistics such as the number
// of stored URLs and total users. On success, StatusCode is http.StatusOK.
// On failure, Error is filled and StatusCode is http.StatusInternalServerError.
func (s *GRPCServer) GetServiceStats(ctx context.Context, req *pb.GetServiceStatsRequest) (*pb.GetServiceStatsResponse, error) {
	var response pb.GetServiceStatsResponse

	stats, err := s.uc.GetServiceStats(ctx)
	if err != nil {
		response.Error = fmt.Sprintf("Не удалось получить статистику: %s", err.Error())
		response.StatusCode = http.StatusInternalServerError
	} else {
		response.Urls = int32(stats.URLs)
		response.Users = int32(stats.Users)
		response.StatusCode = http.StatusOK
	}

	return &response, nil
}
