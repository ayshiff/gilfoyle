all: $(SOURCES)

.PHONY: help
help:
	@echo "Usage: make RULES"
	@echo ""
	@echo "Main rules"
	@echo "  lauch-servers - Launch the s3 and gcs test servers"
	@echo "  lauch-s3-test-server - Launch the s3 test server"
	@echo "  lauch-gcs-test-server - Launch the gcs test server"

	@echo "  clean - Stop, delete containers and delete files"
	@echo "  clean-storage - Delete uploaded files"
	@echo "  clean-docker - Stop and delete containers"
	@echo "  stop - Stop the containers"

.PHONY: clean-docker
clean-docker:
	docker-compose down

.PHONY: clean-storage
clean-storage:
		find volume/gcs/storage/ ! \( -name .gitignore \) -type f -delete
		find volume/minio/ ! \( -name .gitignore \) -type f -delete
		find volume/gcs/storage/ -type d -depth 1 -delete

.PHONY:	clean
clean:	clean-docker clean-storage

.PHONY:	lauch-servers
lauch-servers:	
		docker-compose up

.PHONY:	lauch-s3-test-server
lauch-s3-test-server:	
		docker-compose run minio

.PHONY: lauch-gcs-test-server
lauch-gcs-test-server:	
		docker-compose run gcs