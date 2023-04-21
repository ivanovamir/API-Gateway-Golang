### Golang API Gateway using pure tools ###

The main file for this application is the data.json file, which is located in the ./data folder. All routes to be handled by the api must be specified in this file.

Let's look at a sample query configuration.

```json
{
  "requests": [
    {
      "path": "/api/hello",
      "url": "https://api/v1/hello_world",
      "method": "GET",
      "make_proxy": true,
      "proxy": "https://auth/api/v1/auth"
    }
  ]
}
```

There is a requests object in the file, which is an array of objects, the element of this array must have the following fields:

- `path` - path of the request to be processed
- `url` - url to which the request is to be made
- `method` - request method
- `make_proxy` - flag indicating whether the request should be proxied
- `proxy` - request path to proxy. If `make_proxy: false` flag value or not specified, there will be validation error.

Список из ближайших целей для проекта:

- [x] Добавить логирование
- [ ] Расширить возможности конфигурации файла `data.json`, возможно, добавить хедеры запросов.
- [ ] Покрыть приложение тестами
- [ ] Оптимизировать код