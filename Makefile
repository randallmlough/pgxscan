home				        = 	$(shell home)
software_version	  		=	$(shell cat VERSION)
version_array		    	=	$(subst ., ,$(software_version))
major				        = 	$(word 1,${version_array})
minor				        = 	$(word 2,${version_array})
patch				        = 	$(word 3,${version_array})
pwd 				        = 	$(shell pwd)

patch:
	- @echo "BUMPING PATCH"
	- @echo "Current Version: $(software_version)"
	- $(eval patch=$(shell echo $$(($(patch)+1))))
	- @echo "New Version: $(major).$(minor).$(patch)"
	- @printf $(major).$(minor).$(patch) > VERSION

minor:
	- @echo "BUMPING MINOR"
	- @echo "Current Version: $(software_version)"
	- $(eval minor=$(shell echo $$(($(minor)+1))))
	- @echo "New Version: $(major).$(minor).0"
	- @printf $(major).$(minor).0 > VERSION

major:
	- @echo "BUMPING MAJOR"
	- @echo "Current Version: $(software_version)"
	- $(eval major=$(shell echo $$(($(major)+1))))
	- @echo "New Version: $(major).0.0"
	- @printf $(major).0.0 > VERSION


.PHONY: patch minor major


# run tests using docker
pg_container	= pgxscan_postgres
network_name	= pgxscan_network
pg_user 		= root
pg_pass			= pass
pg_db			= pgxscan

net-up:
	docker network create $(network_name)

net-down:
	docker network rm $(network_name) || true

ifneq ("$(EXPOSE)","")
EXPOSE_PORTS=-p 5435:5432
endif

ifeq ("$(TEST_WHAT)","")
TEST_WHAT=./...
endif

db-up:
	docker run \
		--rm \
		-d \
		--network $(network_name) \
		--name $(pg_container) \
		$(EXPOSE_PORTS) \
		--env POSTGRES_USER=$(pg_user) \
		--env POSTGRES_PASSWORD=$(pg_pass) \
		--env POSTGRES_DB=$(pg_db) \
		postgres:12-alpine

db-down:
	docker stop $(pg_container) || true

test-up: | net-up db-up
test-down: | db-down net-down

test:
	$(MAKE) test-down
	$(MAKE) test-up
	$(MAKE) test-go
	$(MAKE) test-down

test-go:
	docker run \
		--rm \
		--network $(network_name) \
		--volume `pwd`:/test-go/ \
		--workdir /test-go \
		--env PG_URI="postgres://$(pg_user):$(pg_pass)@$(pg_container):5432/$(pg_db)?sslmode=disable" \
		-it golang:1.14-buster \
		go test $(TEST_WHAT)
