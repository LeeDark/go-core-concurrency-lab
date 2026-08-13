# Конкурентність і паралелізм

> Англійське canonical-джерело: [cheatsheet-concurrency.md](cheatsheet-concurrency.md).
> Російський переклад: [cheatsheet-concurrency.ru.md](cheatsheet-concurrency.ru.md).
> Скорочена версія: [cheatsheet-concurrency-simplified.md](cheatsheet-concurrency-simplified.md).
> Технічні терміни Go — `goroutine`, `channel`, `worker`, `handler`, `select` — залишено
> англійською для відповідності коду та API.

## Визначення

**Конкурентність (concurrency)** — здатність організувати програму як незалежні обчислювальні
частини, які можуть просуватися вперед не в строгому або частковому порядку, не змінюючи очікуваний
результат. Частини можуть чергуватися в часі або виконуватися паралельно.

**Паралелізм (parallelism)** — одночасне виконання кількох частин роботи.

Коротко: **конкурентність описує проєктування програми, а паралелізм — її виконання.** Конкурентний
код не стає автоматично швидшим за послідовний.

## Модель конкурентного програмування

Communicating Sequential Processes (CSP) — математична модель, що вплинула на Go. Незалежні
послідовні процеси обмінюються повідомленнями. Обмін може бути синхронним: відправник не продовжує
роботу, доки отримувач не прийме повідомлення. Так працює небуферизований Go `channel`.

У конкурентній програмі кількість можливих порядків виконання швидко зростає. Через це складніше
аналізувати весь стан програми та легше отримати deadlock, data race або помилку порядку.

## Продуктивність і планувальник Go

Конкурентний код має витрати на створення та планування goroutine, синхронізацію, канали й обмін
даними. Конкурентність доречна для незалежної роботи CPU, незалежних I/O-операцій, конвеєрів,
worker pool та обробки запитів. Вона може бути невигідною для дуже малих або строго послідовних
задач.

Планувальник Go часто описують моделлю G-M-P:

```text
G — goroutine
M — OS thread
P — logical processor, потрібний для виконання Go-коду
```

`GOMAXPROCS` обмежує кількість `P` і, відповідно, кількість OS thread, які можуть одночасно
виконувати користувацький Go-код. Це не обмежує загальну кількість OS thread: runtime може створити
додаткові thread, коли деякі з них заблоковані в системних викликах.

Типові стани goroutine:

```text
running   — виконується зараз
runnable  — готова до виконання та очікує процесорного часу
waiting   — заблокована на I/O, channel, mutex, syscall або іншій події
```

## Goroutine

Goroutine — легка одиниця конкурентного виконання, якою керує Go runtime:

```go
go work()
```

Goroutine може отримувати аргументи, але не повертає значення безпосередньо викликувачу. Результати
передають через channel, спільний стан або інший механізм координації.

Goroutine починається після виконання `go` і завершується після повернення її функції. Якщо завершується
`main`, завершується вся програма, навіть якщо інші goroutine ще працюють.

**Витік goroutine** виникає, коли goroutine ніколи не завершується — наприклад, назавжди заблокована
на channel, mutex, I/O або через відсутній сигнал скасування. `sync.WaitGroup` координує завершення,
але не скасовує роботу й не збирає помилки.

`context.Context` передає скасування, дедлайни, тайм-аути та дані запиту через межі API. Він дає
goroutine сигнал, що робота більше не потрібна.

## Channels

Channel — типізований канал даних і водночас механізм синхронізації між goroutine.

### Небуферизований channel

У небуферизованого channel немає місткості. Відправлення блокується, доки інша goroutine не готова
прийняти значення, а отримання — доки хтось не відправить значення. Це точка синхронізації.

### Буферизований channel

У буферизованого channel є фіксована місткість:

- відправлення блокується, коли буфер заповнений;
- отримання блокується, коли буфер порожній;
- буфер змінює блокування, але не скасовує правила володіння та синхронізації.

### Nil і закриття

Нульове значення channel — `nil`. Відправлення в nil channel або отримання з нього блокується
назавжди. У `select` case з nil channel вимкнений.

Закриття channel означає, що нових значень більше не буде. Отримувач спочатку отримує буферизовані
значення, а потім одержує zero value з `ok == false`:

```go
value, ok := <-ch
```

Відправлення в закритий channel і повторне закриття викликають panic. Закривати channel для garbage
collection не потрібно.

### Володіння channel

Сторона, яка відправляє значення, зазвичай відповідає за закриття channel. Отримувач не повинен
закривати channel, яким він не володіє.

```go
var receiveOnly <-chan int
var sendOnly chan<- int
```

Один channel може мати кількох отримувачів, але кожне відправлене значення отримує лише один із них.
Це зручно для розподілу роботи між workers.

### Select, тайм-аут і скасування

`select` очікує на кілька операцій channel. Якщо готові кілька case, обирається один із них.
`default` робить операцію неблокувальною, але необережне використання може створити busy loop,
завантажити CPU або зламати backpressure.

Тайм-аут можна реалізувати через `time.After` або deadline у context. Скасування зазвичай обробляють
через `ctx.Done()`:

```go
select {
case value := <-ch:
	_ = value
case <-ctx.Done():
	return ctx.Err()
}
```

## Worker Pool v1

