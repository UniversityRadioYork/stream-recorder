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

## Docker

-   Set configuration in `config.yml` (we're not doing config through Docker yet, too much work)
-   Build the image: `docker build -t stream-recorder .`
-   Run the container, exporting the port in `config.yml`, and also a volume for the recordings, going to `/usr/src/app/recordings` i.e.:

    `docker run -p 3001:3001 -v /mnt/recordings:/usr/src/app/recordings --name stream-recorder stream-recorder`

###### Michael Grace June 2021, University Radio York
