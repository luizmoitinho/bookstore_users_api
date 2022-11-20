FROM golang:1.18

#docker build -t bookstore_users_api  .
# docker run -d -p 8080:8080 -it bookstore_users_api 

WORKDIR /

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build -o bookstore_users_api

EXPOSE 8080

CMD [ "./bookstore_users_api"]