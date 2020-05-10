package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brunoluiz/goaccess-pixel/handler"
	"github.com/go-chi/chi"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Usage: "Goaccess pixel route",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "port",
				Usage:  "Server port",
				Value:  "80",
				EnvVar: "PORT",
			},
			&cli.StringFlag{
				Name:   "log-file",
				Usage:  "Log file output",
				Value:  "./access.log",
				EnvVar: "LOG_FILE",
			},
			&cli.DurationFlag{
				Name:   "log-max-age",
				Usage:  "Log max age",
				Value:  7 * 24 * time.Hour,
				EnvVar: "LOG_MAX_AGE",
			},
			&cli.DurationFlag{
				Name:   "log-rotation-time",
				Usage:  "Time between each log rotation",
				Value:  24 * time.Hour,
				EnvVar: "LOG_ROTATION_TIME",
			},
			&cli.StringFlag{
				Name:   "pixel-route",
				Usage:  "Pixel route",
				Value:  "/*",
				EnvVar: "PIXEL_ROUTE",
			},
			&cli.StringFlag{
				Name:   "ready-route",
				Usage:  "Ready probe route",
				Value:  "/__/ready",
				EnvVar: "READY_ROUTE",
			},
			&cli.StringFlag{
				Name:   "metrics-route",
				Usage:  "Metrics route",
				Value:  "/__/metrics",
				EnvVar: "METRICS_ROUTE",
			},
		},
		Action: serve,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func serve(c *cli.Context) error {
	r := chi.NewRouter()

	output, err := getLogger(
		c.String("log-file"),
		rotatelogs.WithLinkName(c.String("log-file")),
		rotatelogs.WithMaxAge(c.Duration("log-max-age")),
		rotatelogs.WithRotationTime(c.Duration("log-rotation-time")),
	)
	if err != nil {
		return err
	}

	r.Get(c.String("pixel-route"), handler.PixelLogger(output).ServeHTTP)
	r.Get(c.String("metrics-route"), promhttp.Handler().ServeHTTP)
	r.Get(c.String("ready-route"), func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return http.ListenAndServe(":"+c.String("port"), r)
}

func getLogger(file string, opts ...rotatelogs.Option) (io.Writer, error) {
	if file == "" || file == "/dev/stdout" {
		return os.Stdout, nil
	}

	if file == "/dev/stderr" {
		return os.Stderr, nil
	}

	return rotatelogs.New(file+".%Y%m%d%H%M", opts...)
}
