# Worker Pool v1

## Final Plan

### 01: Read basics

```text
- Effective Go
- Go Code Review Comments
- 100 Go Mistakes
- Concurrency in Go
```

Topics:

* goroutines;
* goroutine lifetimes;
* channels, buffered / unbuffered;
* channel buffer size;
* channel ownership;
* select;
* WaitGroup;
* buffering;
* backpressure;
* worker pool;
* concurrency mistakes;
* channel misuse;
* avoiding leaks;
* race conditions, race detector;

* contexts;
* cancellation, context cancellation;
* timeout;
* shutdown;
* context misuse;

Interview questions:

1. Что такое goroutine и чем она отличается от OS thread?
2. Как остановить goroutine?
3. Что такое goroutine leak?
4. Чем buffered channel отличается от unbuffered?
5. Кто должен закрывать channel?
6. Что произойдёт при чтении из закрытого channel?
7. Что произойдёт при записи в закрытый channel?
8. Как работает select?
9. Для чего нужен default в select и чем он опасен?
10. Как реализовать timeout?
11. Как реализовать cancellation?
12. Чем context cancellation отличается от closing done channel?
13. Что такое worker pool?
14. Что такое fan-out/fan-in?
15. Когда использовать mutex вместо channel?
16. Когда использовать atomic?
17. Что находит race detector?
18. Что race detector не гарантирует?
19. Как сделать graceful shutdown HTTP-сервиса?
20. Как диагностировать блокировки и latency в concurrent Go-приложении?

---

### 02: Write lab code

Минимальная структура:

```text
06-worker-pool-v1/
  README.md
  main.go
  workerpool/
    pool.go
    pool_test.go
```

Первое упражнение:

```text
jobs -> workers -> results
```

```text
Implement a worker pool that:
- accepts jobs;
- runs fixed number of workers;
- supports context cancellation;
- returns results;
- handles errors;
- does not leak goroutines;
- passes go test -race.
```

Последовательность написания:

```text
- Too simple worker pool
    - handles errors;
    - unit tests;
    - accepts jobs, Job struct;
    - runs fixed number of workers;
    - returns results, Result struct;
    - result channel closing;
    - context?;
    - WaitGroup?;
- Supports context cancellation; (v2)
- Does not leak goroutines; (v2)
- Passes go test -race; (v2)
```

И сразу interview angle:

> “I usually avoid unbounded goroutine creation for message/request processing. I prefer bounded worker pools or semaphores when throughput and external dependencies need to be controlled.”

Вот это хорошая Senior-фраза. Не пафосная, не выдуманная, и не звучит как “я вчера посмотрел туториал на YouTube на скорости 1.75x”.

---

## First Plan - День 2 — Effective Go: Concurrency

### Цель

Освежить базовые goroutines/channels/select без углубления в адскую математику scheduling.

### Читать

Из Effective Go:

```text
- Goroutines
- Channels
- Channels of channels
- Parallelization
- Leaky buffer
```

### Практика

Начать:

```text
06-worker-pool-v1/
```

Минимальный worker pool:

```text
jobs -> workers -> results
```

### Шпаргалка

В `docs/cheatsheet-concurrency.md` добавить:

```text
- goroutine lifecycle
- unbuffered channel
- buffered channel
- closing channels
- select
- basic worker pool
```

### Interview questions

Добавить:

```text
1. Что такое goroutine?
2. Buffered vs unbuffered channel?
3. Кто должен закрывать channel?
4. Что будет при чтении из закрытого channel?
5. Как работает select?
```

## ChatGPT Plan - День 6 — Mini review + first concurrency bridge

### Темы

```text
race detector
maps are not thread-safe
mutex preview
context cancellation preview
```

### Lab

Начать:

```text
06-worker-pool/
```

Но очень маленько:

```text
- 3 workers
- jobs channel
- results channel
- WaitGroup
- go test -race
```

Без сложной cancellation пока.

## CONCURRENCY.md - 2. Что покрываем по модулям

### 01-worker-pool

Темы:

* goroutine lifecycle;
* buffered/unbuffered channels;
* worker pool;
* backpressure;
* closing channels;
* error handling;
* context cancellation;
* wait groups;
* race detector.

Практика:

```text
Задача:
Есть очередь задач.
Нужно обработать N jobs с M workers.
Добавить:
- context cancellation;
- timeout;
- result channel;
- error channel или result struct;
- graceful stop;
- тесты;
- go test -race.
```

