package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Chronos struct {
	DeviceID int     `json:"device_id"`
	Time     string  `json:"timestamp"`
	Value    float64 `json:"value"`
}

// ProcessLogs выполняет потоковую фильтрацию и конвертацию логов.
// Входной файл содержит строки в формате:
// JSON Lines с метками времени "YYYY/MM/DD-hh:mm:ss±hhmm".
// Функция должна отфильтровать записи,
// попадающие в диапазон [start, end] включительно.
// Результат записывается в outputPath,
// где поле timestamp становится строкой RFC3339 в таймзоне UTC.
func ProcessLogs(inputPath, outputPath string, start, end time.Time) error {
	// Эталон
	const layout = "2006/01/02-15:04:05-0700"

	// Открываем файл
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	// Закрываем файл при выходе из функции
	defer file.Close()

	// Проверяем, что файл не пустой
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}
	if stat.Size() == 0 {
		return nil
	}

	// Начинаем вывод
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	// Закрываем вывод при выходе из функции
	defer out.Close()

	// Пишем вывод
	w := bufio.NewWriter(out)
	defer w.Flush()

	// Создаём декодер
	decoder := json.NewDecoder(file)
	for {
		var c Chronos
		// Декодируем
		if err := decoder.Decode(&c); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Парсим время
		t, err := time.Parse(layout, c.Time)
		if err != nil {
			return err
		}

		// Если НЕ ДО старта и НЕ ПОСЛЕ конца
		if !t.Before(start) && !t.After(end) {
			c.Time = t.UTC().Format(time.RFC3339)
			b, err := json.Marshal(c)
			if err != nil {
				return err
			}
			_, err = w.Write(b)
			if err != nil {
				return err
			}
			err = w.WriteByte('\n')
			if err != nil {
				return err
			}
		}
	}
	return w.Flush()
}

// Функция main и тесты от вас будут скрыты во время проверки.
func main() {
	if !test1() || !test2() || !test3() || !test4() || !test5() || !test6() || !test7() || !test8() || !test9() {
		os.Exit(1)
	}
	fmt.Println("Все тесты успешно пройдены!")
}

