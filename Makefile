VERSION := $(shell go run main.go version | sed 's/version: /v/')

collector: test
	go build ./.

test: setup
	go test ./...

setup:
	go get ./...

clean-zip:
	find pkg -name '*.zip' | xargs rm

all: setup test
	gox \
	    -os="linux" \
	    -arch="amd64" \
	    -output "pkg/{{.Dir}}_$(VERSION)-{{.OS}}-{{.Arch}}" \
	    .

compress: all clean-zip
	cd pkg && ( find . -perm -u+x -type f -name 'collector*' | gxargs -i zip -m {}.zip {} )

release: compress
	git push origin master
	rm -rf pkg/.DS_Store  
	ghr $(VERSION) pkg
	git fetch origin --tags

