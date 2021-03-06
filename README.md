# Http service

Приложение представляет собой http-сервер с одним хендлером.
Хендлер на вход получает POST-запрос со списком url в json-формате.
Сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате.
Если в процессе обработки хотя бы одного из url получена ошибка, обработка всего списка прекращается и клиенту возвращается текстовая ошибка.

Ограничения:
- использованы можно только компоненты стандартной библиотеки Go (и расширенной стандартной - golang.org/x/*)
- сервер не принимает запрос если количество url в в нем больше 20
- для каждого входящего запроса - не больше 4 одновременных исходящих
- сервер не обслуживает больше чем 100 одновременных входящих http-запросов
- таймаут на запрос одного url - 1 секунда
- обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом
- сервис поддерживает 'graceful shutdown'


Запуск:
```
make run
```

Прогон тестов:
```
make test
```

Параметры, которые нужны для изменения конфигурации сервиса [здесь](.env.example)
