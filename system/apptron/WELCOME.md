# Welcome to Apptron

Thanks for using Apptron, a full blown Linux based developer environment that runs entirely in your web browser.

There is a lot you can do with this environment. If this is your first time using Apptron, this document shows a few features for you to try.

If you close this document, you can bring it back by running this in the terminal:

```
open /apptron/WELCOME.md
```

## Install a package

Apptron is based on Alpine Linux, so you can install any package in the [Alpine package repository](https://pkgs.alpinelinux.org/packages). Try running this command:

```
apk add -u sl
```

This will install the latest in reverse "ls" technology. Once the package is installed try running it!

```
sl
```

Fun for the whole family!

## Publish an HTML file

Among other things, Apptron might just be the quickest way to publish websites:

- In the sidebar, click the file icon with the plus symbol for "New File"
- Name the new file `index.html`
- Copy the code below into `index.html`

```html
<html>
    <head><title>Apptron Demo</title></head>
    <body>
        <h1>Hello world.</h1>
        <pre>I am an HTML file.</pre>
    </body>
</html>
```

- Save the file using Ctrl+S, or Command+S on macOS
- In the upper right, click the "Share" button
- Select the "Publish" tab
- Leave the *Source Path* as `.` for the project root
- Copy the *Public URL* into your clipboard
- Click the "Publish" button

To view the page you just published:

- Open a new tab
- Paste the URL from your clipboard into the tab location bar
- Press the "Return" key to visit that page

Congrats! You successfully published a web page with Apptron.

> NOTE: If you get a 404, it might be because you loaded the URL too soon and now have a negative cache entry for that page. Clear that with a hard reload, holding Shift while reloading.

## Advanced Publishing

The "Share > Publish" UI is a convenience workflow for copying files to the system's `/public` mount. Any files there will be published to the URL that you used above.

Try it out by doing the following in your Terminal:

```
echo "Hello from the command-line" > /public/hello.txt
```

Now add `/hello.txt` to the end of the URL from earlier and load it, you should see the file you just created.

To sync a whole directory to `/public` from the command-line, you can use the `publish` command. Running this in the project root is the same as using the Publish UI using `.` as the source:

```
publish .
```

## Write and build a Go program

While writing and running JavaScript code in the browser is quite common, Apptron is the first environment that lets you write and compile the Go systems language entirely in-browser. 

Start with loading the Go profile by running this command in the terminal:

```
source /etc/goprofile 
```

Now you should have the `go` toolchain. You can see the version with:

```
go version
```

Use the "New File" button to create a file named `hello.go` and then paste this code into it:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hello <name>")
		os.Exit(1)
	}
	fmt.Printf("Hello, %s!\n", os.Args[1])
}
```

Then build the file we just created by running:

```
time go build -o hello hello.go
```

Patiently wait while the Go compiler runs in your browser.

> On an M3 Macbook, this takes about 10 seconds to run the first time. Let us know how long it takes on your system: Copy and paste the results into the Feedback dialog, which you can get by clicking the icon to right of the GitHub and Discord icons in the top bar.

This will create a binary named `hello` that accepts a name as a parameter, run it by typing this command:

```
./hello Apptron
```

You should see this:

```
Hello, Apptron!
```

Go lets you cross compile to any platform/architecture. See if you can figure out how to compile a version of the program that will run natively on your platform if you download it.

Go is the first of many languages Apptron will support. If you have a preference for what we support next, cast your vote in this [GitHub Discussion thread](https://github.com/tractordev/apptron/discussions/215).

## Run Apache from inside your browser

We start by installing Apache:

```
apk add -u apache2
```

Once that is done, run Apache directly in the foreground by running this command:

```
httpd -DFOREGROUND
```

After Apache starts with some non-critical warnings, Apptron will detect it listening on port 80 and give you a public URL to visit. 

It should look like this:

```
AH00557: httpd: apr_sockaddr_info_get() failed for (none)
AH00558: httpd: Could not reliably determine the server's fully qualified domain name, using 127.0.0.1. Set the 'ServerName' directive globally to suppress this message

=> Apptron public URL: https://tcp-80-0a00001e-example.apptron.dev

```

Once you see the Apptron public URL, click or copy-paste it to see Apache running!

This URL is sharable and will work for anyone on the internet as long as your browser tab is running.

What other servers can you run in the browser with Apptron?

## More tips for this environment

* There are 3 persistent mounts in an Apptron environment:
  * `/project` - Files shown in the sidebar for this environment
  * `/home/$USER` - Your home directory available in all environments
  * `/public` - The website for this environment
* All other files, including installed programs, will be reset with every pageload / session.
* An `.apptron` directory in the project root lets you customize the environment with these files:
  * `.apptron/envrc` - Commands here will be run at the start of every session
  * `.apptron/envbuild` - Commands here will be used to rebuild the environment. (experimental)
* Files can be uploaded into the project mount by dragging them into the sidebar
* You can use the `open` command to open more files or folders into the editor UI

## Give us feedback

You can easily send us a message by clicking the Feedback icon to the right of the GitHub and Discord icons next to the Apptron logo in the top bar. We'd love to hear from you!

Even better if you [join our Discord](https://discord.gg/zCrpdAgZAf) or [star the project on GitHub](https://github.com/tractordev/apptron). :)

Enjoy!<br/>
—Apptron team