go.exe test -benchmem -benchtime 10x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_10rep.txt 
go.exe test -benchmem -benchtime 100x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_100rep.txt 
go.exe test -benchmem -benchtime 1000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_1000rep.txt 
go.exe test -benchmem -benchtime 10000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_10000rep.txt 
go.exe test -benchmem -benchtime 100000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_100000rep.txt 
go.exe test -benchmem -benchtime 1000000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_1000000rep.txt 
rem go.exe test -benchmem -benchtime 10000000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_10000000rep.txt 
rem go.exe test -benchmem -benchtime 100000000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_100000000rep.txt 
rem go.exe test -benchmem -benchtime 1000000000x -bench ^(BenchmarkHandler_DeptEmpsGetHandler)$ > 1bench_Handler_DeptEmpsGetHandler_1000000000rep.txt 