// Тест 1: Базовый вариант
func test1() bool {
	inPath := "test1_in.jsonl"
	outPath := "test1_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	content := `{"device_id": 101, "timestamp": "2026/09/06-11:30:00+0300", "value": 15.4}
{"device_id": 102, "timestamp": "2026/09/06-10:00:00+0300", "value": 12.1}
{"device_id": 103, "timestamp": "2026/09/06-13:00:00+0300", "value": 18.2}
`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	loc := time.FixedZone("UTC+3", 3*60*60)
	start := time.Date(2026, 9, 6, 11, 0, 0, 0, loc)
	end := time.Date(2026, 9, 6, 12, 0, 0, 0, loc)

	if err := ProcessLogs(inPath, outPath, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 1 упал с ошибкой: %s\n", err)
		return false
	}

	res, err := os.ReadFile(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 1: Не удалось прочитать файл с результатами: %s\n", err)
		return false
	}

	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	if len(lines) != 1 || lines[0] == "" {
		fmt.Fprintf(os.Stderr, "Тест 1: Ожидалась ровно 1 строка, получено: %d\n", len(lines))
		return false
	}

	// Проверяем, что устройство 101 сохранено, а время переведено в UTC (11:30 +0300 -> 08:30 Z)
	if !strings.Contains(lines[0], `"device_id":101`) || !strings.Contains(lines[0], `"2026-09-06T08:30:00Z"`) {
		fmt.Fprintf(os.Stderr, "Тест 1: Некорректный формат или данные на выходе. Получено: %s\n", lines[0])
		return false
	}
	return true
}

// Тест 2: Граничные значения, проверяем включение start и end ровно секунда в секунду
func test2() bool {
	inPath := "test2_in.jsonl"
	outPath := "test2_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	content := `{"device_id": 201, "timestamp": "2026/09/06-11:00:00+0300", "value": 1.0}
{"device_id": 202, "timestamp": "2026/09/06-12:00:00+0300", "value": 2.0}
`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	loc := time.FixedZone("UTC+3", 3*60*60)
	start := time.Date(2026, 9, 6, 11, 0, 0, 0, loc)
	end := time.Date(2026, 9, 6, 12, 0, 0, 0, loc)

	if err := ProcessLogs(inPath, outPath, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 2 упал с ошибкой: %s\n", err)
		return false
	}

	res, _ := os.ReadFile(outPath)
	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	if len(lines) != 2 || lines[0] == "" {
		fmt.Fprintf(os.Stderr, "Тест 2: Границы должны быть включены. Ожидалось 2 строки, получено: %d\n", len(lines))
		return false
	}
	return true
}

// Тест 3: Пустой файл, при этом код должен отработать вхолостую и не упасть
func test3() bool {
	inPath := "test3_in.jsonl"
	outPath := "test3_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	if err := os.WriteFile(inPath, []byte(""), 0644); err != nil {
		return false
	}

	if err := ProcessLogs(inPath, outPath, time.Now(), time.Now()); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 3: Пустой файл вызвал ошибку: %s\n", err)
		return false
	}
	return true
}

// Тест 4: Ошибки ввода-вывода, передаем фейковый путь, ждем ошибку открытия файла
func test4() bool {
	err := ProcessLogs("non_existent_file_12345.jsonl", "out.jsonl", time.Now(), time.Now())
	if err == nil {
		fmt.Fprintf(os.Stderr, "Тест 4: Ожидалась ошибка открытия несуществующего файла\n")
		return false
	}
	return true
}

// Тест 5: Кривой формат данных, парсинг должен спотыкнуться на невалидном формате времени
func test5() bool {
	inPath := "test5_in.jsonl"
	outPath := "test5_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	content := `{"device_id": 501, "timestamp": "broken-time-format", "value": 9.9}`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	if err := ProcessLogs(inPath, outPath, time.Now(), time.Now()); err == nil {
		fmt.Fprint(os.Stderr, "Тест 5: Ожидалась ошибка парсинга некорректного формата времени\n")
		return false
	}
	return true
}

// Тест 6: Синтаксическая ошибка в JSON
func test6() bool {
	inPath := "test6_in.jsonl"
	outPath := "test6_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	content := `{"device_id": 601, "timestamp": "2026/09/06-11:30:00+0300", "value": 15.4` // Нет закрывающей скобки
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	if err := ProcessLogs(inPath, outPath, time.Now(), time.Now()); err == nil {
		fmt.Fprint(os.Stderr, "Тест 6: Ожидалась ошибка из-за некорректного синтаксиса JSON\n")
		return false
	}
	return true
}

// Тест 7: Смешанные часовые пояса во входных данных (положительные и отрицательные смещения)
func test7() bool {
	inPath := "test7_in.jsonl"
	outPath := "test7_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	// В UTC это:
	// 1) 11:30:00 +0300 -> 08:30:00 UTC (Должен пройти)
	// 2) 03:30:00 -0500 -> 08:30:00 UTC (Должен пройти)
	// 3) 15:30:00 +0900 -> 06:30:00 UTC (Вне диапазона)
	content := `{"device_id": 701, "timestamp": "2026/09/06-11:30:00+0300", "value": 10.0}
{"device_id": 702, "timestamp": "2026/09/06-03:30:00-0500", "value": 20.0}
{"device_id": 703, "timestamp": "2026/09/06-15:30:00+0900", "value": 30.0}
`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	// Ищем строго записи, соответствующие 08:30:00 по UTC
	start := time.Date(2026, 9, 6, 8, 30, 0, 0, time.UTC)
	end := time.Date(2026, 9, 6, 8, 30, 0, 0, time.UTC)

	if err := ProcessLogs(inPath, outPath, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 7 упал с ошибкой: %s\n", err)
		return false
	}

	res, err := os.ReadFile(outPath)
	if err != nil {
		return false
	}

	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	// Должно остаться ровно 2 записи (701 и 702)
	if len(lines) != 2 || lines[0] == "" {
		fmt.Fprintf(os.Stderr, "Тест 7: Ожидалось 2 строки, получено: %d\n", len(lines))
		return false
	}

	if !strings.Contains(string(res), `"device_id":701`) || !strings.Contains(string(res), `"device_id":702`) {
		fmt.Fprintf(os.Stderr, "Тест 7: В выходной файл попали не те устройства. Получено:\n%s\n", string(res))
		return false
	}
	return true
}

// Тест 8: Отрицательное смещение в логах и фильтр в совершенно другой таймзоне
func test8() bool {
	inPath := "test8_in.jsonl"
	outPath := "test8_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	// 07:00:00 -0400 (например, Нью-Йорк EDT) -> это 11:00:00 UTC
	content := `{"device_id": 801, "timestamp": "2026/09/06-07:00:00-0400", "value": 5.5}`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	// UTC+5 это Екатеринбург/Ташкент
	// 15:00:00 +0500 -> 10:00:00 UTC
	// 17:00:00 +0500 -> 12:00:00 UTC
	// Наша запись (11:00:00 UTC) аккурат между ними.
	loc := time.FixedZone("UTC+5", 5*60*60)
	start := time.Date(2026, 9, 6, 15, 0, 0, 0, loc)
	end := time.Date(2026, 9, 6, 17, 0, 0, 0, loc)

	if err := ProcessLogs(inPath, outPath, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 8 упал с ошибкой: %s\n", err)
		return false
	}

	res, _ := os.ReadFile(outPath)
	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	if len(lines) != 1 || lines[0] == "" {
		fmt.Fprintf(os.Stderr, "Тест 8: Запись должна была пройти фильтр. Получено строк: %d\n", len(lines))
		return false
	}

	// Проверяем формат RFC3339 и перевод в UTC на выходе
	if !strings.Contains(lines[0], `"2026-09-06T11:00:00Z"`) {
		fmt.Fprintf(os.Stderr, "Тест 8: Неверный формат времени UTC на выходе. Получено: %s\n", lines[0])
		return false
	}
	return true
}

// Тест 9: Валидный файл, но ни одна запись не подходит под диапазон
func test9() bool {
	inPath := "test9_in.jsonl"
	outPath := "test9_out.jsonl"
	defer os.Remove(inPath)
	defer os.Remove(outPath)

	content := `{"device_id": 901, "timestamp": "2026/09/06-11:30:00+0300", "value": 42.0}`
	if err := os.WriteFile(inPath, []byte(content), 0644); err != nil {
		return false
	}

	// Фильтр указывает на далекое прошлое
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	if err := ProcessLogs(inPath, outPath, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "Тест 9 упал с ошибкой при пустом результате: %s\n", err)
		return false
	}

	res, _ := os.ReadFile(outPath)
	// Файл должен быть пустым
	if len(strings.TrimSpace(string(res))) != 0 {
		fmt.Fprintf(os.Stderr, "Тест 9: Ожидался пустой выходной файл, но получено: %s\n", string(res))
		return false
	}
	return true
}
