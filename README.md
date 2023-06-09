# Serve
CLI for starting blazing fast server for your local sites. Perfect for designing Html, css
templates

## Usage
For now, Serve supports just HTML, CSS and Javascript
```bash
    serve // this starts a server on port 3000

    // To use a custom port run
    serve -p 5173

    // To serve another folder
    serve -root "path-to-folder"
```

For local development, you will need to add this script tag to your html head

```html
    <script src="https://nelwhix-serve.s3.eu-central-1.amazonaws.com/serve.js"></script>
```
this downloads a script that listens for filechange events
from the server

