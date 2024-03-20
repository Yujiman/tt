package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Замените URL на вашу ссылку
	url := "https://s3.eu-central-1.amazonaws.com/stage-cdn.halyk-travel.com/media/translations/frontend/kk.json"

	// Отправляем GET-запрос по указанной ссылке
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer response.Body.Close()

	// Читаем тело ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Создаем структуру данных для распаковки JSON
	var data map[string]interface{}

	// Распаковываем JSON в структуру данных
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Ошибка при распаковке JSON:", err)
		return
	}

	// Теперь вы можете работать с данными, представленными в переменной "data"
	fmt.Println(data)
}
