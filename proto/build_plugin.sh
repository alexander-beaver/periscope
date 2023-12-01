#!/bin/zsh
protoc --go_out=.. --go-grpc_out=.. ./plugin.proto
grpc_tools_node_protoc --js_out=import_style=commonjs,binary:../common/node/static_codegen/ --grpc_out=../common/node/static_codegen --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` plugin.proto
python3 -m grpc_tools.protoc -I . --python_out=../common/py --pyi_out=../common/py --grpc_python_out=../common/py ./plugin.proto