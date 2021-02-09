test:
	$(MAKE) integration_test

compose_file = ./.docker/docker-compose.yml

integration_test:
	docker-compose -f $(compose_file) up --build --abort-on-container-exit

ifeq ("$(TEST_WHAT)","")
TEST_WHAT=./...
endif

db-up:
	docker-compose -f $(compose_file) up -d postgres

destroy:
	docker-compose -f $(compose_file) down -v
