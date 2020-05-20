protoc -I src/ --go_out=src/ src/simple/simple.proto

protoc -I src/ --go_out=src/ src/enumdemo/enums.proto

protoc -I src/ --go_out=src/ src/complexDemo/complex.proto