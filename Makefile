#!make
# include ./packages/api/.env

ifneq (,$(wildcard ./packages/api/.env))
	include ./packages/api/.env
	export
endif

dev:
	@( \
		trap 'kill 0' INT; \
		cd $(CURDIR)/packages/api/cmd && go run main.go & \
		cd $(CURDIR)/packages/www && go run main.go & \
		cd $(CURDIR)/packages/www && npx tailwindcss -i ./src/css/style.css -o ./src/css/output.css --watch \
		wait \
		)
