# Slices, Maps, Defer

## Slices

Slice — это динамическое представление части массива. Он содержит pointer на underlying array, длину `len` и емкость `cap`.

- Slice не хранит данные сам, он ссылается на array.
- `append` может переиспользовать array или создать новый.
- Изменение элемента slice меняет underlying array.
- `len` — текущая длина, `cap` — доступная емкость.
- Subslice может удерживать большой underlying array в памяти.
- `nil slice` и empty slice ведут себя похоже, но не идентичны.
- Slice не thread-safe.
- Для копирования используется `copy`.

### Read

- https://go.dev/blog/slices-intro

```text
append добавляет элементы в существующий backing array, если хватает capacity.
Если capacity не хватает, Go выделяет новый backing array и копирует туда старые элементы.
Поэтому несколько slices могут неожиданно видеть изменения друг друга, если они разделяют один backing array.
Для защиты можно использовать copy или full slice expression s[low:high:max].
```

### Labs

### Cheetsheet, Interview

