default:
	@go build -o builds/gf -ldflags '-s -w'

clean:
	@rm builds -rf

releases:
	@python3 makereleases.py