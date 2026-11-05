package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// fileInfo — описание одного файла: имя и размер.
// Нужен, потому что при печати файла требуется пара «имя + размер»,
// а стандартного типа под эту задачу нет.
type fileInfo struct {
	name string
	size int64
}

// dirInfo — содержимое одной папки: два списка детей.
// dirs хранит имена детей-папок, files — их имена и размеры.
type dirInfo struct {
	dirs  []string
	files []fileInfo
}

// WalkerDir печатает дерево каталога path в out.
// printFiles=false — печатаются только папки.
// prefix — отступ текущего уровня (для корня всегда пустой).
func WalkerDir(out io.Writer, path string, printFiles bool, prefix string) error {
	// Нормализуем путь корня: убираем "./", лишние слэши.
	// filepath.Dir и filepath.Join внутри себя нормализуют пути так же,
	// поэтому ключи карты и путь корня обязаны быть в одном формате —
	// иначе printDir не найдёт корень в карте и ничего не напечатает.
	root := filepath.Clean(path)

	// Карта: путь папки -> её дети. Заполняется на фазе 1, читается на фазе 2.
	children := make(map[string]*dirInfo)

	// Ленивая инициализация: достаём запись о папке, а если её нет — создаём.
	// Храним указатели: map[string]dirInfo отдавал бы КОПИЮ структуры,
	// и append внутри копии не изменил бы то, что лежит в карте.
	get := func(p string) *dirInfo {
		di, ok := children[p]
		if !ok {
			di = &dirInfo{}
			children[p] = di
		}
		return di
	}

	// ФАЗА 1: сбор. Обходим дерево и только заполняем карту — печати нет.
	err := filepath.WalkDir(path, func(internalPath string, d fs.DirEntry, err error) error {
		if err != nil {
			// Ошибка чтения записи (нет прав, файл исчез и т.п.) —
			// возвращаем её, WalkDir прервёт обход, а мы отдадим её наверх.
			return err
		}
		// Корень не вписываем ни в чей список: он никому не ребёнок.
		// Сравниваем с НОРМАЛИЗОВАННЫМ корнем, иначе на "./..." он не совпадёт.
		if internalPath == root {
			return nil
		}

		// Каждый элемент — ребёнок своей родительской папки.
		// Для "project/cmd/server" родитель — "project/cmd".
		parent := filepath.Dir(internalPath)

		if d.IsDir() {
			// WalkDir обходит записи в лексикографическом порядке,
			// поэтому дописывание в конец даёт уже отсортированный список.
			p := get(parent)
			p.dirs = append(p.dirs, d.Name())
			return nil
		}

		// Файл: достаём размер. Файлы собираем ВСЕГДА (независимо от printFiles) —
		// фильтрация произойдёт на фазе печати, там её место.
		info, err := d.Info()
		if err != nil {
			return err
		}
		p := get(parent)
		p.files = append(p.files, fileInfo{name: d.Name(), size: info.Size()})
		return nil
	})
	if err != nil {
		// Сюда попадаем, если корня не существует или это не папка —
		// WalkDir сообщил об этом первой ошибкой в callback.
		return err
	}

	// ФАЗА 2: печать. Рекурсивная функция печатает ДЕТЕЙ папки dirPath с отступом prefix.
	// Объявляем тип заранее, потом присваиваем — иначе внутри тела нельзя
	// вызвать функцию рекурсивно (имя ещё не определено).
	var printDir func(dirPath string, prefix string) error
	printDir = func(dirPath string, prefix string) error {
		// Достаём детей по пути папки. Пустые папки записи в карте не имеют
		// (мы вписывали в карту только родителей, у которых есть дети),
		// поэтому здесь возможен miss — просто нечего печатать внутри.
		di, ok := children[dirPath]
		if !ok {
			return nil
		}

		// ВАЖНО: считаем только тех детей, которые БУДУТ НАПЕЧАТАНЫ.
		// При printFiles=false файлы в подсчёте не участвуют —
		// иначе последней «мнимо» окажется непечатаемая запись,
		// и у настоящей последней папки поедет коннектор.
		n := len(di.dirs)
		if printFiles {
			n += len(di.files)
		}

		// total — порядковый номер текущего элемента среди всех печатаемых детей.
		total := 0

		// --- Папки: идут первыми ---
		for _, name := range di.dirs {
			// Последний элемент — тот, чей номер равен n-1.
			// Сравнение по ОБЩЕМУ счётчику, а не внутри кучки:
			// последний элемент папки может оказаться файлом (.keep в fixtures).
			isLast := total == n-1

			// Коннектор: три чёрточки — так требует эталон Stepik.
			conn := "├───"
			if isLast {
				conn = "└───"
			}

			// Строка папки: отступ + коннектор + имя. Размер у папок не печатается.
			if _, err := fmt.Fprintf(out, "%s%s%s\n", prefix, conn, name); err != nil {
				return err
			}

			// Префикс для поддерева: если ТЕКУЩАЯ папка последняя —
			// на её уровне ниже ничего не будет, ставим табуляцию;
			// если нет — рисуем продолжение вертикальной линии.
			childPrefix := prefix + "\t"
			if !isLast {
				childPrefix = prefix + "│\t"
			}

			// Рекурсия: спускаемся в поддерево ребёнка.
			if err := printDir(filepath.Join(dirPath, name), childPrefix); err != nil {
				return err
			}

			// Учли один элемент. Инкремент после рекурсии — не принципиально,
			// лишь бы ровно один раз на элемент.
			total++
		}

		// --- Файлы: печатаются только с printFiles ---
		if printFiles {
			for _, f := range di.files {
				isLast := total == n-1

				conn := "├───"
				if isLast {
					conn = "└───"
				}

				// Размер: пустой файл — (empty), остальные — (N b) без пробела внутри скобок.
				size := fmt.Sprintf("(%db)", f.size)
				if f.size == 0 {
					size = "(empty)"
				}

				// Строка файла: отступ + коннектор + имя + пробел + размер.
				if _, err := fmt.Fprintf(out, "%s%s%s %s\n", prefix, conn, f.name, size); err != nil {
					return err
				}

				total++
			}
		}

		return nil
	}

	// Запуск: печатаем детей корня с пустым префиксом.
	// Сам корень строкой не выводится — это делает printDir, а не мы.
	return printDir(root, "")
}

