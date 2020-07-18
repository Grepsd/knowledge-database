FROM golang

WORKDIR /app
RUN go get -u github.com/cosmtrek/air
COPY build/air.conf /etc/air.conf
CMD air -d -c /etc/air.conf