# service-demo
Папка configs содержит конфигурационный файл main.yaml, 
где прописан порт и атрибуты для соединения с базами данных PostgreSQL и Redis.

Папка internal содержит пакеты:
  api: RESR API
  config: конфигурация
  database:
    pgsql: пакет для PostgreSQL. До сегодняшнего дня не использовала, но успела 
          реализовать основные методы: Select, Insert, Update.
    redis: пакет для Redis. До сегодняшнего дня не использовала, но успела 
          реализовать основные методы: Set с параметром ttl, Get.
  models: структуры данных и запросы для создания таблиц
  server: запуск сервера
