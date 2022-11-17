Описание

Микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.

Подход

Для достижения атомарности выполняемых операций и изолированности данных использовались транзакции. Для согласованности данных и надежности, использую брокер сообщений Kafka, который гарантирует, что мы точно обработаем каждое сообщение ровно один раз, даже при падениях отдельных частей системы.  

Запуск проекта
1. Go to "deployments" folder
2. docker-compose up
3. Wait until all containers will be up and give them some time to initialize

БД

Таблицы и все необходимые связи создаются при первом запуске из файла init.sql

Структура:
accounts(id, balance, account_type) - используется для работы с балансом счета
orders(id, service_id, price, user_account_id, status, creation_date) - хранит информациию о статусах обработок транзакций заказов.
services(id, service_name, account_id) - тип оказываемой услуги. Каждая услуга имеет свой счет для резервации.
transactions(id, type, amount, sender_id, receiver_id, creation_date, comment) - история всех типов операций со счетами.

API

В корне проекта лежит POSTMAN коллекция с желательной последовательностью запросов и подготовленными телами

1. Метод начисления средств на баланс:
- Method: POST 
- Path: /api/v1/account/deposit
- Request body:

    account_id - (uint) идентификатор счета

    amount - (uint) сумма пополнения
- Description: при отсутствии счета с указанным идентификатором, создается и пополняется новый
- Postman request name: Deposit, Deposit 2
2. Метод резервирования средств с основного баланса на отдельном счете:  
- Method: POST
- Path: /api/v1/account/reserve
- Request body:

  user_account_id - (uint) идентификатор счета

  service_id - (uint) идентификатор услуги

  order_id - (uint) идентификатор заказа

  price - (float) стоимость
- Description: переводит зарезервированную сумму на отдельный счет, создает запись с запонанием статуса обработки заказа, добавляет перевод в историю.
- Postman request name: Reserve
3. Метод признания выручки
- Method: async
- Request body:

  user_account_id - (uint) идентификатор счета

  service_id - (uint) идентификатор услуги

  order_id - (uint) идентификатор заказа

  price - (float) стоимость

  status - (string) статус исполнения заказа от стороннего сервиса
- Description: асинхронное взаимодействие. Подписываемся с помощью Kafka на изменение статуса обработки заказа от сервиса deliver-service(имитируем работу в docker-compose). Если приходит заказ со статусом "COMPLETED" - переводим статус заказа в "COMPLETED", зарезервированные средства переводим на счет с прибылью компании, сохраняем перевод в историю транзакций со статусом "PROFIT". В случае любого другого статуса заказа - переводим статус заказа в "CANCELLED", зарезервированные средства возвращаем на счет пользователя, сохраняем перевод в историю транзакций со статусом "REFUND". Для имитации ответа необходимо выполнить запрос "4"
4. Метод для отправки статуса заказа в очередь:
- Method: POST
- Path: http://localhost:9091/api/v1/orders/process
- Request body:

  user_account_id - (uint) идентификатор счета

  service_id - (uint) идентификатор услуги

  order_id - (uint) идентификатор заказа

  price - (float) стоимость

  status - ("COMPLETED"/other) статус исполнения заказа от стороннего сервиса
- Description: отправка в очередь. Для успешной обработки заказ нужно status = "COMPETED". Все остальные статусы будут обрабатываться, но считаться неуспешными.
- Postman request name: Push order status, Push failed order sttatus
5. Метод получения баланса пользователя:
- Method: GET
- Path: /api/v1/accounts/{account_id}/balance
- Query:
    
  account_id - (uint) идентификатор счета
- Description: переводит зарезервированную сумму на отдельный счет, создает запись с запонанием статуса обработки заказа, добавляет перевод в историю.
- Postman request name: Reserve
6. Метод для получения месячного отчета:
- Method: GET
- Path: /api/v1/report/service/profit
- Query:

  year - (uint) год

  month - (uint) месяц
- Description: Создает отчет по прибыли в разрезе всех категорий услуг по всем пользователям за выбранный период. Файл сохраняется в каталоге reports/ в корне проекта. На запрос статус, а не ссылка. Это сделано для удоства, чтобы не передавать через env в image абсолютный путь.
- Postman request name: Generate CSV
7. Метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса:
- Method: GET
- Path: /api/v1/accounts/1/transactions
- Query:

  date_sort - (asc/desc/нету значения) сортировка по дате

  amount_sort - (asc/desc/нету значения) сортировка по сумме

  page - (uint) - порядковый номер желаемой выборки
- Description: размера элементов в одной странице - 9. Сортировка по дате и сумме - возрастающее/убывающее/нету. Значения выбираются из transactions для указанного счета
- Postman request name: Get transactions with default pagination, Get transactions with date sort, Get transactions with amount sort
8. Перевод средств:
- Method: POST
- Path: /api/v1/account/transfer
- Request body:

  from - (uint) от кого

  amount - (float) сумма перевода

  to - (uint) кому
- Description: списывает деньги со счета одного, прибавляет другому. Добавляет перевод в историю транзакций со типом "TRANSFER". Счет аккаунта получателя должен быть создан заранее депозитом в любую сумму.
- Postman request name: Transfer
9. Вывод средств:
- Method: POST
- Path: /api/v1/account/withdraw
- Request body:

  from - (uint) кто снимает

  amount - (float) сумма перевода
- Description: списывает деньги с указанного счета. Добавляет перевод в историю транзакций со типом "WITHDRAWAL"
- Postman request name: Withdrawal






