# URY Stream Recorder

Records various icecast endpoints when they go live.

## Running

-   Set configuration in `config.yml`
-   Build the frontend:
    -   `cd frontend`
    -   `npm install`
    -   `npm run build`
-   Build the server:
    -   `go build`
    -   `./stream-recorder`
-   View the client to download recordings at the port specified in `config.yml`, default `3001`

The audio files are stored in `/recordings` and their data in `recordings.yml` (will be created if it doesn't exist).

###### Michael Grace June 2021, University Radio York
