package main

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/brunoluiz/goaccess-pixel/middleware"
	"github.com/brunoluiz/goaccess-pixel/tmpl"
	"github.com/go-chi/chi"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "Goaccess pixel route",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Usage:   "Server port",
				Value:   "80",
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:     "host",
				Usage:    "Application hostname (eg: https://pixel.example.com:3000/). Will be used on the script served on script-route",
				EnvVars:  []string{"HOST"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "log-file",
				Usage:   "Log file output",
				Value:   "./access.log",
				EnvVars: []string{"LOG_FILE"},
			},
			&cli.DurationFlag{
				Name:    "log-max-age",
				Usage:   "Log max age",
				Value:   7 * 24 * time.Hour,
				EnvVars: []string{"LOG_MAX_AGE"},
			},
			&cli.DurationFlag{
				Name:    "log-rotation-time",
				Usage:   "Time between each log rotation",
				Value:   time.Hour,
				EnvVars: []string{"LOG_ROTATION_TIME"},
			},
			&cli.StringFlag{
				Name:    "ready-route",
				Usage:   "Ready probe route",
				Value:   "/__/ready",
				EnvVars: []string{"READY_ROUTE"},
			},
			&cli.StringFlag{
				Name:    "metrics-route",
				Usage:   "Metrics route",
				Value:   "/__/metrics",
				EnvVars: []string{"METRICS_ROUTE"},
			},
			&cli.StringFlag{
				Name:    "script-route",
				Usage:   "Pixel script route",
				Value:   "/goaccess-pixel.js",
				EnvVars: []string{"METRICS_ROUTE"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Turn on debug mode",
				EnvVars: []string{"DEBUG"},
			},
		},
		Action: serve,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	if !c.Bool("debug") {
		logrus.SetLevel(logrus.FatalLevel)
	}

	r := chi.NewRouter()

	logf, err := rotatelogs.New(
		c.String("log-file")+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(c.String("log-file")),
		rotatelogs.WithMaxAge(c.Duration("log-max-age")),
		rotatelogs.WithRotationTime(time.Minute),
	)
	if err != nil {
		return err
	}
	r.
		With(middleware.Transform).
		With(middleware.Log(logf)).
		Get("/pixel.png", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	r.Get(c.String("script-route"), func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("script").Parse(tmpl.ScriptTmpl)
		if err != nil {
			logrus.Error(err)
		}

		if err := t.Execute(w, tmpl.ScriptTmplSubs{
			Hostname: c.String("host"),
		}); err != nil {
			logrus.Error(err)
		}
	})
	r.Get(c.String("metrics-route"), promhttp.Handler().ServeHTTP)
	r.Get(c.String("ready-route"), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return http.ListenAndServe(":"+c.String("port"), r)
}
