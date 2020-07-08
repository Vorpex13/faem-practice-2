package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DefaultFile string `json:"default_file"`
}

type Good struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Producer string `json:"producer"`
	Price    int    `json:"price"`
	Count    int    `json:"count"`
}

func main() {
	args := os.Args
	argLength := len(os.Args)
	config := readConfigs()
	goods := readGoods(config.DefaultFile)

	if argLength > 1 {
		command := args[1]

		if strings.Compare(command, "list") == 0 {
			fmt.Println("-----------------------------------------------------------------")
			fmt.Printf("|%-20s|%-20s|%-8s|%-12s|\n", "      Название", "   Производитель", "  Цена", " Количество ")

			for i := 0; i < len(goods); i++ {
				fmt.Println("+--------------------+--------------------+--------+------------+")
				fmt.Printf("|%-20s|%-20s|%8d|%12d|\n", goods[i].Name, goods[i].Producer, goods[i].Price, goods[i].Count)
			}
			fmt.Println("-----------------------------------------------------------------")
		}

		if strings.Compare(command, "-f") == 0 {
			config.DefaultFile = args[2]
			saveConfigs(config)
		}

		if strings.Compare(command, "new") == 0 {
			good := Good{}
			good.Id = nextID(goods)
			fmt.Print("Название товара: ")
			_, _ = fmt.Scanf("%s\n", &good.Name)
			fmt.Print("Производитель: ")
			_, _ = fmt.Scanf("%s\n", &good.Producer)
			fmt.Print("Кол-во товаров: ")
			_, _ = fmt.Scanf("%d\n", &good.Count)
			fmt.Print("Цена товара: ")
			_, _ = fmt.Scanf("%d\n", &good.Price)

			goods = append(goods, good)
			saveGoods(goods, config.DefaultFile)
			fmt.Println("appended good: ", good)
		}

		if strings.Compare(command, "edit") == 0 {
			fmt.Println("edit")
		}

		if strings.Compare(command, "del") == 0 {
			index := findIndexByID(goods, args[2])
			if index != -1 {
				goods = append(goods[:index], goods[index+1:]...)
				saveGoods(goods, config.DefaultFile)
			} else {
				fmt.Println("Неверно указанный ID")
			}
		}

		if strings.Compare(command, "read") == 0 {
			index := findIndexByID(goods, args[2])
			if index != -1 {
				good := goods[index]
				fmt.Println("-----------------------------------------------------------------")
				fmt.Printf("|%-20s|%-20s|%-8s|%-12s|\n", "      Название", "   Производитель", "  Цена", " Количество ")
				fmt.Println("+--------------------+--------------------+--------+------------+")
				fmt.Printf("|%-20s|%-20s|%8d|%12d|\n", good.Name, good.Producer, good.Price, good.Count)
				fmt.Println("-----------------------------------------------------------------")
			} else {
				fmt.Println("Неверно указанный ID")
			}
		}
	}
}

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

func readConfigs() *Config {
	bytes, _ := ioutil.ReadFile("config.json")
	config := Config{}
	_ = json.Unmarshal(bytes, &config)
	return &config
}

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

func readGoods(filename string) []Good {
	bytes, _ := ioutil.ReadFile(filename)
	var products []Good
	_ = json.Unmarshal(bytes, &products)
	return products
}

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

func findIndexByID(goods []Good, id string) int {
	id_, _ := strconv.Atoi(id)
	for i := 0; i < len(goods); i++ {
		if goods[i].Id == id_ {
			return i
		}
	}
	return -1
}
