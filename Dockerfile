# SPDX-License-Identifier: Apache-2.0

#############################################################################
##    docker build --no-cache --target certs -t vela-downstream:certs .    ##
#############################################################################

FROM alpine@sha256:eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978 as certs

RUN apk add --update --no-cache ca-certificates

##############################################################
##    docker build --no-cache -t vela-downstream:local .    ##
##############################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-downstream /bin/vela-downstream

ENTRYPOINT [ "/bin/vela-downstream" ]
