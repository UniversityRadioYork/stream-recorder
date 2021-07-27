FROM node:latest AS webBuild
WORKDIR /app

COPY frontend/package*.json ./
RUN npm install

COPY frontend .

RUN npm run build


FROM golang:1.16
WORKDIR /usr/src/app
COPY . .

COPY --from=webBuild /app/build frontend/build

RUN go build

EXPOSE 3001

CMD [ "./stream-recorder" ]