Что должен уметь объяснить на интервью:

* зачем worker pool;
* почему нельзя запускать бесконечно goroutine на каждый request/message;
* чем buffered channel отличается от unbuffered;
* кто закрывает channel;
* почему закрывать channel должен producer;
* как не получить deadlock;
* как не получить goroutine leak.

## CONCURRENCY.md - 3. Cheat sheet: структура на 1–3 страницы

Я бы сделал не длинный конспект, а **боевую шпаргалку для собеседования**.

### Разделы cheat sheet

```text
1. Goroutines
   - что это;
   - стоимость;
   - lifecycle;
   - leaks.

2. Channels
   - unbuffered vs buffered;
   - send/receive/close;
   - ownership;
   - common deadlocks.

3. select
   - multiplexing;
   - timeout;
   - cancellation;
   - default case risks.

4. Context
   - cancellation;
   - timeout;
   - deadline;
   - values: осторожно.

5. Patterns
   - worker pool;
   - fan-out/fan-in;
   - pipeline;
   - rate limiting;
   - graceful shutdown.

6. Shared state
   - mutex;
   - RWMutex;
   - atomic;
   - channels;
   - race detector.

7. Production checklist
   - bounded concurrency;
   - cancellation;
   - timeouts;
   - error handling;
   - shutdown;
   - metrics/logging;
   - pprof/block profile.
```

## CONCURRENCY.md - 8. Первый шаг

Начинать лучше с `06-worker-pool-v1`, потому что он связывает половину тем сразу:

```text
goroutines
channels
context
WaitGroup
buffering
backpressure
error handling
shutdown
race detector
```

Минимальная структура:

```text
01-worker-pool/
  README.md
  main.go
  workerpool/
    pool.go
    pool_test.go
```

Первое упражнение:

```text
Implement a worker pool that:
- accepts jobs;
- runs fixed number of workers;
- supports context cancellation;
- returns results;
- handles errors;
- does not leak goroutines;
- passes go test -race.
```

И сразу interview angle:

> “I usually avoid unbounded goroutine creation for message/request processing. I prefer bounded worker pools or semaphores when throughput and external dependencies need to be controlled.”

Вот это хорошая Senior-фраза. Не пафосная, не выдуманная, и не звучит как “я вчера посмотрел туториал на YouTube на скорости 1.75x”.

## CONCURRENCY.md - День 1 — Setup + mental model

### Reading, 20–30 минут

Прочитать выборочно:

```text
Effective Go:
- Goroutines
- Channels
- Errors, если останется время

Go Code Review Comments:
- Goroutine Lifetimes
- Channel Buffer Size
```

### Practice

Создать структуру проекта.

```bash
mkdir go-core-concurrency-lab
cd go-core-concurrency-lab

go mod init github.com/your-name/go-core-concurrency-lab

mkdir -p docs
mkdir -p labs/01-worker-pool/cmd/demo
mkdir -p labs/01-worker-pool/workerpool
mkdir -p labs/02-pipeline-cancellation/cmd/demo
mkdir -p labs/02-pipeline-cancellation/pipeline
mkdir -p labs/03-rate-limited-fetcher
mkdir -p labs/04-shared-state-mutex-vs-channel
mkdir -p labs/05-graceful-shutdown
```

Создать первый простой worker pool:

```text
jobs -> workers -> results
```

Пока без идеального shutdown. Просто чтобы код заработал.

### Done

```text
[ ] repo создан
[ ] go test ./... работает
[ ] есть простейший worker pool
[ ] есть README.md в корне
[ ] есть docs/cheatsheet.md
```

---

## CONCURRENCY.md - День 2 — Worker Pool v1

### Practice

Довести `06-worker-pool-v1` до нормальной учебной версии.

Добавить:

```text
- фиксированное количество workers
- Job struct
- Result struct
- context.Context
- обработку ошибок
- WaitGroup
- корректное закрытие result channel
```

Пример концепции:

```go
type Job struct {
    ID    int
    Value int
}

type Result struct {
    JobID int
    Value int
    Err   error
}
```

### Interview notes

Записать короткие ответы:

```text
1. Что такое worker pool?
2. Почему не всегда стоит запускать goroutine на каждую задачу?
3. Что такое bounded concurrency?
4. Кто должен закрывать channel?
```

### Done

```text
[ ] worker pool принимает jobs
[ ] workers останавливаются после закрытия jobs channel
[ ] results channel закрывается корректно
[ ] есть 2–3 unit tests
```

---