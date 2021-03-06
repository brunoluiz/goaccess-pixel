# goaccess pixel

Privacy focused pixel tracker server, based on Apache logs -- best used with Goaccess.io.

## Usage

### Using as server

Install it using one of the provided methods and use the following script on your page, replacing GOACCESS_PIXEL_ADDRESS with your goaccess-pixel address.

```js
  <script>
    const path = "u=" + encodeURIComponent(window.location.pathname);
    const referrer = document.referrer ? "&r=" + encodeURIComponent(document.referrer) : ""

    const _pixel = new Image(1, 1);
    _pixel.src = `http://GOACCESS_PIXEL_ADDRESS/?${path}${referrer}`;
  </script>
```

Once you run `goaccess-pixel`, it will output all the requests in log files, following Apache combined format. These can be used together with Goaccess using `goaccess ./access.log --log-format COMBINED`. By default, all logs are purged after a week (please refer to the parameters section for more information).

Have a look at `/example` for a suggested config using `docker-compose`

### Using as lib

If you want to run it on your own HTTP server, you will only need to use `handler.PixelLogger(io.Writer)`. On `cmd/goaccess-pixel/main.go` you can see an example of how to create a rotated log or you can use `os.Stdout`.

## Install

### MacOS

Use `brew` to install it

```
brew tap brunoluiz/tap
brew install goaccess-pixel
```

### Linux and Windows

[Check the releases section](https://github.com/brunoluiz/goaccess-pixel/releases) for more information details 

### go get

Install using `GO111MODULES=off go get github.com/brunoluiz/goaccess-pixel/cmd/goaccess-pixel` to get the latest version. This will place it in your `$GOPATH`, enabling it to be used anywhere in the system.

**⚠️ Reminder**: the command above download the contents of master, which might not be the stable version. [Check the releases](https://github.com/brunoluiz/goaccess-pixel/releases) and get a specific tag for stable versions.

### Docker

The tool is available as a Docker image as well. Please refer to [Docker Hub page](https://hub.docker.com/r/brunoluiz/goaccess-pixel/tags) to pick a release

```
docker run -p 80:80 \
  --env-file .env.sample \
  -v $(PWD)/access.log:/access.log \
  brunoluiz/goaccess-pixel
```

## Parameters

```
--port value               Server port (default: "80") [$PORT]
--log-file value           Log file output (default: "./access.log") [$LOG_FILE]
--log-max-age value        Log max age (default: 168h0m0s) [$LOG_MAX_AGE]
--log-rotation-time value  Time between each log rotation (default: 24h0m0s) [$LOG_ROTATION_TIME]
--pixel-route value        Pixel route (default: "/*") [$PIXEL_ROUTE]
--ready-route value        Ready probe route (default: "/__/ready") [$READY_ROUTE]
--metrics-route value      Metrics route (default: "/__/metrics") [$METRICS_ROUTE]
```

## To-do

- npm package with a simple to use script
