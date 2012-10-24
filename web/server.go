/*
	The web package contains all the code to provide Inbucket's web GUI
*/
package web

import (
	"fmt"
	"github.com/goods/httpbuf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jhillyerd/inbucket/config"
	"github.com/jhillyerd/inbucket/log"
	"net/http"
	"time"
)

type handler func(http.ResponseWriter, *http.Request, *Context) error

var Router *mux.Router

var sessionStore sessions.Store

func setupRoutes(cfg config.WebConfig) {
	log.Info("Theme templates mapped to '%v'", cfg.TemplateDir)
	log.Info("Theme static content mapped to '%v'", cfg.PublicDir)

	r := mux.NewRouter()
	// Static content
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/",
		http.FileServer(http.Dir(cfg.PublicDir))))

	// Root
	r.Path("/").Handler(handler(RootIndex)).Name("RootIndex").Methods("GET")
	r.Path("/status").Handler(handler(RootStatus)).Name("RootStatus").Methods("GET")
	r.Path("/mailbox").Handler(handler(MailboxIndex)).Name("MailboxIndex").Methods("GET")
	r.Path("/mailbox/list/{name}").Handler(handler(MailboxList)).Name("MailboxList").Methods("GET")
	r.Path("/mailbox/show/{name}/{id}").Handler(handler(MailboxShow)).Name("MailboxShow").Methods("GET")
	r.Path("/mailbox/html/{name}/{id}").Handler(handler(MailboxHtml)).Name("MailboxHtml").Methods("GET")
	r.Path("/mailbox/source/{name}/{id}").Handler(handler(MailboxSource)).Name("MailboxSource").Methods("GET")
	r.Path("/mailbox/delete/{name}/{id}").Handler(handler(MailboxDelete)).Name("MailboxDelete").Methods("POST")

	// Register w/ HTTP
	Router = r
	http.Handle("/", Router)
}

// Start() the web server
func Start() {
	cfg := config.GetWebConfig()
	setupRoutes(cfg)

	// TODO Make configurable
	sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

	addr := fmt.Sprintf("%v:%v", cfg.Ip4address, cfg.Ip4port)
	log.Info("HTTP listening on TCP4 %v", addr)
	s := &http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Error("HTTP server failed: %v", err)
	}
}

// ServeHTTP builds the context and passes onto the real handler
func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create the context
	ctx, err := NewContext(req)
	if err != nil {
		log.Error("Failed to create context: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ctx.Close()

	// Run the handler, grab the error, and report it
	buf := new(httpbuf.Buffer)
	log.Trace("Web: %v %v %v %v", req.RemoteAddr, req.Proto, req.Method, req.RequestURI)
	err = h(buf, req, ctx)
	if err != nil {
		log.Error("Error handling %v: %v", req.RequestURI, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the session
	if err = ctx.Session.Save(req, buf); err != nil {
		log.Error("Failed to save session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply the buffered response to the writer
	buf.Apply(w)
}