Worker Pool v1 запускає фіксовану кількість goroutine для потоку незалежних задач:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
```

Це **bounded concurrency**: одночасно виконується не більше `workerCount` handlers.

### API

```go
type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID int
	Value string
	Err   error
}

func Run(workerCount int, jobs <-chan Job, handle func(Job) Result) <-chan Result
```

Контракт v1:

- `workerCount <= 0` перетворюється на одного worker;
- caller відправляє jobs і закриває `jobs`;
- pool відправляє results і закриває `results`;
- `handle` має бути non-nil і не повинен викликати panic;
- pool не відновлює panic handler;
- nil або ніколи не закритий `jobs` залишає workers в очікуванні назавжди;
- порядок results не гарантується: вони приходять у міру завершення workers, а не за порядком
  відправлення або `Job.ID`;
- caller повинен читати `results`, доки channel не закриється.

### Безпечний порядок роботи

Обидва channel у v1 небуферизовані. Producer і consumer повинні працювати одночасно:

```go
jobs := make(chan Job)
results := Run(3, jobs, handle)

go func() {
	defer close(jobs)
	for _, job := range submittedJobs {
		jobs <- job
	}
}()

for result := range results {
	consume(result)
}
```

Не можна спочатку синхронно відправити всі jobs, а потім читати results: workers можуть заблокуватися
на відправленні result, а producer — на відправленні наступної job. Це може призвести до deadlock.

### Володіння channel і завершення

`jobs`:

- caller відправляє та закриває;
- workers отримують;
- workers не закривають.

`results`:

- workers відправляють;
- caller отримує;
- pool закриває після завершення всіх workers через `sync.WaitGroup`.

V1 зберігає помилки окремих jobs у `Result.Err`, але не скасовує pool і не повторює роботу.
Cancellation, тайм-аути, graceful shutdown, leak checks і розширена політика помилок належать v2.

## Race detector / memory model / scheduler

У наступних темах потрібно розглянути:

- data race;
- happens-before;
- mutex;
- atomic;
- `GOMAXPROCS`;
- CPU-bound проти I/O-bound роботи.

## Перевірка Worker Pool v1

```bash
go test ./07-worker-pool-v1/workerpool
go test -race ./07-worker-pool-v1/workerpool
```

## Worker Pool v2: керування життєвим циклом

V2 зберігає фіксований worker pool і додає кооперативне скасування через `context.Context`:

```go
func Run(
	ctx context.Context,
	workerCount int,
	jobs <-chan Job,
	handle func(context.Context, Job) Result,
) <-chan Result
```

Caller створює та скасовує `ctx`, надсилає jobs і закриває `jobs`, а також читає `results`. Pool
ніколи не закриває і не drain-ить `jobs`, яким володіє caller; він надсилає в `results` і закриває
`results`. Непозитивне число workers нормалізується до одного, порядок results не гарантується.

Workers мають враховувати скасування і під час очікування jobs, і під час публікації results:

```go
for {
	select {
	case job, ok := <-jobs:
		if !ok {
			return
		}
		result := handle(ctx, job)
		if ctx.Err() != nil {
			return
		}
		select {
		case results <- result:
		case <-ctx.Done():
			return
		}
	case <-ctx.Done():
		return
	}
}
```

Скасування є кооперативним. Handler має спостерігати `ctx.Done()` або використовувати операції,
що підтримують context; pool не може примусово перервати код, який ігнорує context. Producer також
має обирати між надсиланням job і `ctx.Done()`, а потім закривати `jobs`, якщо він володіє цим
життєвим циклом:

```go
for _, job := range submittedJobs {
	select {
	case jobs <- job:
	case <-ctx.Done():
		return
	}
}
close(jobs)
```

Після скасування jobs, що залишилися в черзі, можуть бути не оброблені. Якщо надсилання result і
скасування стають готовими одночасно, result може бути надісланий або відкинутий. Тому caller має
сприймати скасування як завершення операції й не очікувати останній result для кожної job.

Pool зберігає правило v1: `sync.WaitGroup` відстежує завершення workers, а coordinator закриває
`results` лише після `wg.Wait()`. Скасування не закриває `results` безпосередньо.

Для загального тайм-ауту операції caller створює похідний context:

```go
ctx, cancel := context.WithTimeout(parent, timeout)
defer cancel()
results := workerpool.Run(ctx, workers, jobs, handle)
```

Помилки окремих jobs залишаються в `Result.Err`; окремого errors channel немає. Загальний тайм-аут
видимий через `context.DeadlineExceeded` у context handler’а. Тайм-аут окремої job використовує
дочірній context і не скасовує весь pool:

```go
handle := workerpool.WithJobTimeout(jobTimeout, func(jobCtx context.Context, job Job) Result {
	return doJob(jobCtx, job)
})
results := workerpool.Run(ctx, workers, jobs, handle)
```

Wrapped handler має спостерігати `jobCtx.Done()`. Якщо після завершення дочірнього context він
повернув result без власної помилки, wrapper записує помилку context у `Result.Err`. Тайм-аут однієї
job не скасовує весь pool. Retries, metrics, tracing і persistent queues не входять до цієї
лабораторії.

Перевірка:

```bash
go test ./08-worker-pool-v2/workerpool
go test -race ./08-worker-pool-v2/workerpool
```
