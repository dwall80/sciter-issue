# Build

If you need to recompile you will need to install go - golang for windows https://golang.org/dl/
Set up a GOPATH in env vars to point to a working directory and clone this repo into:
{GOPATH}/github.com/swiftdv8/sciter-issue

The only dependency is the go-sciter sdk so just run:

`go get github.com/sciter-sdk/go-sciter`

and drop the sciter.dll next to the binary in the repo

you should no be able to build with 

`go build`



# Usage

You can edit the template.html to switch between frame and div.

I've set up a simple custom asset handler that reads the files in from disk.
There is then a load loop which attempts to load the 1/2.html files alternately
into $(#content)

You can run the binary as follows to bypass the asset handler and use sciter
loading from file:

`./sciter --skip-loader`

What I observe is that if run with `--skip-loader` frame.load() and div.load()
both work as you have described

using the custom loader div.load() works as expected but frame.load() always
returns empty and never hits the custom asset handler
