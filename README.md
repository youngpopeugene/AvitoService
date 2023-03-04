# Тестовое задание на позицию стажера-бекендера

## Микросервис для работы с балансом пользователей.

### Решенные задачи
Реализованы следующие методы:
* Метод получения текущего баланса пользователя.
* Метод начисления и списания средств.
* Метод резервирования средств с основного баланса на отдельном счете. 
* Метод признания выручки.
* Развернут docker compose

## Как запустить ?
```
git clone https://github.com/youngpopeugene/GolangTask.git

cd GolangTask

docker-compose build && docker-compose up
```

### Метод получения текущего баланса пользователя
Требуется отправить GET запрос на localhost:8181/get_balance?user_id={ число } с указанием id пользователя

#### Особенности
* Принимает id пользователя. (Значение указано в рублях, все операции подразумевают работу с рублем, что является положительным числом, за исключением операции списания средств, она подразумевает указание отрицательного числа).

#### Пример 1
>localhost:8181/get_balance?user_id=1

RESPONSE:
```json
{
  "type": "success",
  "data": [
    {
      "user_id": 1,
      "user_balance": 0
    }
  ],
  "message": ""
}
```
#### Пример 2
>localhost:8181/get_balance?user_id=2

RESPONSE:
```json
{
  "type": "error",
  "data": null,
  "message": "User with id=2 wasn't found"
}
```

### Метод начисления и списания средств.
Требуется отправить POST запрос на localhost:8181/update_balance
#### Особенности
* Принимает id пользователя и сколько средств зачислить/списать. 
* Если пользователь с указанным id существует его баланс обновится в соответсвии с параметрами. Если нет - пользователь будет добавлен в базу данных. 

#### Пример 1
>localhost:8181/update_balance

BODY:
```json : 
{
  "user_id": 1,
  "user_balance": 200
}
```
RESPONSE:
```json
{
  "type": "success",
  "data": [
    {
      "user_id": 1,
      "user_balance": 200
    }
  ],
  "message": "User's balance was updated"
}
```
#### Пример 2
>localhost:8181/update_balance

BODY:
```json : 
{
  "user_id": 2,
  "user_balance": 300
}
```
RESPONSE:
```json
{
  "type": "success",
  "data": [
    {
      "user_id": 2,
      "user_balance": 300
    }
  ],
  "message": "User was created"
}
```

### Метод резервирования средств с основного баланса на отдельном счете.
Требуется отправить POST запрос на localhost:8181/from_user_to_reserve. 

#### Особенности
* Принимает id пользователя, с которого нужно списать средства, id услуги, id заказа, и стоимость. 
* Метод делает проверку на существование пользователя с указаным id, на наличие у него достаточного количество средств, чтобы провести резервирование, на уникальность id услуги и id заказа (не может существовать две резервации с одинаковыми id услуги и id заказа).

#### Пример 1
> localhost:8181/from_user_to_reserve

BODY:
```json :
{
    "user_id": 1,
    "service_id": 1,
    "order_id": 1,
    "price": 50
}
```
RESPONSE:
```json
{
    "type": "success",
    "data": null,
    "message": "Some money from balance of user with id=1 was transferred to reserve, new reserve created"
}
```

#### Пример 2
> localhost:8181/from_user_to_reserve

BODY:
```json :
{
    "user_id": 2,
    "service_id": 1,
    "order_id": 1,
    "price": 50
}
```
RESPONSE:
```json
{
  "type": "error",
  "data": null,
  "message": "Cannot create new reserve - this combination of order_id and service_id already exist"
}
```

#### Пример 3
> localhost:8181/from_user_to_reserve

BODY:
```json :
{
    "user_id": 5,
    "service_id": 1,
    "order_id": 3,
    "price": 50
}
```
RESPONSE:
```json
{
    "type": "error",
    "data": null,
    "message": "There is no such user_id in 'users' table"
}
```

### Метод признания выручки.
Требуется отправить POST запрос на localhost:8181/from_reserve_to_user.

#### Особенности
* Принимает id пользователя, которому нужно начислить выручку, id услуги, id заказа, и сумму. 
* Метод делает проверку на существование резерва в базе данных по уникальной связке id услуги + id заказа. 
* Возможна ситуация, когда запрашиваемая сумма меньше, чем стоимость, указанная в резерве, тогда оставшиеся деньги вернутся пользователю, который "создал" этот резерв. Если же запрашиваемая сумма больше стоимости, указанной в резерве, тогда операция отклоняется. 
* Если пользователя, который собирается получить прибыль, не существует в базе данных, то он будет добавлен с балансом равным прибыли, в ином случае его баланс просто обновится.

#### Пример 1
> localhost:8181/from_reserve_to_user

BODY:
```json :
{
    "user_id": 5,
    "service_id": 1,
    "order_id": 1,
    "price": 50
}
```
RESPONSE:
```json
{
  "type": "success",
  "data": null,
  "message": "User with id=5 was created with some money on balance, reserve was deleted"
}
```

#### Пример 2
> localhost:8181/from_reserve_to_user

BODY:
```json :
{
    "user_id": 1,
    "service_id": 1,
    "order_id": 2,
    "price": 20
}
```
RESPONSE:
```json
{
  "type": "success",
  "data": null,
  "message": "User with id=1 increased his balance, user with id=2 got some money back, reserve was deleted"
}
```

#### Пример 3
> localhost:8181/from_reserve_to_user

BODY:
```json :
{
    "user_id": 5,
    "service_id": 3,
    "order_id": 3,
    "price": 50
}
```
RESPONSE:
```json
{
  "type": "error",
  "data": null,
  "message": "No reserves with these service_id and order_id"
}
```
