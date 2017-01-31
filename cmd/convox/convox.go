package main

import "github.com/convox/rack/cmd/convox/stdcli"

func init() {
	stdcli.Tags["app"] = stdcli.RenderAttributes(39)
	stdcli.Tags["config"] = stdcli.RenderAttributes()
	stdcli.Tags["description"] = stdcli.RenderAttributes(245)
	stdcli.Tags["fail"] = stdcli.RenderAttributes(160)
	stdcli.Tags["file"] = stdcli.RenderAttributes(249)
	stdcli.Tags["link"] = stdcli.RenderAttributes(4)
	stdcli.Tags["linenumber"] = stdcli.RenderAttributes(235)
	stdcli.Tags["release"] = stdcli.RenderAttributes(39)
	stdcli.Tags["security"] = stdcli.RenderAttributes(226)
	stdcli.Tags["service"] = stdcli.RenderAttributes(39)
	stdcli.Tags["database"] = stdcli.RenderAttributes(141)
	stdcli.Tags["success"] = stdcli.RenderAttributes(10)
	stdcli.Tags["unsupported"] = stdcli.RenderAttributes(220)
	stdcli.Tags["warning"] = stdcli.RenderAttributes(172)
}
