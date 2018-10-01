# Setup name variables for the package/tool
NAME := assume
PKG := github.com/s3than/$(NAME)

CGO_ENABLED := 0

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild: