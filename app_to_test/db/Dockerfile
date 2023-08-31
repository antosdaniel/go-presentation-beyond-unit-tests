FROM postgres:15.2-alpine
ENV VERSION v4.15.2

WORKDIR /db
COPY . .

RUN apk --no-cache add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/$VERSION/migrate.linux-amd64.tar.gz \
    | tar xvz

CMD ["./run-dev-migrations"]
