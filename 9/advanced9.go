package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Translator хранит загруженные переводы.
// Подсказка: подумайте над использованием map[string]map[string]string
type Translator struct {
	Blocks  map[string]map[string]string
	Default string
}

// NewTranslator читает файл и инициализирует структуру Translator.
// Возвращает ошибку, если файл не удалось прочитать.
func NewTranslator(filePath string, defaultLang string) (*Translator, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	t := Translator{
		Blocks:  make(map[string]map[string]string),
		Default: defaultLang,
	}
	var block string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, "@") {
			continue
		}
		if strings.HasPrefix(line, "@") {
			if line[1:] == "" {
				continue
			}
			block = line[1:]
			continue
		}
		if block == "" {
			continue
		}
		if t.Blocks[block] == nil {
			t.Blocks[block] = make(map[string]string)
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		t.Blocks[block][key] = val
	}
	return &t, nil
}

// Get возвращает перевод для заданного ключа и языка.
// Реализуйте логику fallback на defaultLang, если язык не найден.
// Если и defaultLang отсутствует - вернуть сам key.
func (t *Translator) Get(key string, lang string) string {
	for block, m := range t.Blocks {
		if block == key {
			for l, v := range m {
				if l == lang {
					return v
				}
			}
			for l, v := range m {
				if l == t.Default {
					return v
				}
			}
		}
	}
	return key
}

// Функция main и все тесты ниже будут скрыты от вас на сайте
func main() {
	// Запуск всех тестов. Если хоть один упадет - выходим с кодом 1
	if !test1() || !test2() || !test3() || !test4() || !test5() ||
		!test6() || !test7() || !test8() || !test9() || !test10() ||
		!test11() || !test12() || !test13() || !test14() || !test15() ||
		!test16() || !test17() {
		os.Exit(1)
	}
	fmt.Println("Все тесты успешно пройдены!")
}

// writeTemp - хелпер для тестов. Записывает текст во временный файл и возвращает путь к нему.
func writeTemp(content string) string {
	f, err := os.CreateTemp("", "locales-*.txt")
	if err != nil {
		return ""
	}
	_, _ = f.WriteString(content)
	name := f.Name()
	_ = f.Close()
	return name
}

// Обычный перевод для существующего языка
func test1() bool {
	file := writeTemp("@greeting\nen=Hello\nru=Привет\n")
	defer os.Remove(file)

	t, err := NewTranslator(file, "en")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 1: не удалось создать Translator: %v\n", err)
		return false
	}

	if res := t.Get("greeting", "ru"); res != "Привет" {
		fmt.Fprintf(os.Stderr, "Тест 1 (ru): ожидалось 'Привет', получили '%s'\n", res)
		return false
	}
	if res := t.Get("greeting", "en"); res != "Hello" {
		fmt.Fprintf(os.Stderr, "Тест 1 (en): ожидалось 'Hello', получили '%s'\n", res)
		return false
	}
	return true
}

// Откат на defaultLang, если запрошенного языка нет
func test2() bool {
	file := writeTemp("@greeting\nen=Hello\nru=Привет\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("greeting", "fr"); res != "Hello" {
		fmt.Fprintf(os.Stderr, "Тест 2: для языка 'fr' ожидался дефолтный перевод 'Hello', получили '%s'\n", res)
		return false
	}
	return true
}

// Несуществующий ключ, возвращаем сам key
func test3() bool {
	file := writeTemp("@greeting\nen=Hello\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("unknown", "en"); res != "unknown" {
		fmt.Fprintf(os.Stderr, "Тест 3: для несуществующего ключа ожидался возврат самого ключа 'unknown', получили '%s'\n", res)
		return false
	}
	return true
}

// Ключ есть, но ни запрошенного языка, ни defaultLang нет, должны вернуть key
func test4() bool {
	file := writeTemp("@only_ru\nru=Только русский\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("only_ru", "en"); res != "only_ru" {
		fmt.Fprintf(os.Stderr, "Тест 4: когда нет ни перевода, ни дефолтного языка в ключе, должен вернуться сам ключ 'only_ru', получили '%s'\n", res)
		return false
	}
	return true
}

// Комментарии и пустые строки должны игнорироваться
func test5() bool {
	file := writeTemp("# comment\n\n@greeting\nen=Hello\n# another comment\nru=Привет\n\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("greeting", "ru"); res != "Привет" {
		fmt.Fprintf(os.Stderr, "Тест 5: парсер не проигнорировал комментарии/пустые строки. На 'ru' ожидали 'Привет', получили '%s'\n", res)
		return false
	}
	return true
}

// Файла нет, ожидаем ошибку
func test6() bool {
	_, err := NewTranslator("/nonexistent/path/that/does/not/exist.txt", "en")
	if err == nil {
		fmt.Fprintln(os.Stderr, "Тест 6: ожидалась ошибка при попытке открыть несуществующий файл, но ошибка не вернулась")
		return false
	}
	return true
}

// Пустой файл не падает, любой запрос возвращает key
func test7() bool {
	file := writeTemp("")
	defer os.Remove(file)

	t, err := NewTranslator(file, "en")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 7: пустой файл вызвал ошибку инициализации: %v\n", err)
		return false
	}
	if res := t.Get("anything", "en"); res != "anything" {
		fmt.Fprintf(os.Stderr, "Тест 7: для пустого файла любой запрос должен возвращать сам ключ 'anything', получили '%s'\n", res)
		return false
	}
	return true
}

