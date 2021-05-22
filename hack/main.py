for i in range(100000, 9999, -10000):
    a = 'sudo perf stat -e L1-dcache-load-misses go test -bench=^BenchmarkLockFreeLookupAndIns -benchtime=%sx -cpuprofile profile_file.out' % i
    b = 'go test -bench=. -benchtime=%sx' % i
    print(a)
    print(b)
    print()
