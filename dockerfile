ARG BIN_NAME="upv-gym-reservations"
ARG RUNTIME_USER="user"
ARG WORKDIR="/app"

FROM golang:alpine as build
# https://github.com/moby/moby/issues/37622#issuecomment-412101935
ARG BIN_NAME
ARG WORKDIR

WORKDIR ${WORKDIR}
COPY . .
RUN go build -o ${BIN_NAME} .

FROM alpine

ARG BIN_NAME
ARG RUNTIME_USER
ARG WORKDIR

RUN adduser -h ${WORKDIR} -s /bin/sh -D ${RUNTIME_USER}

USER ${RUNTIME_USER}
WORKDIR ${WORKDIR}

# Copy executable binary
COPY --from=build --chown=${RUNTIME_USER}:${RUNTIME_USER} ${WORKDIR}/${BIN_NAME} .

# Include exmaple config file
COPY --chown=${RUNTIME_USER}:${RUNTIME_USER} config.example.json . 

ENV WORKDIR ${WORKDIR}
ENV BIN_NAME ${BIN_NAME}
ENTRYPOINT ${WORKDIR}/${BIN_NAME}