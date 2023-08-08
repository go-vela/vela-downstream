# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

#############################################################################
##    docker build --no-cache --target certs -t vela-downstream:certs .    ##
#############################################################################

FROM alpine@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as certs

RUN apk add --update --no-cache ca-certificates

##############################################################
##    docker build --no-cache -t vela-downstream:local .    ##
##############################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-downstream /bin/vela-downstream

ENTRYPOINT [ "/bin/vela-downstream" ]
