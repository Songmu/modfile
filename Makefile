.PHONY: place
place:
	go run _tools/place.go

.PHONY: clean
clean:
	rm -rf *.go internal/ testdata/