// Несколько "=" в значении, срезаем только по первому
func test8() bool {
	file := writeTemp("@formula\nen=a=b=c\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("formula", "en"); res != "a=b=c" {
		fmt.Fprintf(os.Stderr, "Тест 8: строка со множественными '=' распарсилась неверно. Ожидалось 'a=b=c', получили '%s'\n", res)
		return false
	}
	return true
}

// UTF-8 - кириллица, эмодзи, японский
func test9() bool {
	file := writeTemp("@emoji\nen=Hello 🫥\nru=Привет 🫪\njp=こんにちは 🫩\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("emoji", "ru"); res != "Привет 🫪" {
		fmt.Fprintf(os.Stderr, "Тест 9 (ru): проблемы с UTF-8. Ожидалось 'Привет 🫪', получили '%s'\n", res)
		return false
	}
	if res := t.Get("emoji", "jp"); res != "こんにちは 🫩" {
		fmt.Fprintf(os.Stderr, "Тест 9 (jp): проблемы с UTF-8. Ожидалось 'こんにちは 🫩', получили '%s'\n", res)
		return false
	}
	return true
}

// Дубликаты ключей, побеждает последнее значение
func test10() bool {
	file := writeTemp("@greeting\nen=Hello\n@greeting\nen=Hi\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("greeting", "en"); res != "Hi" {
		fmt.Fprintf(os.Stderr, "Тест 10: при дублировании ключей/языков должно побеждать последнее значение. Ожидалось 'Hi', получили '%s'\n", res)
		return false
	}
	return true
}

// Пустое значение (lang=) - это валидный перевод, просто пустая строка
func test11() bool {
	file := writeTemp("@empty\nen=\nru=Привет\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("empty", "en"); res != "" {
		fmt.Fprintf(os.Stderr, "Тест 11: конструкция 'lang=' должна считываться как валидная пустая строка. Ожидалось '', получили '%s'\n", res)
		return false
	}
	return true
}

// Кривые строки без "=", просто пропускаем
func test12() bool {
	file := writeTemp("garbage\n@greeting\nen=Hello\nno-equals-sign\nru=Привет\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("greeting", "ru"); res != "Привет" {
		fmt.Fprintf(os.Stderr, "Тест 12: строки без знака '=' ломают парсинг. Ожидалось 'Привет', получили '%s'\n", res)
		return false
	}
	return true
}

// Пробелы в начале строк, TrimSpace должен спасти
func test13() bool {
	file := writeTemp("  @greeting\n  en=Hello\n  ru=Привет\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("greeting", "ru"); res != "Привет" {
		fmt.Fprintf(os.Stderr, "Тест 13: пробелы по краям строк не были обрезаны. Ожидалось 'Привет', получили '%s'\n", res)
		return false
	}
	return true
}

// Мусор до первого "@", игнорируем
func test14() bool {
	// Добавляем знак "=", чтобы спровоцировать запись в мапу у неверных реализаций
	file := writeTemp("invalid_lang=some garbage before any @\n@greeting\nen=Hello\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")

	// Проверим что последующий код работает
	if res := t.Get("greeting", "en"); res != "Hello" {
		fmt.Fprintf(os.Stderr, "Тест 14: мусор в начале файла до первого тега '@' сломал логику. Ожидалось 'Hello', получили '%s'\n", res)
		return false
	}

	// Мусора не должно быть, речь про перевод для пустого ключа
	if res := t.Get("", "invalid_lang"); res != "" {
		fmt.Fprintf(os.Stderr, "Тест 14: мусор до первого '@' со знаком '=' был ошибочно сохранен в память. Ожидался пустой ответ '', получили '%s'\n", res)
		return false
	}
	return true
}

// "@" без имени, игнорируем, следующие строки не должны никуда попасть
func test15() bool {
	file := writeTemp("@\nen=Должны проигнорить\n@Снимаю шляпу\nen=Преклоняюсь :)\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")

	// Проверяем что последующий код работает
	if res := t.Get("Снимаю шляпу", "en"); res != "Преклоняюсь :)" {
		fmt.Fprintf(os.Stderr, "Тест 15: одиночный символ '@' без имени ключа должен игнорироваться вместе со связанными переводами. Ожидалось 'Hello', получили '%s'\n", res)
		return false
	}

	// Тоже, как в предыдущем тесте, данные после пустого '@' не должны быть привязаны к пустой строке
	if res := t.Get("", "en"); res != "" {
		fmt.Fprintf(os.Stderr, "Тест 15: строки после пустого '@' были ошибочно сохранены для пустого ключа. Ожидался пустой ответ '', получили '%s'\n", res)
		return false
	}
	return true
}

// Запрос существующего языка, когда defaultLang для этого ключа отсутствует
func test16() bool {
	file := writeTemp("@only_ru\nru=Только русский\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("only_ru", "ru"); res != "Только русский" {
		fmt.Fprintf(os.Stderr, "Тест 16: запросили ru, а дефолтного en нет. Ожидался 'Только русский', получили '%s'\n", res)
		return false
	}
	return true
}

// Разные ключи не должны смешиваться в одну мапу
func test17() bool {
	file := writeTemp("@first\nru=Первый\n@second\nru=Второй\n")
	defer os.Remove(file)

	t, _ := NewTranslator(file, "en")
	if res := t.Get("first", "ru"); res != "Первый" {
		fmt.Fprintf(os.Stderr, "Тест 16: переводы смешались. Для first/ru ожидался 'Первый', получили '%s'\n", res)
		return false
	}
	return true
}
