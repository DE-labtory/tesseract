go build -buildmode=plugin -o  ./icode.so ../../test/icode_test/icode.go
go run ../cellcode.go ./icode.so $1