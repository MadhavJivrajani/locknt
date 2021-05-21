sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=100000x -cpuprofile profile_file.out
go test -bench=. -benchtime=100000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=90000x -cpuprofile profile_file.out
go test -bench=. -benchtime=90000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=80000x -cpuprofile profile_file.out
go test -bench=. -benchtime=80000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=70000x -cpuprofile profile_file.out
go test -bench=. -benchtime=70000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=60000x -cpuprofile profile_file.out
go test -bench=. -benchtime=60000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=50000x -cpuprofile profile_file.out
go test -bench=. -benchtime=50000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=40000x -cpuprofile profile_file.out
go test -bench=. -benchtime=40000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=30000x -cpuprofile profile_file.out
go test -bench=. -benchtime=30000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=20000x -cpuprofile profile_file.out
go test -bench=. -benchtime=20000x

sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=10000x -cpuprofile profile_file.out
go test -bench=. -benchtime=10000x

