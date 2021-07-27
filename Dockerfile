FROM node:latest AS webBuild
WORKDIR /app

COPY frontend/package*.json ./
RUN npm install
RUN npm run build

FROM golang:1.16
WORKDIR /go/src/app
COPY . .

COPY --from=webBuild /app/frontend/build frontend/build

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3001

CMD [ "app" ]


