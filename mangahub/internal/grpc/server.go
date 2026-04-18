package grpcservice

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "mangahub/proto"
)

// MangaGRPCServer implements the MangaService gRPC server
type MangaGRPCServer struct {
	pb.UnimplementedMangaServiceServer
	db *sql.DB
}

func NewMangaGRPCServer(db *sql.DB) *MangaGRPCServer {
	return &MangaGRPCServer{db: db}
}

func (s *MangaGRPCServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	log.Printf("[gRPC] Ping received: %s", req.Message)
	return &pb.PongResponse{
		Message:  "pong: " + req.Message,
		ServerId: "mangahub-grpc-v1",
	}, nil
}

func (s *MangaGRPCServer) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, title, author, status, chapter_count, cover_url, rating, description
		 FROM manga WHERE id=$1`, req.Id)

	resp := &pb.MangaResponse{}
	err := row.Scan(&resp.Id, &resp.Title, &resp.Author, &resp.Status,
		&resp.ChapterCount, &resp.CoverUrl, &resp.Rating, &resp.Description)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("manga not found: %s", req.Id)
	}
	return resp, err
}

func (s *MangaGRPCServer) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.MangaListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, author, status, chapter_count, cover_url, rating, description
		 FROM manga WHERE title ILIKE $1 ORDER BY popularity_rank LIMIT $2 OFFSET $3`,
		"%"+req.Query+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &pb.MangaListResponse{Page: page}
	for rows.Next() {
		m := &pb.MangaResponse{}
		rows.Scan(&m.Id, &m.Title, &m.Author, &m.Status,
			&m.ChapterCount, &m.CoverUrl, &m.Rating, &m.Description)
		resp.Manga = append(resp.Manga, m)
		resp.Total++
	}
	return resp, rows.Err()
}

func (s *MangaGRPCServer) GetUserStats(ctx context.Context, req *pb.UserStatsRequest) (*pb.UserStatsResponse, error) {
	resp := &pb.UserStatsResponse{UserId: req.UserId}
	err := s.db.QueryRowContext(ctx, `
		SELECT
		  COUNT(*) FILTER (WHERE 1=1)               AS total,
		  COUNT(*) FILTER (WHERE status='reading')   AS reading,
		  COUNT(*) FILTER (WHERE status='completed') AS completed,
		  COALESCE(SUM(current_chapter),0)           AS chapters,
		  COALESCE(AVG(rating) FILTER (WHERE rating IS NOT NULL),0) AS avg_rating
		FROM library_entries WHERE user_id=$1`, req.UserId,
	).Scan(&resp.TotalManga, &resp.Reading, &resp.Completed,
		&resp.TotalChapters, &resp.AverageRating)
	return resp, err
}

// StartGRPCServer starts the gRPC server on the given port
func StartGRPCServer(port string, db *sql.DB) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("gRPC listen failed: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterMangaServiceServer(s, NewMangaGRPCServer(db))
	reflection.Register(s) // enables grpc_cli/evans inspection

	log.Printf("[gRPC] listening on :%s", port)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("[gRPC] server error: %v", err)
		}
	}()
	return nil
}
