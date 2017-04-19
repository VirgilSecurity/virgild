FROM alpine:3.5
MAINTAINER Virgil <support@VirgilSecurity.com>
RUN apk add --update ca-certificates
ARG GIT_COMMIT=unkown
ARG GIT_BRANCH=unkown
LABEL git-commit=$GIT_COMMIT
LABEL git-branch=$GIT_BRANCH
ADD virgild .

ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/virgild"]
