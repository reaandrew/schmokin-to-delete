.PHONY: proto
proto:
	 cd server && protoc --go_out=plugins=grpc:. *.proto
	 
