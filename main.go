package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/SENERGY-Platform/gin-middleware"
	"github.com/SENERGY-Platform/go-service-base/sql-db-hdl"
	"github.com/SENERGY-Platform/go-service-base/srv-info-hdl"
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/SENERGY-Platform/go-service-base/watchdog"
	"github.com/SENERGY-Platform/mgw-device-manager/api"
	"github.com/SENERGY-Platform/mgw-device-manager/handler/devices_hdl"
	"github.com/SENERGY-Platform/mgw-device-manager/handler/http_hdl"
	"github.com/SENERGY-Platform/mgw-device-manager/handler/storage_hdl"
	lib_model "github.com/SENERGY-Platform/mgw-device-manager/lib/model"
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"github.com/SENERGY-Platform/mgw-device-manager/util/db"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

var version string

func main() {
	srvInfoHdl := srv_info_hdl.New("device-manager", version)

	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	util.ParseFlags()

	config, err := util.NewConfig(util.Flags.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		ec = 1
		return
	}

	logFile, err := util.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *sb_util.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	util.Logger.Printf("%s %s", srvInfoHdl.GetName(), srvInfoHdl.GetVersion())

	util.Logger.Debugf("config: %s", sb_util.ToJsonStr(config))

	watchdog.Logger = util.Logger
	wtchdg := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	sql_db_hdl.Logger = util.Logger
	db, err := db.New(config.Database.Path)
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	defer db.Close()

	deviceHdl := devices_hdl.New(storage_hdl.New(db), time.Duration(config.Database.Timeout))

	mApi := api.New(deviceHdl)

	gin.SetMode(gin.ReleaseMode)
	httpHandler := gin.New()
	staticHeader := map[string]string{
		lib_model.HeaderApiVer:  srvInfoHdl.GetVersion(),
		lib_model.HeaderSrvName: srvInfoHdl.GetName(),
	}
	httpHandler.Use(gin_mw.StaticHeaderHandler(staticHeader), requestid.New(requestid.WithCustomHeaderStrKey(lib_model.HeaderRequestID)), gin_mw.LoggerHandler(util.Logger, http_hdl.GetPathFilter(), func(gc *gin.Context) string {
		return requestid.Get(gc)
	}), gin_mw.ErrorHandler(util.GetStatusCode, ", "), gin.Recovery())
	httpHandler.UseRawPath = true

	http_hdl.SetRoutes(httpHandler, mApi)
	util.Logger.Debugf("routes: %s", sb_util.ToJsonStr(http_hdl.GetRoutes(httpHandler)))

	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(config.ServerPort), 10))
	if err != nil {
		util.Logger.Error(err)
		ec = 1
		return
	}
	server := &http.Server{Handler: httpHandler}
	srvCtx, srvCF := context.WithCancel(context.Background())
	wtchdg.RegisterStopFunc(func() error {
		if srvCtx.Err() == nil {
			ctxWt, cf := context.WithTimeout(context.Background(), time.Second*5)
			defer cf()
			if err := server.Shutdown(ctxWt); err != nil {
				return err
			}
			util.Logger.Info("http server shutdown complete")
		}
		return nil
	})
	wtchdg.RegisterHealthFunc(func() bool {
		if srvCtx.Err() == nil {
			return true
		}
		util.Logger.Error("http server closed unexpectedly")
		return false
	})

	wtchdg.Start()

	dbCtx, dbCF := context.WithCancel(context.Background())
	wtchdg.RegisterStopFunc(func() error {
		dbCF()
		return nil
	})

	go func() {
		defer dbCF()
		if err = sql_db_hdl.InitDB(dbCtx, db, config.Database.SchemaPath, time.Second*5, time.Duration(config.Database.Timeout)); err != nil {
			util.Logger.Error(err)
			ec = 1
			wtchdg.Trigger()
			return
		}
	}()

	go func() {
		defer srvCF()
		util.Logger.Info("starting http server ...")
		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			util.Logger.Error(err)
			ec = 1
			return
		}
	}()

	ec = wtchdg.Join()
}
