package tmpl

// ScriptParams params to render ScriptTmpl
type ScriptTmplSubs struct {
	Hostname string
}

// ScriptTmpl script to be used to call the pixel route
const ScriptTmpl = `(function () {
	const path = "u=" + encodeURIComponent(window.location.pathname);
	const referrer = document.referrer ? "&r=" + encodeURIComponent(document.referrer) : ""
	const time = "&t=" + (new Date()).getTime()

	const _pixel = new Image(1, 1);
	_pixel.src = '://{{ .Hostname }}' + '/pixel.png?' + path + referrer + time;
	console.log(_pixel.src)
})()`
