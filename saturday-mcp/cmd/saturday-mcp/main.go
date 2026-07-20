package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/adapters/otel"
	"github.com/orieken/saturday-mcp/internal/domain"
	"github.com/orieken/saturday-mcp/internal/logging"
	saturdayServer "github.com/orieken/saturday-mcp/internal/server"
)

const (
	serverName    = "saturday-mcp"
	serverVersion = "0.1.0"

	// otelEndpointEnv is the OTel-conventional env var for the OTLP
	// endpoint (host:port for gRPC). Presence toggles the real tracer
	// on; absence keeps the noop tracer.
	otelEndpointEnv = "OTEL_EXPORTER_OTLP_ENDPOINT"

	// otelInsecureEnv gates plaintext gRPC transport (local collectors,
	// dev clusters). Any truthy value ("true"/"1") enables it.
	otelInsecureEnv = "OTEL_EXPORTER_OTLP_INSECURE"

	// shutdownTimeout bounds the OTel provider's flush-and-close at
	// process exit; long enough to drain a queued batch, short enough
	// that a stuck exporter doesn't hold up shutdown.
	shutdownTimeout = 5 * time.Second
)

func main() {
	logger := logging.NewLogger(os.Stderr)
	logger.Info("Starting Saturday MCP Server", "version", serverVersion)

	tracer, shutdown := buildTracer(context.Background(), logger)
	defer flushTracer(shutdown, logger)

	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithLogging(),
	)

	saturdayHandler, err := saturdayServer.NewHandler(logger, saturdayServer.WithTracer(tracer))
	if err != nil {
		logger.Error("Failed to initialize server handler", "error", err)
		os.Exit(1)
	}

	if err := saturdayHandler.RegisterTools(mcpServer); err != nil {
		logger.Error("Failed to register tools", "error", err)
		os.Exit(1)
	}

	if err := saturdayHandler.RegisterResources(mcpServer); err != nil {
		logger.Error("Failed to register resources", "error", err)
		os.Exit(1)
	}

	if err := saturdayHandler.RegisterPrompts(mcpServer); err != nil {
		logger.Error("Failed to register prompts", "error", err)
		os.Exit(1)
	}

	logger.Info("Saturday MCP Server initialized successfully")

	if err := server.ServeStdio(mcpServer); err != nil {
		logger.Error("Server error", "error", err)
		os.Exit(1)
	}
}

// buildTracer chooses between the real OTel-gRPC tracer and the noop
// based on OTEL_EXPORTER_OTLP_ENDPOINT. Returns (tracer, shutdown) —
// the shutdown callback is nil for the noop path, non-nil for the real
// path so main can defer a bounded flush at process exit.
//
// If endpoint construction fails at startup (e.g. malformed URL), the
// error is logged and the process falls back to noop rather than
// exiting — an observability failure should never take down the MCP
// server, only degrade what it exports.
func buildTracer(ctx context.Context, logger *logging.Logger) (domain.Tracer, func(context.Context) error) {
	endpoint := os.Getenv(otelEndpointEnv)
	if endpoint == "" {
		logger.Info("OTel tracer disabled (no endpoint configured)", "env", otelEndpointEnv)
		return otel.NewNoopTracer(), nil
	}

	cfg := otel.Config{
		Endpoint:    endpoint,
		ServiceName: serverName,
		Insecure:    isTruthy(os.Getenv(otelInsecureEnv)),
	}
	tracer, err := otel.NewOTelTracer(ctx, cfg)
	if err != nil {
		logger.Error("Failed to build OTel tracer, falling back to noop", "error", err, "endpoint", endpoint)
		return otel.NewNoopTracer(), nil
	}
	logger.Info("OTel tracer enabled", "endpoint", endpoint, "insecure", cfg.Insecure)
	return tracer, tracer.Shutdown
}

// flushTracer runs the OTel provider's shutdown under a bounded context.
// A no-op when shutdown is nil (the noop path). Errors are logged but
// swallowed — shutdown races with process exit and a partial flush is
// preferable to a stalled exit.
func flushTracer(shutdown func(context.Context) error, logger *logging.Logger) {
	if shutdown == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		logger.Error("OTel tracer shutdown reported an error", "error", err)
	}
}

// isTruthy accepts the common opt-in strings without pulling in a bigger
// parser. Anything else — including empty, "false", "0", "no" — reads
// as false, matching how most Go OTel examples treat this flag.
func isTruthy(v string) bool {
	switch strings.ToLower(v) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
