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
