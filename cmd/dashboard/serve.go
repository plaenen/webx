package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/cmd/dashboard/internal"
	"github.com/spf13/cobra"
)

func serveCmd() *cobra.Command {
	var (
		port int
		pro  bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the showcase HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(port, pro)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3001, "port to listen on (0 for random)")
	cmd.Flags().BoolVar(&pro, "pro", false, "use Datastar Pro (requires BYOL files in byol/datastar/)")

	return cmd
}

func serve(port int, pro bool) error {
	r := chi.NewRouter()
	internal.SetupRoutes(r,
		internal.WithPro(pro),
	)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on port %d: %w", port, err)
	}

	dsMode := "open-source"
	if pro {
		dsMode = "pro"
	}
	slog.Info("server started", "address", fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port), "datastar", dsMode)

	return http.Serve(ln, r)
}
