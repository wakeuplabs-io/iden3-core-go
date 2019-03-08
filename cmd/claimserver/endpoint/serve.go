package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/services/adminsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/rootsrv"

	log "github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server {
	api, serviceapi := genericserver.NewServiceAPI("/api/unstable")
	serviceapi.POST("/claims", handlePostClaim)                  // Get relay claim proof
	serviceapi.GET("/claims/:hi/proof", handleGetClaimProofByHi) // Get relay claim proof

	serviceapisrv := &http.Server{Addr: genericserver.C.Server.ServiceApi, Handler: api}
	go func() {
		log.Info("API server at ", genericserver.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server {
	api, adminapi := genericserver.NewAdminAPI("/api/unstable", stopch)
	adminapi.POST("/claims/basic", handleAddClaimBasic)

	adminapisrv := &http.Server{Addr: genericserver.C.Server.AdminApi, Handler: api}
	go func() {
		log.Info("ADMIN server at ", genericserver.C.Server.AdminApi)
		if err := adminapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(rs rootsrv.Service, cs claimsrv.Service, as adminsrv.Service) {

	genericserver.Claimservice = cs
	genericserver.Rootservice = rs
	genericserver.Adminservice = as

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	// start servers
	genericserver.Rootservice.Start()
	serviceapisrv := serveServiceApi()
	adminapisrv := serveAdminApi(stopch)

	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}

	if err := adminapisrv.Shutdown(context.Background()); err != nil {
		log.Error("AdminApi Shutdown:", err)
	} else {
		log.Info("AdminApi stopped")
	}

}