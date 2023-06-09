# GoServe
CLI for starting blazing fast server for your local sites. Also perfect for designing Html, css
templates

## Installation
You can install goserve globally from npm
```bash
    npm i -g @nelwhix/goserve
```

To confirm it is downloaded, run
```bash
goserve -h
```

## Usage
For now, Serve supports just HTML, CSS and Javascript
```bash
    goserve # this starts a server on port 3000

    # To use a custom port run
    goserve -p 5173

    # To serve another folder
    goserve -root "/Desktop/my-cool-site/dist"

    # goserve -root "/Desktop/my-cool-site/dist/index.html" won't work because the root flag needs a directory to serve not a file
```

For local development, you will need to add this script tag to your html head

```html
    <script src="https://nelwhix-serve.s3.eu-central-1.amazonaws.com/serve.js"></script>
```
this downloads a script that listens for filechange events
from the server

### Troubleshooting
On Windows, if after install the goserve command is not available. Confirm that the nodejs bin directory is in your system's environment variables PATH