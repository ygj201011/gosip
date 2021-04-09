package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ygj201011/gosip"
	"github.com/ygj201011/gosip/log"
	"github.com/ygj201011/gosip/sip"
)

var (
	logger log.Logger
)

func init() {
	logger = log.NewDefaultLogrusLogger().WithPrefix("Server")
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	srvConf := gosip.ServerConfig{
		UserAgent: "gaizi-server",
	}
	srv := gosip.NewServer(srvConf, nil, nil, logger)
	srv.OnRequest(sip.REGISTER, func(req sip.Request, tx sip.ServerTransaction) {
		//check information
		ah := req.GetHeaders("Authorization")
		if len(ah) == 0 {
			wwwAuthenheader := make([]sip.Header, 0)
			wwwAuthenheader = append(wwwAuthenheader, &sip.GenericHeader{HeaderName: "WWW-Authenticate", Contents: fmt.Sprintf("Digest realm=\"sip_reg\",nonce=\"%s\",stale=true,algorithm=MD5", sip.GenerateNonce())})
			srv.RespondOnRequest(req, 401, "Unauthorized", "", wwwAuthenheader)
		} else {
			var (
				auth = sip.AuthFromValue(ah[0].Value())
			)
			//remote response
			rr := auth.GetResponse()
			auth.SetUsername("1000")
			auth.SetPassword("1234")
			auth.SetMethod("REGISTER")
			auth.CalcResponse()
			//local response
			cr := auth.GetResponse()
			if rr != cr {
				srv.RespondOnRequest(req, 403, "OK", "", nil)
			} else {
				srv.RespondOnRequest(req, 200, "OK", "", nil)
			}
			logger.Errorf("%+v", req.Headers)
		}
	})
	srv.Listen("udp", "0.0.0.0:5091")

	<-stop

	srv.Shutdown()
}
