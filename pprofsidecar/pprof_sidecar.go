package pprofsidecar

import (
	"context"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

// Run starts an errgroup.Group that runs a http server for pprof along with whatever
// function/service is supposed to run
//
// This should not be used for http servers
func Run(ctx context.Context, srvAddr string, interrupt chan os.Signal, run func(context.Context) error) error {
	if interrupt == nil {
		interrupt = make(chan os.Signal, 3)
		defer close(interrupt)
		signal.Notify(interrupt, os.Interrupt)
	}

	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	srv := &http.Server{
		Addr:              srvAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	} // the pprof debug server

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// watches for interrupts
	g.Go(func() error {
		select {
		case <-interrupt:
			cancel()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// runs the function
	g.Go(func() error {
		defer cancel()
		return run(ctx)
	})

	// runs the pprof server
	g.Go(srv.ListenAndServe)

	// kills the pprof server when necessary
	g.Go(func() error {
		<-ctx.Done() // something said we are done

		shutdownCtx, cncl := context.WithTimeout(context.Background(), 2*time.Second)
		defer cncl()

		return srv.Shutdown(shutdownCtx) //nolint:contextcheck // want to use a fresh one here
	})

	return g.Wait()
}
