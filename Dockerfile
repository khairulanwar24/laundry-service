FROM alpine:latest

RUN mkdir /app

COPY SSO-App /app
COPY app/.env /app/.env
COPY app/.env /app/.init.sh

CMD [ "/app/SSO-App"]
