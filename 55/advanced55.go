package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SensitiveData string

func (sd SensitiveData) MarshalJSON() ([]byte, error) {
	return []byte(`"***"`), nil
}

func (sd *SensitiveData) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*sd = SensitiveData(s + "_processed")
	return nil
}

type LogRecord struct {
	Component string `json:"component"`
	Level     int    `json:"level"`
	SensitiveData
}

// ----------------------------------------------------------------------
// Данная строка и все что выше будет скрыто от вас при проверке на сайте

func (lr LogRecord) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(struct {
		Component     string        `json:"component"`
		Level         int           `json:"level"`
		SensitiveData SensitiveData `json:"sensitiveData"`
	}{
		Component:     lr.Component,
		Level:         lr.Level,
		SensitiveData: lr.SensitiveData,
	})
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (lr *LogRecord) UnmarshalJSON(data []byte) error {
	var raw struct {
		Component     *string        `json:"component"`
		Level         *int           `json:"level"`
		SensitiveData *SensitiveData `json:"sensitiveData"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw.Component != nil {
		lr.Component = *raw.Component
	}
	if raw.Level != nil {
		lr.Level = *raw.Level
	}
	if raw.SensitiveData != nil {
		lr.SensitiveData = *raw.SensitiveData
	}
	return nil
}

// Данная строка и все что ниже будет скрыто от вас при проверке на сайте
// ----------------------------------------------------------------------
func main() {
	if !test1() || !test2() || !test3() || !test4() || !test5() {
		os.Exit(1)
	}
	fmt.Println("Все тесты успешно пройдены!")
}

// Маршалинг, маскируем секрет, сохраняем остальные поля
func test1() bool {
	rec := LogRecord{Component: "Auth", Level: 1, SensitiveData: "super_secret"}
	data, err := json.Marshal(rec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка маршалинга: %v\n", err)
		return false
	}

	expected := `{"component":"Auth","level":1,"sensitiveData":"***"}`
	if string(data) != expected {
		fmt.Fprintf(os.Stderr, "Неверный JSON. Ожидали: %s, получили: %s\n", expected, string(data))
		return false
	}
	return true
}

// Анмаршалинг, проверяем секрет и заполнение полей структуры
func test2() bool {
	input := []byte(`{"component":"Worker","level":42,"sensitiveData":"raw_token"}`)
	var rec LogRecord
	if err := json.Unmarshal(input, &rec); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка анмаршалинга: %v\n", err)
		return false
	}

	if rec.Component != "Worker" || rec.Level != 42 || rec.SensitiveData != "raw_token_processed" {
		fmt.Fprintf(os.Stderr, "Данные побились при анмаршалинге. Получили: Component=%s, Level=%d, Data=%s\n", rec.Component, rec.Level, rec.SensitiveData)
		return false
	}
	return true
}

// Пустые значения и нулевые уровни, на всякий случай проверяем
func test3() bool {
	rec := LogRecord{}
	data, err := json.Marshal(rec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка маршалинга пустой структуры: %v\n", err)
		return false
	}

	expected := `{"component":"","level":0,"sensitiveData":"***"}`
	if string(data) != expected {
		fmt.Fprintf(os.Stderr, "Неверный JSON для пустых полей. Получили: %s\n", string(data))
		return false
	}
	return true
}

// Смотрим что будет при неверном синтаксисе JSON, ожидаем ошибку от парсера
func test4() bool {
	badJSON := []byte(`{"component":"API", "level": "строка_вместо_инта"}`)
	var rec LogRecord
	if err := json.Unmarshal(badJSON, &rec); err == nil {
		fmt.Fprintln(os.Stderr, "Ожидали ошибку валидации типов JSON, но парсер проглотил некорректные данные")
		return false
	}
	return true
}

// Частичный анмаршалинг, проверяем что пропущенные в JSON поля не затирают текущие значения структуры
func test5() bool {
	rec := LogRecord{Component: "KeepMe", Level: 99, SensitiveData: "Initial"}
	input := []byte(`{"level":10}`)
	if err := json.Unmarshal(input, &rec); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при частичном анмаршалинге: %v\n", err)
		return false
	}

	if rec.Component != "KeepMe" {
		fmt.Fprintf(os.Stderr, "Ошибка: поле Component было затерто дефолтным значением при частичном апдейте\n")
		return false
	}
	if rec.Level != 10 {
		fmt.Fprintf(os.Stderr, "Ошибка: поле Level не обновилось. Ожидали 10, получили %d\n", rec.Level)
		return false
	}
	if rec.SensitiveData != "Initial" {
		fmt.Fprintf(os.Stderr, "Ошибка: поле SensitiveData изменилось, хотя отсутствовало в JSON. Получили: %s\n", rec.SensitiveData)
		return false
	}
	return true
}
