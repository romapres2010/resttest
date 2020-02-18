package json

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

// PopulateDeptCache populate JSON Cahce and Bolt DB
// =====================================================================
func (s *Service) PopulateDeptCache() {
	cnt := "json (s *Service)  PopulateDeptCache" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	var err error

	// Считаем все PK
	if testDeptPK == nil {
		// считаем все PK для dept
		testDeptPK, err = s.deptService.GetDeptsPK()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetDeptsPK()", cnt)
			log.Printf(errM)
			return
		}
	}

	//	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU() / 2
	log.Printf("[INFO] %v - runtime.NumCPU()= '%v'", cnt, numCPU)

	lenPK := len(testDeptPK)
	log.Printf("[INFO] %v - lenPK= '%v'", cnt, lenPK)

	var waitgroup sync.WaitGroup // для ожидания всех запущенных горутин

	for i := 0; i < numCPU; i++ {
		log.Printf("[INFO] %v - start gorutine '%v'", cnt, i)
		waitgroup.Add(1) // добавляем в список ожидания горутин

		go func(index int) {
			defer PrintElapsed(fmt.Sprintf("[INFO] %v - gorutine '%v' - takes", cnt, index))()
			for j := 0; j < lenPK/numCPU; j++ {
				pk := testDeptPK[j+(index*(lenPK/numCPU))].Deptno
				buf, _ := s.GetDeptEmps(pk)
				s.PutByteSlice(buf) // вернем в pool для повторного использования. Safe for nil pointer
			}
			waitgroup.Done() // сигнал завершения горутины
		}(i)

	}
	waitgroup.Wait() // ожидание завершения всех горутин

	log.Printf("[INFO] %v - SUCCESS", cnt)
	return
}

// PrintElapsed при отложенном запуске в начале процедуры измеряет время ее выполнение (с закрытием)
func PrintElapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s%v", what, time.Since(start))
	}
}
