package helpers

import (
	"fmt"
	"regexp"
	"sync/atomic"
)

// Counter структура для увеличения значения счетчика
type Counter struct {
	value int64
}

// NewCounter создает новый экземпляр счетчика
func NewCounter(initialValue int64) *Counter {
	return &Counter{value: initialValue}
}

// Increment увеличивает значение счетчика на 1 и возвращает новое значение
func (c *Counter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// IsValidUUID проверяет валидность UUID
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// FormatID форматирует число в строку с заданным префиксом и нулевым заполнением
func FormatID(prefix string, value int64) string {
	return fmt.Sprintf("%s%06d", prefix, value)
}

// OrderIDCounter глобальный счетчик для Order ID
var OrderIDCounter = NewCounter(1)

// GetNextOrderID возвращает следующее значение Order ID
func GetNextOrderID() string {
	nextValue := OrderIDCounter.Increment()
	return FormatID("O-", nextValue)
}

// ProductIDCounter глобальный счетчик для Product ID
var ProductIDCounter = NewCounter(1)

// GetNextProductID возвращает следующее значение Product ID
func GetNextProductID() string {
	nextValue := ProductIDCounter.Increment()
	return FormatID("P-", nextValue)
}
