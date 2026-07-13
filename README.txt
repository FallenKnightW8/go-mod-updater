# Go Module Updater CLI

## Описание
Утилита командной строки для анализа Go-модулей в Git-репозиториях. 
Выводит имя модуля, версию Go и список зависимостей с доступными обновлениями.

## Требования
1. Установленный Go (версия 1.21 или выше)
2. Установленный Git
3. Доступ в интернет для клонирования репозиториев

## Установка

### Шаг 1: Сборка исполняемого файла
```cmd
go build -o go-mod-updater.exe
Шаг 2: Проверка работы
cmd
go-mod-updater.exe --help
Настройка списка репозиториев
Шаг 1: Создайте файл repos.txt
Создайте текстовый файл repos.txt в папке с программой.

Шаг 2: Добавьте репозитории
Укажите репозитории построчно в одном из форматов:

gin-gonic/gin
gorilla/mux
https://github.com/labstack/echo
git@github.com:gofiber/fiber.git
Пример содержимого repos.txt:
gin-gonic/gin
gorilla/mux
labstack/echo
gofiber/fiber
kubernetes/kubernetes
Важно:

Один репозиторий на строку

Можно использовать короткий формат user/repo для GitHub

Можно указать полный URL для любых Git-серверов

Пустые строки игнорируются

Использование
analyze.bat отккрывает cmd меню
Выбор репозитория
Введите номер репозитория из списка
Нажмите Enter

Выбор формата вывода
После выбора репозитория доступны опции:

Нажать Enter - показать таблицу с зависимостями на экране

Ввести 1 - сохранить результаты в JSON-файл в папке results\

После анализа нажмите любую клавишу для возврата в меню

Выход из программы
В главном меню введите 0 для выхода

Форматы вывода
Табличный формат (по умолчанию)
Показывает:

Имя модуля

Версия Go

Список зависимостей

Текущая версия

Доступная версия

Статус (Up to date / Update available)

Тип зависимости (direct / indirect)

JSON формат
Сохраняется в папку results\ с именем файла:

text
results/имя_репозитория_ДД-ММ-ГГ_ЧЧ-ММ.json
Примеры использования
Пример 1: Просмотр зависимостей gin
text
1. Запустить analyze.bat
2. Выбрать номер репозитория gin-gonic/gin
3. Нажать Enter для табличного вывода
4. Просмотреть результат
5. Нажать любую клавишу для возврата в меню
Пример 2: Сохранение в JSON
text
1. Запустить analyze.bat
2. Выбрать номер репозитория gin-gonic/gin
3. Ввести 1 для сохранения в JSON
4. Файл сохранен в results/gin-gonic_gin_13-07-26_14-30.json
5. Нажать любую клавишу для возврата в меню
Структура JSON-файла
json
{
  "module_name": "github.com/gin-gonic/gin",
  "go_version": "1.21",
  "dependencies": [
    {
      "name": "github.com/gin-contrib/sse",
      "current_version": "v0.1.0",
      "latest_version": "v0.2.0",
      "is_direct": true,
      "can_update": true
    }
  ]
}
Возможные ошибки и решения
Ошибка: "repos.txt not found"
Причина: Файл repos.txt отсутствует
Решение: Программа создаст пример файла автоматически. Отредактируйте его.

Ошибка: "git clone failed"
Причина:

Не установлен Git

Нет доступа в интернет

Неверный URL репозитория
Решение: Проверьте подключение к интернету и правильность URL

Ошибка: "go.mod not found"
Причина: Репозиторий не является Go-проектом
Решение: Убедитесь, что анализируете Go-проект с файлом go.mod

Ошибка: "go list failed"
Причина: Проблемы с Go-окружением
Решение: Проверьте установку Go командой go version

Дополнительные команды
Прямой анализ без меню
cmd
go-mod-updater.exe --repo https://github.com/gin-gonic/gin
Только прямые зависимости
cmd
go-mod-updater.exe --repo https://github.com/gin-gonic/gin --direct
JSON вывод в консоль
cmd
go-mod-updater.exe --repo https://github.com/gin-gonic/gin --json
Удаление
cmd
del go-mod-updater.exe
del analyze.bat
del repos.txt
rmdir /s /q results
Техническая информация
Язык: Go 1.21

Буфер клонирования: 1 коммит (--depth 1)

Временные файлы: удаляются автоматически

Кодировка: UTF-8

Зависимости: github.com/spf13/cobra

text

## Инструкция по созданию файлов

1. **Сохраните `analyze.bat`** в `S:\go-mod-updater\analyze.bat`
2. **Сохраните `README.md`** в `S:\go-mod-updater\README.md`
3. **Создайте `repos.txt`** в `S:\go-mod-updater\repos.txt`