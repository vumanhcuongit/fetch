IMAGE_NAME = fetch

fetch:
	@if ! docker image inspect $(IMAGE_NAME) > /dev/null 2>&1 ; then \
        docker build -t $(IMAGE_NAME) . ; \
    fi
	@docker run --rm \
		-v $(CURDIR):/app \
        -v $(CURDIR)/assets:/assets \
        -w /app \
        fetch \
        $(URLS) \
        $(if $(filter true,$(METADATA)),--metadata)

build:
	@docker build -t $(IMAGE_NAME) . ;