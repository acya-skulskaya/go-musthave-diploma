# go-musthave-diploma-tpl

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.

_________


# Lint
orig: https://github.com/s-shpak/praktikum-golangci-lint?tab=readme-ov-file

# Mocks
1. Install mockery https://vektra.github.io/mockery/latest/installation/
2. Generate config `$ mockery init github.com/vektra/mockery/v3/internal/fixtures`
3. Add packages to generate mocks, example

```yaml 
  packages:
    github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance:
      interfaces:
        Balance:
```
   
4. Generate mocks `$ mockery`
