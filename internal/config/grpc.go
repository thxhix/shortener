package config

// GRPCConfig config fields for gRPC server.
type GRPCConfig struct {
	Address string `env:"GRPC_ADDRESS" envDefault:":3200"`
}
