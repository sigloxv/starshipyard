## Starship Yard Framework - Web Application Framework                       ##
###############################################################################
## BUILD
###############################################################################
TARGET="starshipyard"
DEV_BIN_PATH="./bin"
GO_PATH="/home/$(USER)/go"
GO_BUILD_FLAGS=""    
GO_BINARY="go" 
INSTALL_BIN_PATH="/usr/local/bin"
###############################################################################
#PACKAGES="core database log parser"
PERMISSIONS="CAP_NET_BIND_SERVICE=+eip"
###############################################################################
# NOTE [GO_BUILD_FLAGS] Variable
# Obstuficate could be added here
# NOTE [GO_BINARY] Variable
# Perhaps have the ability to try it with TinyGo, or a Multiverse OS 
# implementation of SSA based size reduction, and elimination of all non-Linux 
# standard library code. 
# NOTE [GO_PATH] Variable 
# We used to do `export GOPATH=$(PWD)` to keep all the dependencies segregated
# and to make it function more like cargo or rubygems. This may still be
# preferable for some. 
#
###############################################################################
# 
# Multiverse OS: Starship Yard Web Framework 
#  
# More than inspired by Rails, essentially a feature-for-feature implementation
# of `ruby` Rails in Go. Assuming the Rails was setup with Rails-composer, and
# came with very popular, but not default, plugins. 
#
# The remaining functionality missing is the controllers and models; which have 
# not been completed because developers have been experimenting with a variety
# of different implementations. 
#
# Leading to two popular paths forward that developers are still not decided
# upon: 
#
# 	1) Implementing the Controllers and Models using Ruby, that is either 
# 	   parsed and translated to Go functions for speed, and portability. 
#
# 	   Or possibly, including the `ruby` binary and using a socket connection  
# 	   or libruby (similar to FFI but faster) based hand off similar to how Rack
# 	   works, for handling the full MVC component, returning HTML ready for
# 	   a Go based asset pipeline, that includes HTML whitelisting, minification,
# 	   and other common pipeline functionality. 
#
# 	2) A pure Go implementation that is functionality and organized exactly like
# 	   Rails. Including not just comparable libraries to ActiveSupport, and
# 	   ActiveRecord, and ActiveController;  but rebuild these libraires almost
# 	   in their entirity in Go, so that controller, Model and view functions
# 	   look nearly identical to their Ruby counterparts. 
#
# 	   Then when the Multiverse OS Ruby VM in Go is complete, re-implement these
# 	   components using Go Ruby, to have a nearly identical developer experience
# 	   to `ruby` Rails while still being able to compile, and all but the code
# 	   exposed to the web developer implemented in Go, leaving only the
# 	   web-developer code in Ruby, with the ability to use inline Go, or even
# 	   the Go alternatives as needed or desired. 
#
#
###############################################################################
# Helpers
DEV_BINARY="$(DEV_BIN_PATH)/$(TARGET)"
INSTALLED_BINARY="$(INSTALL_BIN_PATH)/$(TARGET)"
###############################################################################

# TODO: Prefix build with clean, and only clean if starship yard exists. 

.PHONY: all 
all: build

build:
		@echo ".===================================."
		@echo "|Building Starship Yard Framework   |"
		@echo "'==================================="
		@echo "Downloading starship yard dependencies..." 
		@cd ./framework && go get 
		# TODO: Will need to get all the submodules too if they are missing 
		@echo "Building the binary, output to `$(DEV_BINARY)`"
		@go build $(GO_BUILD_FLAGS) -o $(DEV_BINARY)
		@echo "Setting the binary permissions enabling it to use port 80 and 443"
		sudo setcap $(PERMISSIONS) $(DEV_BINARY)
		# TODO Should actually check if the build occured and it has the 
		# permissions and if it failed, cleanup. 
		@echo "Build complete"

clean:
		# TODO: Only preform if bin/starship-yard exists
		@echo ".===================================."
		@echo "|Installing Starship Yard Framework |"
		@echo "'===================================" 
		@echo "Cleaning up using `$(GO_BINARY) clean` command"
		@go clean
		@echo "Removing the binary at `$(DEV_BINARY)`"
		@rm -f $(INSTALLED_BINARY)
		# TODO: Remove any fines in binary path, perhaps delete it, and only 
		# create it when it is populated.
		#@rm "$(INSTALL_BIN_PATH)/*" 
		@echo "Cleanup successfull, all artificates removed."

install:
		@echo ".===================================."
		@echo "|Installing Starship Yard Framework |"
		@echo "'==================================="
		@echo "Copying binary into `$(INSTALLED_BINARY)`"
		@cp $(DEV_BINARY) $(INSTALL_BIN_PATH)

uninstall: 
		@echo ".====================================."
		@echo "|Uinstalling Starship Yard Framework |"
		@echo "'===================================="
		# TODO: Should instlal documentation, auto-complete, and other useful
		#       utilities. 
		@echo "Removing binary from \$PATH `$(INSTALLED_BINARY`"
		@cp $(INSTALLED_BINARY)
