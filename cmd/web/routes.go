package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	/*Pat matches patterns in the order that they are registered. In our application, 
	a HTTP request to GET "/snippet/create" is actually a valid match for two routes — it’s 
	an exact match for /snippet/create, and a wildcard match for /snippet/:id (the "create" 
	part of the path would be treated as the :id parameter). So to ensure that the exact match 
	takes preference, we need to register the exact match routes before any wildcard routes.*/
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/mySnippets", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showSnippetByUser))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Post("/snippet/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteSnippet))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
	mux.Post("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	fileServer := http.FileServer(http.Dir("M:/Projects/snippetbox/ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	
	return standardMiddleware.Then(mux)
}
