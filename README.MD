Задание на поцизицию Junior Golang

**Используемые технологии:**

- Go
- JWT
- MongoDB

**Задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя сидентификатором (GUID) указанным в параметре
  запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refreshтокенов

**Требования:**

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базеисключительно в виде bcrypt хеша, должен быть
защищен от изменения настороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном
который был выдан вместе с ним.

**Результат:**

Результат выполнения задания нужно предоставить в виде исходного кода на Github.

## Мои дороботки:

- Было решено изменить немного задание, запрос получения токенов пускай принимает логин и пароль, а не UID пользователя,
  хотя первоначально я сделал чисто через GUID в url как параметр
- Был добавлен маршрут чтобы зарегистрировать клиента (email является уникальным)
- Для тестирования еще был добавлен метод блокировки рефреш токена, например когда токен был скомпрометирован и
  необходимо прекратить его дейсвтие

### Маршруты
- Authentication
  - /api/v1/login
    - ``{"password": "","email": ""}``
  - /api/v1/refresh
    - ``{"access_token": "","refresh_token": ""}``
  
  - /api/v1/block_token
    - ``{"password": "","email": ""}``
- User
  - /api/v1/register
    - ``{ "password": "","email": "","username": ""}``