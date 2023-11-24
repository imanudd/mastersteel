package api

import (
	"context"
	"finalproject/config"
	"finalproject/infra/db"
	"finalproject/internal/chatapp"
	customMiddleware "finalproject/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	cfg    *config.Config
	router *echo.Echo
	db     *pg.DB
}

func New() *API {
	cfg := config.LoadDefault()

	db := db.NewGoPG(cfg)

	router := echo.New()
	api := &API{
		cfg:    cfg,
		router: router,
		db:     db,
	}
	return api
}

func (api API) BuildHandler() *echo.Echo {
	api.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	api.router.Use(customMiddleware.RequestIDContext())
	// api.router.HTTPErrorHandler = CustomHTTPErrorHandler(api.cfg, api.log, httplog.NewHTTPLog(api.db))

	chatapp.RegisterAPI(
		*api.router.Group(""),
		api.cfg,
		chatapp.NewService(api.cfg, api.db,
			chatapp.NewRepository(api.db, api.cfg)),
	)

	return api.router
}

func (api API) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	handleSigterm(func() {
		cancel()
	})

	// Start server
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", api.cfg.Server.PORT),
		Handler: api.BuildHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	fmt.Printf("server is running at port: %v [env: %v]", api.cfg.Server.PORT, api.cfg.Server.ENV)

	gracefulShutdownServer(ctx, &server)
}

func gracefulShutdownServer(ctx context.Context, srv *http.Server) {

	<-ctx.Done()

	fmt.Println("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+s", err)
	}

	fmt.Println("server exited properly")

}

func handleSigterm(exitFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		exitFunc()
	}()
}
