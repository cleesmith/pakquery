package route

import (
	"net/http"

	"github.com/cleesmith/pakquery/controller"
	// "github.com/cleesmith/pakquery/route/middleware/acl"
	hr "github.com/cleesmith/pakquery/route/middleware/httprouterwrapper"
	// "github.com/cleesmith/pakquery/route/middleware/logrequest"
	// "github.com/cleesmith/pakquery/shared/session"

	"github.com/gorilla/context"
	// "github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Load the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// Load the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// Load the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

func routes() *httprouter.Router {
	r := httprouter.New()

	// Set 404 handler
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)

	// Serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))

	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.Index)))

	r.GET("/help", hr.Handler(alice.
		New().
		ThenFunc(controller.HelpGET)))

	r.GET("/records", hr.Handler(alice.
		New().
		ThenFunc(controller.RecordsGET)))

	r.POST("/records", hr.Handler(alice.
		New().
		ThenFunc(controller.RecordsPOST)))

	r.GET("/recordshow", hr.Handler(alice.
		New().
		ThenFunc(controller.RecordShowGET)))

	return r
}

func middleware(h http.Handler) http.Handler {
	// // Prevents CSRF and Double Submits
	// cs := csrfbanana.New(h, session.Store, session.Name)
	// cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	// cs.ClearAfterUsage(true)
	// cs.ExcludeRegexPaths([]string{"/static(.*)"})
	// csrfbanana.TokenLength = 32
	// csrfbanana.TokenName = "token"
	// csrfbanana.SingleToken = false
	// h = cs

	// // Log every request
	// h = logrequest.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
