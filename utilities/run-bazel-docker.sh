#!/usr/bin/env sh

######################################################################
# @author      : hung0913208 (hung0913208@gmail.com)
# @file        : run-bazel-docker
# @created     : Wednesday Mar 10, 2021 21:03:19 +07
#
# @description : Run bazel command using docker
######################################################################

if which docker &> /dev/null; then
	CMD=$1

	if [ $CMD = 'test' ]; then
		PARAMS='--test-tag-filters=selftest'
		CMD='test'
	elif [ $CMD = 'bench' ]; then
		PARAMS='--test-tag-filters=benchtest'
		CMD='test'
	fi

	shift
	docker run 				\
		-e USER="$(id -u)" 		\
	    	-u="$(id -u)" 			\
		-v $(pwd):$(pwd) 		\
		-v $(pwd):$(pwd) 		\
		-w /src/workspace 		\
		l.gcr.io/google/bazel:latest 	\
		--output_user_root=$(pwd) 	\
		$CMD $PARAMS $@
else
	error "please install docker first"
fi
