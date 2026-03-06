package tgapi

import (
	"context"
	"sync"
)

// workerPool — приватная структура, управляющая пулом воркеров.
// Внешний код не может создавать или напрямую взаимодействовать с этой структурой.
// Используется только через экспортируемые методы newWorkerPool, start, stop, submit.
type workerPool struct {
	taskCh    chan requestEnvelope // канал для принятия задач (буферизованный)
	queueSize int                  // максимальный размер очереди
	workers   int                  // количество воркеров (горутин)
	wg        sync.WaitGroup       // синхронизирует завершение всех воркеров при остановке
	quit      chan struct{}        // канал для сигнала остановки
	started   bool                 // флаг, указывающий, запущен ли пул
	startedMu sync.Mutex           // мьютекс для безопасного доступа к started
}

// requestEnvelope — приватная структура, инкапсулирующая задачу и канал для результата.
// Используется только внутри пакета для передачи задач воркерам.
type requestEnvelope struct {
	doFunc   func(context.Context) (any, error) // функция, выполняющая запрос
	resultCh chan requestResult                 // канал, через который воркер вернёт результат
}

// requestResult — приватная структура, представляющая результат выполнения задачи.
// Внешний код получает его через канал, но не знает структуры — только через <-chan requestResult.
type requestResult struct {
	value any   // значение, возвращённое задачей
	err   error // ошибка, если возникла
}

// newWorkerPool создаёт новый пул воркеров с заданным количеством горутин и размером очереди.
// Это единственный способ создать workerPool — внешний код не может создать его напрямую.
func newWorkerPool(workers int, queueSize int) *workerPool {
	if workers <= 0 {
		workers = 1 // защита от некорректных значений
	}
	if queueSize <= 0 {
		queueSize = 100 // разумный дефолт
	}

	return &workerPool{
		taskCh:    make(chan requestEnvelope, queueSize),
		queueSize: queueSize,
		workers:   workers,
		quit:      make(chan struct{}),
	}
}

// start запускает воркеры (горутины), которые будут обрабатывать задачи из очереди.
// Метод идемпотентен: если пул уже запущен — ничего не делает.
// Должен вызываться перед первым вызовом submit.
func (p *workerPool) start(ctx context.Context) {
	p.startedMu.Lock()
	defer p.startedMu.Unlock()
	if p.started {
		return // уже запущен — ничего не делаем
	}
	p.started = true

	// Запускаем воркеры — каждый будет обрабатывать задачи в бесконечном цикле
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx) // запускаем горутину с контекстом
	}
}

// stop останавливает пул воркеров.
// Отправляет сигнал остановки через quit-канал и ждёт завершения всех активных задач.
// Безопасно вызывать многократно — после остановки повторные вызовы не имеют эффекта.
func (p *workerPool) stop() {
	close(p.quit) // сигнал для всех воркеров — выйти из цикла
	p.wg.Wait()   // ждём, пока все воркеры завершатся
}

// submit отправляет задачу в очередь и возвращает канал, через который будет получен результат.
// Если очередь переполнена — возвращает ErrPoolQueueFull.
// Канал результата имеет буфер 1, чтобы не блокировать воркера при записи.
// Контекст используется для отмены задачи, если клиент отменил запрос до отправки.
func (p *workerPool) submit(ctx context.Context, do func(context.Context) (any, error)) (<-chan requestResult, error) {
	// Проверяем, не превышена ли очередь
	if len(p.taskCh) >= p.queueSize {
		return nil, ErrPoolQueueFull
	}

	// Создаём канал для результата — буферизованный, чтобы не блокировать воркера
	resultCh := make(chan requestResult, 1)

	// Создаём обёртку задачи
	envelope := requestEnvelope{
		doFunc:   do,
		resultCh: resultCh,
	}

	// Пытаемся отправить задачу в очередь
	select {
	case <-ctx.Done():
		// Клиент отменил операцию до отправки — возвращаем ошибку отмены
		return nil, ctx.Err()
	case p.taskCh <- envelope:
		// Успешно отправлено — возвращаем канал для чтения результата
		return resultCh, nil
	default:
		// Очередь переполнена — не должно происходить при проверке len(p.taskCh), но на всякий случай
		return nil, ErrPoolQueueFull
	}
}

// worker — приватная горутина, выполняющая задачи из очереди.
// Каждый воркер работает в бесконечном цикле, пока не получит сигнал остановки.
// При получении задачи:
//   - вызывает doFunc с контекстом
//   - записывает результат в resultCh
//   - закрывает канал, чтобы клиент мог прочитать и завершить
//
// После закрытия quit-канала — воркер завершает работу.
func (p *workerPool) worker(ctx context.Context) {
	defer p.wg.Done() // уменьшаем WaitGroup при завершении горутины

	for {
		select {
		case <-p.quit:
			// Получен сигнал остановки — выходим из цикла
			return

		case envelope := <-p.taskCh:
			// Выполняем задачу с переданным контекстом (клиентский или общий)
			value, err := envelope.doFunc(ctx)

			// Записываем результат в канал — не блокируем, т.к. буфер 1
			envelope.resultCh <- requestResult{
				value: value,
				err:   err,
			}
			// Закрываем канал — клиент знает, что результат пришёл и больше не будет
			close(envelope.resultCh)
		}
	}
}
