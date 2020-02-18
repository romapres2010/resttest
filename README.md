# Тестирование REST API на Golang. 14000 GET/sec per 1 core не предел?

## Почему Golang

Компания, в которой я работаю, специализируется на решениях Oracle. Обычно REST API мы делаем с использованием [Oracle REST Data Services](https://www.oracle.com/database/technologies/appdev/rest.html). Но импортозамещение заставило нас обратить внимание в сторону свободного ПО. 
Параллельно шел проект на Blockcahin, контракты для которого планировалось писать на Go. Почему бы не попробовать Go и для REST API.   
Тут еще на глаза попалась не сообо позитивное сравнение [Java vs GO. Тестирование большим числом пользователей](https://habr.com/ru/post/424649/).  

Решил сам проверить, действительно ли так все не радужно с Go.
Забегая вперед скажу, что при кешировании структуры в памяти и формировании JSON удалось получить до 14000 GET/sec на 1 физическое ядро.  
И не одного Failed requests вплоть до 4096 одновременных подключений. 

Параметры стенда для тестирования:
- БД PostgreSQL под Windows 10
- Два таблицы - мастер до 10 000 000 строк, подчиненная таблица 100 000 000 строк
- Размер JSON сообщения 1500 байт, формируется из 1 мастер строки и 10 подчиненных строк
- компьютер - core i7 - 4 core (8 thread), 16 Гбайт оперативка, HDD / SSD диск
- Тестирование велось через ApacheBench, concurrency level от 1 до 4096, 1 000 000 запросов со случайными ID

Базовый сценарий GET запроса: 
- Если данные найдены в in memory кэше и они валидные, то формируем JSON из структуры
- Если данных в кеше нет, то ищем их в Bolt DB, если находим, то считываем готовый JSON 
- Если данных нет в Bolt DB, то запрашиваем их из БД, сохраняем их в in memory кэше
- Данные в in memory кэше накапливаются в буферном канале, после накопления около 10000 элементов они сбрасываются единым save в Bolt DB
- Если данные в БД менялись (update / insert) то через pg_notify передавается уведомление и данные в кэше помечались как невалидные и при следующем обращении они считывались заново из БД 

В ходе тестирования проверялись следующие граничные случаи:
- GET с прямым доступом к PostgreSQL (при каждом запросе идет чтение БД, JSON формируется из структуры)
- GET с кеширование в BOLT DB (JSON предварительно сохранен в BOLT DB)
- GET с кешированием структуры в памяти (JSON формируется из структуры) 

Под катом результаты тестирования. 
Кому интересен код - добро пожаловать в репозиторий.



```
type Dept struct {
	Deptno      int            `db:"deptno" json:"deptNumber" validate:"required"`
	Dname       string         `db:"dname" json:"deptName"`
	Loc         sql.NullString `db:"loc" json:"deptLocation, nullempty"`
	Emps        []*Emp         `json:"emps"` // массив указателей на дочерние emp
	EmpsPK      []*EmpPK       `json:"-"`    // массив PK дочерних emp
	IsPutToPool bool           `json:"-"`    // признак, что структура возвращена в pool
}

type Emp struct {
	Empno       int            `db:"empno" json:"empNo"`
	Ename       sql.NullString `db:"ename" json:"empName, nullempty"`
	Job         sql.NullString `db:"job" json:"job, nullempty"`
	Mgr         sql.NullInt64  `db:"mgr" json:"mgr, omitempty"`
	Hiredate    sql.NullString `db:"hiredate" json:"hiredate, nullempty"`
	Sal         sql.NullInt64  `db:"sal" json:"sal, nullempty"`
	Comm        sql.NullInt64  `db:"comm" json:"comm, nullempty"`
	Deptno      sql.NullInt64  `db:"deptno" json:"deptNumber, nullempty"`
	Dept        *Dept          `json:"-"` // связь с Dept для кеширования
	IsPutToPool bool           `json:"-"` // признак, что структура возвращена в pool
}

```


Тестовый сценарий был достаточно простой: 


При кешировании из памяти удалось получить до 14 000 Get / сек на 1 ядро.  
И не одного Failed requests вплоть до 4096 одновременных подключений.  
Ниже приведены результаты тестирования. Падение показателей при росте Concurrency Level связано с тем, что ApacheBench запускался на том же компьютере и активно "отъедал процессорное время".

## Результаты тестирования
- **GET с прямым доступом к PostgreSQL SSD (при каждом запросе идет чтение БД)** 

Concurrency Level | Requests per second | Time per req. | 99% percentile
------------ | ------------- | -------------  | ------------- 
1  | 713.82 [#/sec] (mean)| 1.401 [ms] (mean) |  2 [ms]
2  | 1914.61 [#/sec] (mean)|1.045  [ms] (mean) |  2 [ms]
4  | 3326.52 [#/sec] (mean)| 1.202  [ms] (mean) |  2 [ms]
8  | 4599.95 [#/sec] (mean)| 1.739  [ms] (mean) |  4 [ms]
16  | 4599.80 [#/sec] (mean)| 3.478 [ms] (mean) |  9 [ms]
128  | 5243.76 [#/sec] (mean)| 24.410 [ms] (mean) |  102 [ms]
512  | 5354.35 [#/sec] (mean)| 95.623 [ms] (mean) | 506 [ms]
2048  | 5285.83 [#/sec] (mean)| 387.451 [ms] (mean) | 2871 [ms]

- **GET с кеширование в BOLT DB (JSON предварительно сохранен в BOLT DB)**

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

- **GET с кешированием структуры в памяти** 

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

