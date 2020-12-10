# Copyright 2020 Apulic, Inc.  All rights reserved.
# 
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
# 
#      http://www.apache.org/licenses/LICENSE-2.0
# 
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
FROM golang:1.13.7-alpine3.11

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add git pkgconfig build-base && mkdir -p /go/src/apulis/AIArtsBackend
WORKDIR /go/src/apulis/AIArtsBackend
ADD ./go.mod .
ADD ./go.sum .
#RUN go mod download && go get -u github.com/swaggo/swag/cmd/swag
ADD . .
#RUN swag init && GO111MODULE=${GO111MODULE} go build -o /go/bin/AIArtsBackend
RUN go mod download && GO111MODULE=${GO111MODULE} go build -o /go/bin/AIArtsBackend

FROM alpine:3.11
RUN apk --no-cache add ca-certificates libdrm
WORKDIR /root/
COPY --from=0 /go/bin/AIArtsBackend .
CMD ["./AIArtsBackend"]



