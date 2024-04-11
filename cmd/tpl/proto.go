package tpl

func ProtoMakefile() string {
	return `.PHONY: proto
proto: clean
	buf generate

.PHONY: clean
clean:
	@rm -rf ../internal/proto && rm -rf ../internal/http/swagger	
`
}

func ProtoBufYaml() string {
	return `version: v1
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
`
}

func ProtoBufGenYaml() string {
	return `version: v1
plugins:
  - name: go
    out: ../internal/proto
    opt: paths=source_relative

  - name: go-grpc
    out: ../internal/proto
    opt: paths=source_relative

  - name: grpc-gateway
    out: ../internal/proto
    opt:
      - paths=source_relative
      - warn_on_unbound_methods=true

  - name: openapiv2
    out: ../internal/http/swagger
    opt:
      - proto3_optional_nullable=true
`
}

func PublicProto() string {
	return `syntax = "proto3";
package public;
option go_package = "./;public";

import "google/api/annotations.proto";
import "google/api/client.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "protoc-gen-openapiv2/options/openapiv2.proto";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
	host : "api.ruguoapp.com" // define your api host
	base_path : "/"
	schemes : HTTPS
	info : {
		title : "Public API"
		description : "{{ .Description }}"
		version : "1.0"
	}
};

service PublicService {
	rpc Echo(EchoRequest)
			returns (EchoResponse) {
		option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
			summary : "echo your request"
			tags : "接口示例"
		};
		option (google.api.http) = {
			post : "/1.0/go-coco/echo" // define your api route
			body : "*"
		};
	}
}

message EchoRequest {
	string data = 1
			[ (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
				description : "输入参数"
			} ];
}

message EchoResponse {
	string data = 1
			[ (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
				description : "输出参数"
			} ];
}
`
}
