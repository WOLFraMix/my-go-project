package main

import (
	"fmt"
)

type Contact struct {
	Name  string
	Phone string
}

func main() {
	m := make(map[string]Contact)

	var n int
	for {
		fmt.Print(`1 - Добавить
2 - Найти
3 - Все контакты
0 - Выход
Выбор: `)
		fmt.Scan(&n)

		name := ""
		phone := ""
		switch n {
		case 1:
			fmt.Print("Имя: ")
			fmt.Scan(&name)
			fmt.Print("Телефон: ")
			fmt.Scan(&phone)
			m[name] = Contact{
				Name:  name,
				Phone: phone,
			}
			fmt.Println()
			continue
		case 2:
			fmt.Print("Имя: ")
			fmt.Scan(&name)
			contact, ok := m[name]
			if !ok {
				fmt.Println("err contact does not exist")
				fmt.Println()
				continue
			}
			fmt.Printf("Телефон: %s\n", contact.Phone)
			fmt.Println()
			continue
		case 3:
			for _, v := range m {
				fmt.Printf("%s: %s\n", v.Name, v.Phone)
				fmt.Println()
			}
			continue
		case 0:
			return
		}
	}
}
