ab -k -c 1 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00001_parallel_long3.txt
ab -k -c 2 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00002_parallel_long3.txt
ab -k -c 4 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00004_parallel_long3.txt
ab -k -c 8 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00008_parallel_long3.txt
ab -k -c 16 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00016_parallel_long3.txt
ab -k -c 32 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00032_parallel_long3.txt
ab -k -c 64 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00064_parallel_long3.txt
ab -k -c 128 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00128_parallel_long3.txt
ab -k -c 256 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00256_parallel_long3.txt
ab -k -c 512 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_00512_parallel_long3.txt
ab -k -c 1024 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_01024_parallel_long3.txt
ab -k -c 2048 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_02048_parallel_long3.txt
ab -k -c 4096 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_04096_parallel_long3.txt
ab -k -c 8192 -n 1000000 "http://10.20.0.4:3000/depts/random" > ab_08192_parallel_long3.txt

