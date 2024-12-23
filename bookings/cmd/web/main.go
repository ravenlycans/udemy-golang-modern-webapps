package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/driver"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/handlers"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/helpers"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/models"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/render"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/routes"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const portNumber = 8080

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", portNumber), Handler: routes.Run()}

	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.ErrorLog.Fatalf("startHttpServer: %s\n", err.Error())
		}
	}()

	// Returning reference so caller can call Shutdown()
	return srv
}

// main is the application entrypoint
func main() {
	// Lets load out .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Cannot load .env file, please make sure it is in the root directory.")
	}

	// Register the complex types we want to store in the sessions.
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	// TODO: Change this to true when in production.
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Let's connect to the db and construct the dsn string.
	app.InfoLog.Printf("Connecting to database at %s:%s\n", os.Getenv("DB_SRV_URL"), os.Getenv("DB_SRV_PORT"))
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s", os.Getenv("DB_SRV_USER"), os.Getenv("DB_SRV_PASS"),
		os.Getenv("DB_SRV_URL"), os.Getenv("DB_SRV_PORT"), os.Getenv("DB_NAME"))
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatalf("main: Cannot connect to database. Err: %s\n", err.Error())
	}
	defer db.SQL.Close()
	app.InfoLog.Printf("Connected to database at %s:%s and pinged alive!!\n", os.Getenv("DB_SRV_URL"), os.Getenv("DB_SRV_PORT"))

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateCache()
	if err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	app.TemplateCache = tc
	app.UseCache = false
	render.New(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.New(repo)

	// initialize the routes package
	routes.New(repo)

	// initialize the helpers package.
	helpers.New(&app)

	// Register the middlewares used.
	routes.AddMiddleware(middleware.Recoverer)
	routes.AddMiddleware(NoSurf)
	routes.AddMiddleware(SessionLoad)

	// Add our static route.
	err = routes.SetStaticDir("/static", "./static/")
	if err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	// Register the routes available.
	err = routes.RegisterRoute("/", handlers.Repo.Home, "GET")
	err = routes.RegisterRoute("/about", handlers.Repo.About, "GET")
	err = routes.RegisterRoute("/favicon.ico", handlers.Repo.FavIcon, "GET")
	err = routes.RegisterRoute("/rooms/generals-quarters", handlers.Repo.RoomsGenerals, "GET")
	err = routes.RegisterRoute("/rooms/majors-suite", handlers.Repo.RoomsMajors, "GET")
	err = routes.RegisterRoute("/make-reservation", handlers.Repo.MakeReservation, "GET")
	err = routes.RegisterRoute("/make-reservation-ep", handlers.Repo.MakeReservationEP, "POST")
	err = routes.RegisterRoute("/reservation-summary", handlers.Repo.ReservationSummary, "GET")
	err = routes.RegisterRoute("/search-availability", handlers.Repo.SearchAvailability, "GET")
	err = routes.RegisterRoute("/search-availability-ep", handlers.Repo.SearchAvailabilityEP, "POST")
	err = routes.RegisterRoute("/search-availability-ep-json", handlers.Repo.SearchAvailabilityEPJSON, "POST")
	err = routes.RegisterRoute("/contact", handlers.Repo.Contact, "GET")

	if err != nil {
		app.ErrorLog.Printf("main: %s", err.Error())
		os.Exit(1)
	}

	app.InfoLog.Printf("Server is listening on port %d\n", portNumber)

	bNotQuitting := true
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startHttpServer(httpServerExitDone)

	// bNotQuitting is always true here, but will be set to false in the loop.
	for bNotQuitting {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press Q to quit\n")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "Q" || text == "q" {
			app.InfoLog.Printf("Quitting server, bye bye!!\n")
			bNotQuitting = false
			_ = srv.Shutdown(nil)
			httpServerExitDone.Wait()
			break
		}

		time.Sleep(1 * time.Second)
	}
}
