package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var ErrExpired = errors.New("file has expired")

// ExpiringFile реализует io.ReadCloser
// с ограничением по времени работы после первого чтения.
type ExpiringFile struct {
	File      os.File
	Duration  time.Duration
	IsOpen    bool
	Start     time.Time
	FirstRead bool
}

// NewExpiringFile открывает файл и инициализирует структуру.
// Если файл не существует или не может быть открыт, возвращает ошибку.
func NewExpiringFile(filePath string, duration time.Duration) (*ExpiringFile, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("err open file")
	}

	return &ExpiringFile{
		File:      *file,
		Duration:  duration,
		IsOpen:    true,
		FirstRead: false,
	}, nil
}

// Read считывает данные.
// При первом вызове фиксирует время начала.
// Если лимит времени превышен - возвращает ошибку ErrExpired.
func (ef *ExpiringFile) Read(p []byte) (int, error) {
	if !ef.IsOpen {
		return 0, io.ErrClosedPipe
	}
	if !ef.FirstRead {
		ef.Start = time.Now()
		ef.FirstRead = true
	}
	if time.Since(ef.Start) > ef.Duration {
		return 0, ErrExpired
	}

	n, err := ef.File.Read(p)
	if err == io.EOF {
		return 0, io.EOF
	}
	if err != nil {
		return 0, errors.New("err read file")
	}
	return n, nil
}

// Close закрывает дескриптор файла.
func (ef *ExpiringFile) Close() error {
	if !ef.IsOpen {
		return nil
	}
	err := ef.File.Close()
	if err != nil {
		return errors.New("err close file")
	}
	ef.IsOpen = false
	return nil
}

// Функция main будет скрыта от вас при проверке на сайте
func main() {
	if !test1() || !test2() || !test3() || !test4() || !test5() || !test6() || !test7() || !test8() {
		os.Exit(1)
	}
	fmt.Println("Все тесты успешно пройдены!")
}

// Вспомогательная функция, создает временный файл.
// Возвращает путь к файлу и функцию для его удаления.
func createTempFile(content string) (string, func()) {
	tmp, err := os.CreateTemp("", "expire_test_*.txt")
	if err != nil {
		panic(err)
	}
	_, _ = tmp.WriteString(content)
	_ = tmp.Close()
	return tmp.Name(), func() { os.Remove(tmp.Name()) }
}

// Базовый случай, успешное чтение в рамках лимита
func test1() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()

	ef, err := NewExpiringFile(path, 1*time.Hour)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 1 провален: не удалось открыть файл: %v\n", err)
		return false
	}
	defer ef.Close()

	buf := make([]byte, 5)
	n, err := ef.Read(buf)
	if err != nil || n != 5 || string(buf[:n]) != "01234" {
		fmt.Fprintf(os.Stderr, "Тест 1 провален: неверное чтение. n=%d, err=%v, data=%s\n", n, err, string(buf[:n]))
		return false
	}
	return true
}

// Проверка блокировки по истечении времени
func test2() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()

	ef, err := NewExpiringFile(path, 50*time.Millisecond)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 2 провален: не удалось открыть файл: %v\n", err)
		return false
	}
	defer ef.Close()

	buf := make([]byte, 2)
	_, err = ef.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Тест 2 провален на первом чтении: %v\n", err)
		return false
	}

	time.Sleep(300 * time.Millisecond)

	_, err = ef.Read(buf)
	if !errors.Is(err, ErrExpired) {
		fmt.Fprintf(os.Stderr, "Тест 2 провален: ожидалась ошибка ErrExpired, получено %v\n", err)
		return false
	}
	return true
}

// Чтение из закрытого файла
func test3() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()

	ef, err := NewExpiringFile(path, 1*time.Hour)
	if err != nil {
		return false
	}
	_ = ef.Close()

	buf := make([]byte, 5)
	_, err = ef.Read(buf)
	if !errors.Is(err, io.ErrClosedPipe) {
		fmt.Fprintln(os.Stderr, "Тест 3 провален: ожидалась ошибка io.ErrClosedPipe при чтении из закрытого файла")
		return false
	}
	return true
}

// Несуществующий файл
func test4() bool {
	_, err := NewExpiringFile("non_existent_file_xyz_777.txt", 1*time.Hour)
	if !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "Тест 4 провален: ожидалась ошибка при попытке открыть несуществующий файл")
		return false
	}
	return true
}

// Идемпотентность метода Close (безопасное повторное закрытие)
func test5() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()

	ef, err := NewExpiringFile(path, 1*time.Hour)
	if err != nil {
		return false
	}

	err1 := ef.Close()
	err2 := ef.Close()
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "Тест 5 провален: ошибка при повторном Close: %v, %v\n", err1, err2)
		return false
	}
	return true
}

// Проверка ленивого старта, таймаут не должен тикать с момента создания структуры
func test6() bool {
	path, cleanup := createTempFile("abcdef")
	defer cleanup()

	// Срок жизни всего 50 мс
	ef, err := NewExpiringFile(path, 50*time.Millisecond)
	if err != nil {
		return false
	}
	defer ef.Close()

	// Ждем 500 мс до первого чтения. Если таймер запустился в конструкторе - чтение упадет.
	time.Sleep(500 * time.Millisecond)

	buf := make([]byte, 3)
	n, err := ef.Read(buf)
	if err != nil || n != 3 || string(buf[:n]) != "abc" {
		fmt.Fprintf(os.Stderr, "Тест 6 провален: время пошло до первого вызова Read. Ошибка: %v, прочитано: %d\n", err, n)
		return false
	}
	return true
}

func test7() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()
	ef, err := NewExpiringFile(path, 1*time.Hour)
	if err != nil {
		return false
	}
	defer ef.Close()
	buf := make([]byte, 100) // буфер больше содержимого файла
	n, err := ef.Read(buf)
	// По контракту io.Reader на последнем чтении допустимы оба варианта: (10, io.EOF) или (10, nil)
	if n != 10 || (err != nil && !errors.Is(err, io.EOF)) {
		fmt.Fprintf(os.Stderr, "Тест 7 провален: ожидалось (10, io.EOF или nil) - такой контракт io.Reader, получено (%d, %v)\n", n, err)
		return false
	}
	// А вот следующий вызов обязан вернуть ровно (0, io.EOF)
	n2, err2 := ef.Read(buf)
	if n2 != 0 || !errors.Is(err2, io.EOF) {
		fmt.Fprintf(os.Stderr, "Тест 7 провален: следующий Read после конца файла должен вернуть (0, io.EOF), получено (%d, %v)\n", n2, err2)
		return false
	}
	return true
}

func test8() bool {
	path, cleanup := createTempFile("0123456789")
	defer cleanup()
	ef, err := NewExpiringFile(path, 10*time.Millisecond)
	if err != nil {
		return false
	}
	defer ef.Close()
	buf := make([]byte, 2)
	if _, err = ef.Read(buf); err != nil {
		return false
	}
	time.Sleep(50 * time.Millisecond)
	if _, err = ef.Read(buf); !errors.Is(err, ErrExpired) {
		fmt.Fprintf(os.Stderr, "Тест 8 провален: ожидалась ErrExpired, получено %v\n", err)
		return false
	}
	time.Sleep(50 * time.Millisecond)
	if _, err = ef.Read(buf); !errors.Is(err, ErrExpired) {
		fmt.Fprintf(os.Stderr, "Тест 8 провален: повторное чтение после истечения должно снова давать ErrExpired, получено %v\n", err)
		return false
	}
	return true
}