// Функция main от вас будет скрыта при проверке, она запускает ту же тестовую структуру,
// на которой решение будет проверяться на Stepik.
func main() {
	root := "./project_tree_test"

	os.RemoveAll(root)

	dirs := []string{
		"cmd/server",
		"internal/auth",
		"internal/db/migrations",
		"pkg/utils",
		"testdata/fixtures/users",
		"testdata/fixtures/orders",
		"docs/api",
		".git/objects",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	files := map[string]string{
		// Корень
		"go.mod":     "module example.com\n\ngo 1.21",
		"readme.md":  "# Project Title\n\nThis is a test project for tree utility.",
		".gitignore": "*.log\n/tmp",

		// cmd
		"cmd/server/main.go": "package main\n\nfunc main() {\n\t// Server start\n}",

		// internal
		"internal/auth/token.go":              "package auth\n\nfunc Generate() string { return \"\" }",
		"internal/auth/token_test.go":         "package auth\n\nfunc TestGenerate(t *testing.T) {}",
		"internal/db/connection.go":           "package db\n\nvar Conn *sql.DB",
		"internal/db/migrations/001_init.sql": "CREATE TABLE users (id INT);",

		// pkg
		"pkg/utils/helpers.go": "package utils\n\nfunc Helper() {}",

		// testdata (смешанные файлы и папки)
		"testdata/fixtures/users/admin.json": `{"role": "admin"}`,
		"testdata/fixtures/users/guest.json": `{"role": "guest"}`,
		"testdata/fixtures/orders/list.csv":  strings.Repeat("id,date,amount\n", 10), // 150b
		"testdata/fixtures/.keep":            "",                                     // пустой файл

		// docs
		"docs/api/swagger.yaml": "openapi: 3.0.0\ninfo:\n  title: API",

		// Пустые файлы в разных местах для проверки (empty)
		"internal/db/migrations/.gitkeep": "",
		".git/objects/pack":               "",
	}

	for path, content := range files {
		fullPath := filepath.Join(root, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Println("Тест 1: Полное дерево (файлы + папки).")
	if err := WalkerDir(os.Stdout, root, true, ""); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println()

	fmt.Println("Тест 2: Только директории.")
	if err := WalkerDir(os.Stdout, root, false, ""); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
