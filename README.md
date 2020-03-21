#Bidding system.

#Dependency injection:
go mod download

#Steps to run
go run main.go

#steps to run benchmark tests
go test -bench=Main -cpuprofile=cpu.out  

