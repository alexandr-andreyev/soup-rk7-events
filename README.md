# Soup RK7 Events - Windows Service

Сервис получает и обрабатывает события от кассового сервера rkeeper через интерфейс notify.

Задача сервиса отправлять события на внешние ресурсы в требуемом формате. Например изменение стоп листа или изменение заказа на кассе.

## Как начать

Для запуска есть следующие параметры

* debug - запускает в debug режиме
* install - устанавливает как службу windows
* remove - удаляет службу
* start
* stop
* pause
* continue

## Запуск сервиса

After compiling an executable, the service can be installed from an Administrative command prompt.  Typing

    YourExecutable.EXE install 

will install the service.

To update the service, stop the service, replace the executble and restart the service.

The service can be removed from an Administrative command prompt by typing:

    YourExecutable.EXE remove 


<Prop name="{79D8D0EC-1DE1-4F21-9105-388384215A90}" value="190cc89af297b8315b9"/> - originalOrderId






