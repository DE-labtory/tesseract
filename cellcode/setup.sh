go build -buildmode=plugin -o /go/src/icode.so /icode/icode.go
go build -o /go/src/cellcode /cellcode/cellcode.go
/go/src/cellcode /go/src/icode.so
