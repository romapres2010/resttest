На глаза попалась не особо позитивное сравнение [Java vs GO. Тестирование большим числом пользователей](https://habr.com/ru/post/424649/).  

Решил сам проверить, действительно ли так все не радужно с Go.
Забегая вперед скажу, что при кэшировании в памяти и формировании JSON "на лету" удалось получить до 56000 [#/sec] на 4 физических ядра.  

Базовый сценарий GET запроса: 
- Если данные найдены в in memory кэше и они валидные, то формируем JSON из структуры
- Если данных в кэше нет, то ищем их в Bolt DB, если находим, то считываем готовый JSON 
- Если данных нет в Bolt DB, то запрашиваем их из БД, сохраняем их в in memory кэше
- Данные в in memory кэше накапливаются в буферном канале, после накопления около 10000 элементов они сбрасываются единым save в Bolt DB
- Если данные в БД менялись (update / insert) то через pg_notify передается уведомление и данные в кэше помечаются как невалидные, при следующем обращении они считываются заново из БД 

Ни каких RAID не использовалось. Интересно было замерить максимальные показатели на одном диске.
Краткие выводы: 
- использование кэширования JSON в Bolt DB дает прибавку в скорости примерно в 2 раза как на HDD, так и на SSD дисках
- перенос БД PostgreSQL с HDD на SDD дает прибавку минимум в 100 раз
- перенос Bolt DB на SDD с HDD так же дает прибавку минимум в 100 раз
- кэширование структуры в памяти и формирование JSON "на лету" дает прибавку в 10 раз по сравнению с  БД PostgreSQL на SDD
- стабильность отличная - ни каких Failed requests вплоть до 4096 одновременных подключений. 

Под катом результаты тестирования, и код тестового стенда [GitHub](https://github.com/romapres2010/resttest)
<cut />

Результаты тестирования, доступны по [ссылке](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark) 
Кому интересен код тестового стенда - добро пожаловать на [GitHub](https://github.com/romapres2010/resttest).   
Кода много около 4000 строк, но ни чего супер интересного использовано не было.  

Единственно, разве что кастомный [JSON encoder/decode](https://github.com/francoispqt/gojay) - это ускорило формирование JSON более чем в 2 раза. 

HTTP Get handler со стандартным пакетом JSON
```
BenchmarkHandler_DeptGetHandler-2   	 1000000	     41277 ns/op	4845.26 MB/s	   10941 B/op	     259 allocs/op
```
HTTP Get handler с пакетом francoispqt/gojay
```
BenchmarkHandler_DeptGetHandler-2   	 1000000	     13300 ns/op	15036.97 MB/s	    3265 B/op	      11 allocs/op
```

## Результаты тестирования

Параметры стенда для тестирования:
- БД PostgreSQL 11 под Windows 10, размер БД - до 30 ГБайт.
- Два таблицы - мастер до 10 000 000 строк, подчиненная таблица 100 000 000 строк
- Размер JSON сообщения 1500 байт, формируется из 1 мастер строки и 10 подчиненных строк
- компьютер - core i7-3770 - 4 core (8 thread), 16 Гбайт оперативка, HDD (WD 2.0TB WD2000fyyz),  SSD (Intel 530 Series) 
- Тестирование велось через ApacheBench, concurrency level от 1 до 4096, 1 000 000 запросов случайными образом. 
- Что бы минимизировать кэширование Windows, перед каждым тестом машина перегружалась.

В ходе тестирования проверялись следующие граничные случаи:
- [GET с прямым доступом к PostgreSQL - HDD диск (при каждом запросе идет чтение БД, JSON формируется из структуры). Один шпиндель](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark/10%20000%20000/1%20-%20pass%20no%20cache/hdd)
- [GET с прямым доступом к PostgreSQL - SSD диск (при каждом запросе идет чтение БД, JSON формируется из структуры). Один SSD диск](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark/10%20000%20000/1%20-%20pass%20no%20cache/ssd)
- [GET с кэширование в BOLT DB - HDD диск (JSON предварительно сохранен в BOLT DB)](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark/10%20000%20000/2%20-%20bolt%20populate%20no%20cache/HDD)
- [GET с кэширование в BOLT DB - SSD диск (JSON предварительно сохранен в BOLT DB)](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark/10%20000%20000/2%20-%20bolt%20populate%20no%20cache/SSD/long%20run)
- [GET с кэшированием структуры в памяти (JSON формируется из структуры)](https://github.com/romapres2010/resttest/tree/master/cmd/restDept/Apach%20benchmark/5%20000%20000/Core%20i7%203770/3%20-%20memory%20populate%20no%20cache)

Падение показателей при росте Concurrency Level связано с тем, что ApacheBench запускался на том же компьютере и активно "отъедал процессорное время".

### GET с прямым доступом к PostgreSQL - HDD (при каждом запросе идет чтение БД)

Concurrency Level | Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 36.84 [#/sec] (mean)| 27.145 [ms] (mean) |  116 [ms]
2  | 38.27 [#/sec] (mean)| 52.262  [ms] (mean) |  165 [ms]
4  | 41.34 [#/sec] (mean)| 96.755 [ms] (mean) |  268 [ms]
8  | 45.02 [#/sec] (mean)| 177.687  [ms] (mean) |  420 [ms]
16  | 47.56 [#/sec] (mean)| 336.428 [ms] (mean) |  813 [ms]
128  | 49.19 [#/sec] (mean)| 2602.228 [ms] (mean) |  13055 [ms]
512  | 50.5 [#/sec] (mean)| 10122.343 [ms] (mean) | 19394 [ms]
2048  | 51.91 [#/sec] (mean)| 39453.681 [ms] (mean) | 57018 [ms]

### GET с прямым доступом к PostgreSQL - SSD (при каждом запросе идет чтение БД)

Concurrency Level | Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 713.82 [#/sec] (mean)| 1.401 [ms] (mean) |  2 [ms]
2  | 1914.61 [#/sec] (mean)| 1.045  [ms] (mean) |  2 [ms]
4  | 3326.52 [#/sec] (mean)| 1.202  [ms] (mean) |  2 [ms]
8  | 4599.95 [#/sec] (mean)| 1.739  [ms] (mean) |  4 [ms]
16  | 4599.80 [#/sec] (mean)| 3.478 [ms] (mean) |  9 [ms]
128  | 5243.76 [#/sec] (mean)| 24.410 [ms] (mean) |  102 [ms]
512  | 5354.35 [#/sec] (mean)| 95.623 [ms] (mean) | 506 [ms]
2048  | 5285.83 [#/sec] (mean)| 387.451 [ms] (mean) | 2871 [ms]

### GET с кэширование в BOLT DB - HDD диск (JSON предварительно сохранен в BOLT DB)

Concurrency Level| Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 81.55 [#/sec] (mean)| 12.262  [ms] (mean) |  38 [ms]
2  | 67.04 [#/sec] (mean)| 29.832  [ms] (mean) |  97 [ms]
4  | 72.51 [#/sec] (mean)| 55.167  [ms] (mean) |  183 [ms]
8  | 92.48 [#/sec] (mean)| 86.502  [ms] (mean) |  291 [ms]
16  | 89.42 [#/sec] (mean)| 178.923  [ms] (mean) |  550 [ms]
128  | 86.76 [#/sec] (mean)| 1475.378 [ms] (mean) |  11280 [ms]
512  | 92.38 [#/sec] (mean)| 2771.145 [ms] (mean) | 10238 [ms]
2048  | 100.56 [#/sec] (mean)| 20366.847 [ms] (mean) | 18632 [ms]
4096  | 106.98 [#/sec] (mean)| 38289.063 [ms] (mean) | 37074[ms]

### GET с кэширование в BOLT DB - SSD диск (JSON предварительно сохранен в BOLT DB)

Concurrency Level| Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 3411.07 [#/sec] (mean)| 0.293  [ms] (mean) |  1 [ms]
2  | 7468.21 [#/sec] (mean)| 0.268  [ms] (mean) |  1 [ms]
4  | 9501.15 [#/sec] (mean)| 0.421  [ms] (mean) |  2 [ms]
8  | 10481.68 [#/sec] (mean)| 0.763  [ms] (mean) |  3 [ms]
16  | 10052.14 [#/sec] (mean)| 1.592  [ms] (mean) |  5 [ms]
128  | 10754.02 [#/sec] (mean)| 11.903 [ms] (mean) |  20 [ms]
512  | 11030.61 [#/sec] (mean)| 46.416 [ms] (mean) | 66 [ms]
2048  | 10634.72 [#/sec] (mean)| 192.577 [ms] (mean) | 362 [ms]
4096  | 10659.04 [#/sec] (mean)| 384.275 [ms] (mean) | 720 [ms]

### GET с кэшированием структуры в памяти

Concurrency Level| Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 9178.22 [#/sec] (mean)| 0.109 [ms] (mean) |  1 [ms]
2  | 22580.40 [#/sec] (mean)| 0.089  [ms] (mean) |  1 [ms]
4  | 36163.33 [#/sec] (mean)| 0.111  [ms] (mean) |  1 [ms]
8  | 56109.17 [#/sec] (mean)| 0.143  [ms] (mean) |  1 [ms]
16  | 43942.75 [#/sec] (mean)| 0.364  [ms] (mean) |  2 [ms]
128  | 55005.53 [#/sec] (mean)| 2.327 [ms] (mean) |  6 [ms]
512  | 35338.01 [#/sec] (mean)| 14.489 [ms] (mean) | 25 [ms]
2048  | 38090.35 [#/sec] (mean)| 53.767 [ms] (mean) | 228 [ms]
4096  | 30196.47 [#/sec] (mean)| 135.645 [ms] (mean) | 609 [ms]

# Профилирование 

По данным Benchmark, HTTP Get handler отрабатывает за 13300 ns
``` go
BenchmarkHandler_DeptGetHandler-2   	 1000000	     13300 ns/op	15036.97 MB/s	    3265 B/op	      11 allocs/op
```

По данным ApacheBench, лучшее mean выполнения запроса - 0.088 ms.  
Стало интересно где разница 0.088 - 0.013 = 0.055 ms.  
Включил профилирование - [результаты](https://raw.githubusercontent.com/romapres2010/resttest/master/img/pprof.png)

На верхнем уровне:
- net/http.(*conn).serve - (61.85%)
- net/http.(*connReader).backgroundRead - (7.04%)
- runtime.gcBgMarkWorker - (18.30%)

Чем занимался HTTP server из (61.85%):  
- Считывание входящих запросов - net/http.(*conn).readRequest - (9.29%). Входящие запросы были без тела.
- Запись исходящих ответов - net/http.(*response).finishRequest - (21.74%). Запись тела ответа.
- Собственно обработчик запроса - net/http.(*ServeMux).ServeHTTP - (23.33%).

![](https://raw.githubusercontent.com/romapres2010/resttest/master/img/http.conn.serve.png)

Из суммарных затрат на readRequest, finishRequest и backgroundRead (39.66%), на обработку системных вызовов IO Windows - internalpoll.(ioSrv).ExecIO - пришлось (24.78%).  

Посмотрим чем занимался наш основной обработчик net/http.(*ServeMux).ServeHTTP:
- Парсинг URL и обработка параметров - github.com/gorilla/mux.(*Route).Match - (2.29%)
- Формирование JSON из структуры в памяти - github.com/francoispqt/gojay.marshal - (9.92%)
- Работка с кэшем в памяти - (4.53%)

![](https://raw.githubusercontent.com/romapres2010/resttest/master/img/http.ServeMux.ServeHTTP.png)