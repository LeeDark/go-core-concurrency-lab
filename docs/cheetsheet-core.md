# Slices

## Definition

Slice — это динамическое представление части массива. Он содержит pointer на backing array, длину `len` и емкость `cap`.

- Slice не хранит данные сам, он ссылается на array.
- `append` может переиспользовать array или создать новый.
- Изменение элемента slice меняет backing array.
- `len` — текущая длина, `cap` — доступная емкость.
- Subslice может удерживать большой backing array в памяти.
- `nil slice` и empty slice ведут себя похоже, но не идентичны.
- Slice не thread-safe.
- Для копирования используется `copy`.

## Internal model

- https://go.dev/blog/slices-intro

## len/cap

```text
slice — это не “массив переменной длины”.
Это маленький header, который держит ссылку на массив, длину и capacity.
Отсюда почти все эти чудеса, из-за которых взрослые люди сидят на собеседованиях и обсуждают cap.
```

## append

```text
append добавляет элементы в существующий backing array, если хватает capacity.
Если capacity не хватает, Go выделяет новый backing array и копирует туда старые элементы.
Поэтому несколько slices могут неожиданно видеть изменения друг друга, если они разделяют один backing array.
Для защиты можно использовать copy или full slice expression s[low:high:max].
```

## copy

```text
copy копирует элементы в уже существующий destination slice и не меняет его длину.
append добавляет элементы и может использовать старый backing array или создать новый.
Поэтому append может неожиданно изменить другие slices, если они разделяют backing array.
Для безопасного копирования используют make + copy или append([]T(nil), s...),
но второй вариант может превратить empty slice в nil slice.
```

## nil vs empty

```go
var nilSlice []int
emptySlice := []int{}
emptySliceMake := make([]int, 0)
```

## aliasing

```go
a := []int{1, 2, 3, 4}
// subslice aliasing bug
b := a[:2]
b = append(b, 99)

c := []int{1, 2, 3, 4}
// full slice expression
d := c[:2:2] // len = 2, cap = 2
d = append(d, 99)
```

## common mistakes
## interview questions
