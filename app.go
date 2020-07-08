package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// структура для хранения конфигураций
type Config struct {
	DefaultFile string `json:"default_file"`
}

// структура товара
type Good struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Producer string `json:"producer"`
	Price    int    `json:"price"`
	Count    int    `json:"count"`
}

// точка входа в программу
func main() {
	// инициализация основных переменных
	scanner := bufio.NewScanner(os.Stdin)
	args := os.Args
	argLength := len(os.Args)
	config := readConfigs()
	goods := readGoods(config.DefaultFile)

	if argLength > 1 {
		command := args[1]

		// реализация команды вывода информации по всем товарам
		if strings.Compare(command, "list") == 0 {
			fmt.Println("----------------------------------------------------------------------")
			fmt.Printf("|%-4s|%-20s|%-20s|%-8s|%-12s|\n", " ID", "      Название", "   Производитель", "  Цена", " Количество ")
			for i := 0; i < len(goods); i++ {
				fmt.Println("+----+--------------------+--------------------+--------+------------+")
				fmt.Printf("|%-4d|%-20s|%-20s|%8d|%12d|\n", goods[i].Id, goods[i].Name, goods[i].Producer, goods[i].Price, goods[i].Count)
			}
			fmt.Println("----------------------------------------------------------------------")
		}

		// реализация команды установки рабочего файла
		if strings.Compare(command, "-f") == 0 {
			config.DefaultFile = args[2]
			saveConfigs(config)
		}

		// реализация команды добавления новых данных
		if strings.Compare(command, "new") == 0 {
			good := Good{}
			good.Id = nextID(goods)

			fmt.Print("Название товара: ")
			scanner.Scan()
			good.Name = scanner.Text()

			fmt.Print("Производитель: ")
			scanner.Scan()
			good.Producer = scanner.Text()

			for {
				fmt.Print("Цена: ")
				scanner.Scan()
				num, err := strconv.Atoi(scanner.Text())
				good.Price = num
				if err == nil {
					break
				}
			}

			for {
				fmt.Print("Кол-во товаров: ")
				scanner.Scan()
				num, err := strconv.Atoi(scanner.Text())
				good.Count = num
				if err == nil {
					break
				}
			}

			goods = append(goods, good)
			saveGoods(goods, config.DefaultFile)
			fmt.Println("appended good: ", good)
		}

		// реализация команды редактирования
		if strings.Compare(command, "edit") == 0 {

			index := findIndexByID(goods, args[2])
			fmt.Println("Редактировать поле:\n1. Название\n2. Производитель\n3. Цена\n4. Количество")
			var choose int
			_, _ = fmt.Scanf("%d\n", &choose)
			if choose == 1 {
				fmt.Print("Название товара: ")
				scanner.Scan()
				goods[index].Name = scanner.Text()
			} else if choose == 2 {
				fmt.Print("Производитель: ")
				scanner.Scan()
				goods[index].Producer = scanner.Text()
			} else if choose == 3 {
				for {
					fmt.Print("Цена: ")
					scanner.Scan()
					num, err := strconv.Atoi(scanner.Text())
					goods[index].Price = num
					if err == nil {
						break
					}
				}
			} else if choose == 4 {
				for {
					fmt.Print("Кол-во товаров: ")
					scanner.Scan()
					num, err := strconv.Atoi(scanner.Text())
					goods[index].Count = num
					if err == nil {
						break
					}
				}
			}
			saveGoods(goods, config.DefaultFile)
		}

		// реализация команды удаления
		if strings.Compare(command, "del") == 0 {
			index := findIndexByID(goods, args[2])
			if index != -1 {
				goods = append(goods[:index], goods[index+1:]...)
				saveGoods(goods, config.DefaultFile)
			} else {
				fmt.Println("Неверно указанный ID")
			}
		}

		// реализация команда вывода указанного товара
		if strings.Compare(command, "read") == 0 {
			index := findIndexByID(goods, args[2])
			if index != -1 {
				good := goods[index]
				fmt.Println("----------------------------------------------------------------------")
				fmt.Printf("|%-4s|%-20s|%-20s|%-8s|%-12s|\n", " ID", "      Название", "   Производитель", "  Цена", " Количество ")
				fmt.Println("----------------------------------------------------------------------")
				fmt.Printf("|%-4d|%-20s|%-20s|%8d|%12d|\n", good.Id, good.Name, good.Producer, good.Price, good.Count)
				fmt.Println("----------------------------------------------------------------------")
			} else {
				fmt.Println("Неверно указанный ID")
			}
		}
	}
}

// функция сохранения конфигураций
func saveConfigs(config *Config) {
	goodsJson, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Cannot encode to JSON", err)
	}

	err = ioutil.WriteFile("config.json", goodsJson, 0777)
	if err != nil {
		log.Fatal("Cannot write data to file", err)
	}
}

// функция чтения конфигураций
func readConfigs() *Config {
	bytes, _ := ioutil.ReadFile("config.json")
	config := Config{}
	_ = json.Unmarshal(bytes, &config)
	return &config
}

// функция сохранения товаров
func saveGoods(goods []Good, filename string) {
	goodsJson, err := json.Marshal(goods)
	if err != nil {
		log.Fatal("Cannot encode to JSON", err)
	}

	err = ioutil.WriteFile(filename, goodsJson, 0777)
	if err != nil {
		log.Fatal("Cannot write data to file", err)
	}
}

// функция чтения товаров
func readGoods(filename string) []Good {
	bytes, _ := ioutil.ReadFile(filename)
	var products []Good
	_ = json.Unmarshal(bytes, &products)
	return products
}

// функция генерации нового (следующего) айдишника
func nextID(goods []Good) int {
	maxID := -1
	for i := 0; i < len(goods); i++ {
		if maxID < goods[i].Id {
			maxID = goods[i].Id
		}
	}
	maxID++
	return maxID
}

// функция, возвращающая индекс товара в массиве по указанному айдишнику
func findIndexByID(goods []Good, id string) int {
	id_, _ := strconv.Atoi(id)
	for i := 0; i < len(goods); i++ {
		if goods[i].Id == id_ {
			return i
		}
	}
	return -1
}
