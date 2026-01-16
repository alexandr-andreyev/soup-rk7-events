1. Get Order Status From Notify Sent, ready, taken
2. Get Order Cancel status (как понять что заказ отменен)
3. Log Get Events and send events with format (переделать формат сообщения для отправки) Sent = Cooking, ready = Ready, taken = IssuedOut
4. Добавить подпись для WSA Events
4. Send Event to url in list
5. Send event with network error
6. Send event after reboot



Алгоритм проверки статусов и отправки событий во внешние системы

1. Как сделать несколько подписчиков для отправки событий (вид подписки, тип событий, фильтры для событий)
2. При получении события проверять подписку и обрабатывать согласно настройкам подписки

Отправка статусов заказа
1. При включение получить список заказов согласно фильтру (проверять список заказов раз в 5 минут)
2. Если вернулись заказы то записать их в бд
3. При получении события OrderChanged проверить заказ в бд и статус, если статус в событии отличается, то сохранить в бд новый статус и отправить событие по подписке
4. Статусы заказа Created, Cooking, Ready, Issued Out (должны идти именно в таком порядке), и статус Canceled (отмена заказа может происходить при любом статусе)

Автообновление модуля:
1. При включении проверять версию приложения, если есть версия новее, то скачивать в фоне и обновлять

Настройки из центра:
1. При включении пробовать загрузить настройки из бэкенда (если в файле конфигурации есть такая настройка)


Основные проблемы

События могут приходить не по порядку (сетевые задержки)
Внешняя система может быть недоступна (нужны retry)
Дубликаты событий (at-least-once delivery)
Сбои в процессе обработки (потеря событий)

Рекомендуемый подход: Outbox Pattern + State Machine

// Состояние отправки для каждого подписчика отдельно
type SubscriberOrderState struct {
    SubscriberID   string
    OrderID        string
    LastSentStatus Status
    UpdatedAt      time.Time
}

// Outbox сообщение привязано к подписчику
type OutboxMessage struct {
    ID            int64
    SubscriberID  string    // ключевое поле!
    OrderID       string
    Status        Status
    Payload       []byte
    CreatedAt     time.Time
    ProcessedAt   *time.Time
    RetryCount    int
    NextRetryAt   time.Time
}

// Цепочка статусов (state machine)
type Status string

const (
    StatusCreated    Status = "created"
    StatusConfirmed  Status = "confirmed"
    StatusPreparing  Status = "preparing"
    StatusReady      Status = "ready"
    StatusDelivering Status = "delivering"
    StatusCompleted  Status = "completed"
    StatusCancelled  Status = "cancelled"
)

// Порядок статусов для валидации цепочки
var statusOrder = map[Status]int{
    StatusCreated:    1,
    StatusConfirmed:  2,
    StatusPreparing:  3,
    StatusReady:      4,
    StatusDelivering: 5,
    StatusCompleted:  6,
}

// Какие статусы можно пропустить при отправке
// (например, если заказ уже в "ready", а мы отправили только "created")
func getStatusChain(from, to Status) []Status {
    chain := []Status{}
    fromOrder := statusOrder[from]
    toOrder := statusOrder[to]
    
    for status, order := range statusOrder {
        if order > fromOrder && order <= toOrder {
            chain = append(chain, status)
        }
    }
    // сортировка по порядку
    sort.Slice(chain, func(i, j int) bool {
        return statusOrder[chain[i]] < statusOrder[chain[j]]
    })
    return chain
}
