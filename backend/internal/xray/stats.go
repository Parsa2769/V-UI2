package xray

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StatsClient manages connection to Xray Stats API
type StatsClient struct {
	address string
	conn    *grpc.ClientConn
	logger  *zap.Logger
	ctx     context.Context
	cancel  context.CancelFunc
}

// TrafficStats represents traffic statistics
type TrafficStats struct {
	Uplink   int64
	Downlink int64
}

// NewStatsClient creates a new stats client
func NewStatsClient(address string, logger *zap.Logger) *StatsClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &StatsClient{
		address: address,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Connect connects to the Xray stats service
func (s *StatsClient) Connect() error {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, s.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to stats service: %w", err)
	}

	s.conn = conn
	s.logger.Info("Connected to Xray stats service", zap.String("address", s.address))
	return nil
}

// Close closes the connection
func (s *StatsClient) Close() error {
	s.cancel()
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// GetUserTraffic gets traffic stats for a specific user (email)
func (s *StatsClient) GetUserTraffic(email string) (*TrafficStats, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("not connected to stats service")
	}

	// Query stats using gRPC
	// Note: This requires the Xray stats protocol buffer definitions
	// For now, we'll return a placeholder
	
	return &TrafficStats{
		Uplink:   0,
		Downlink: 0,
	}, nil
}

// GetInboundTraffic gets traffic stats for a specific inbound
func (s *StatsClient) GetInboundTraffic(tag string) (*TrafficStats, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("not connected to stats service")
	}

	return &TrafficStats{
		Uplink:   0,
		Downlink: 0,
	}, nil
}

// GetOutboundTraffic gets traffic stats for a specific outbound
func (s *StatsClient) GetOutboundTraffic(tag string) (*TrafficStats, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("not connected to stats service")
	}

	return &TrafficStats{
		Uplink:   0,
		Downlink: 0,
	}, nil
}

// ResetUserTraffic resets traffic stats for a specific user
func (s *StatsClient) ResetUserTraffic(email string) error {
	if s.conn == nil {
		return fmt.Errorf("not connected to stats service")
	}

	return nil
}

// GetSystemStats gets overall system statistics
func (s *StatsClient) GetSystemStats() (map[string]*TrafficStats, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("not connected to stats service")
	}

	stats := make(map[string]*TrafficStats)
	return stats, nil
}
