FROM alpine:3.4
MAINTAINER Virgil <support@VirgilSecurity.com>
RUN apk add --update ca-certificates
ARG GIT_COMMIT=unkown
ARG GIT_BRANCH=unkown
LABEL git-commit=$GIT_COMMIT
LABEL git-branch=$GIT_BRANCH
ADD build/linux_x64/virgild .
VOLUME ["/srv"]

ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/virgild"]
CMD ["--config=/srv/virgild.conf"]
