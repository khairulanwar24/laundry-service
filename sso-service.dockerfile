FROM alpine:latest

RUN mkdir /app

COPY ssoApp /app
COPY app/.env /app/.env
COPY app/.env /app/.init.sh

CMD [ "/app/ssoApp"]