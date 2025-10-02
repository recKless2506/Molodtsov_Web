package repository

import (
	"fmt"
	"strings"
)

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Product struct {
	ID          int
	Title       string
	ImageURL    string
	Price       string
	Description string
	Specs       string
}

type ZayavkaDefaults struct {
	Area        string
	TempOutside string
	TempInside  string
	CarrierVol  string
	Result      string
}

func (r *Repository) GetProducts() ([]Product, error) {
	products := []Product{
		{ID: 1, Title: "Электрический котёл", ImageURL: "http://127.0.0.1:9000/webdesign/teplodar_sputnik_elektro_6_belyj.jpg", Price: "Цена за 1 кВт·ч = 7,87 руб.", Description: "Мощный электрический котёл для отопления дома.", Specs: "Мощность: 2000 Вт, Напряжение: 220 В"},
		{ID: 2, Title: "Газовый котёл", ImageURL: "http://127.0.0.1:9000/webdesign/xMINI-NIKE-24-002-260x260.jpg.pagespeed.ic.xblKe_TacA.jpg", Price: "Цена за 1 м³ газа = 8,54 руб.", Description: "Современный газовый котёл с высоким КПД.", Specs: "Диапазон мощности: 5-25 кВт"},
		{ID: 3, Title: "Электрический конвектор", ImageURL: "http://127.0.0.1:9000/webdesign/P2312111420054464853_zoom.webp", Price: "Цена за 1 кВт·ч = 7,87 руб.", Description: "Компактный электрический конвектор для дома.", Specs: "Мощность: 1500 Вт, Регулировка температуры"},
		{ID: 4, Title: "Чугунный радиатор", ImageURL: "http://127.0.0.1:9000/webdesign/standaard4h.jpg", Price: "Цена за 1 кВт·ч = 7,87 руб.", Description: "Надёжный чугунный радиатор для длительной эксплуатации.", Specs: "Количество секций: 8, Рабочее давление: 10 бар"},
		{ID: 5, Title: "Водяной тепловентилятор", ImageURL: "http://127.0.0.1:9000/webdesign/P2311130852045782136_zoom.webp", Price: "Цена за 1 кВт·ч = 7,87 руб., Цена за 1 м³ газа = 8,54 руб.", Description: "Многофункциональный водяной тепловентилятор.", Specs: "Расход воды: 2 л/мин, Мощность: 2000 Вт"},
		{ID: 6, Title: "Инфракрасная панель", ImageURL: "http://127.0.0.1:9000/webdesign/503264-rrecht_2.webp", Price: "Цена за 1 кВт·ч = 7,87 руб.", Description: "Компактная инфракрасная панель для точечного отопления.", Specs: "Мощность: 800 Вт, Напряжение: 220 В"},
	}
	return products, nil
}

func (r *Repository) GetProduct(id int) (Product, error) {
	products, err := r.GetProducts()
	if err != nil {
		return Product{}, err
	}

	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}
	return Product{}, fmt.Errorf("товар не найден")
}

func (r *Repository) GetProductsByTitle(title string) ([]Product, error) {
	products, err := r.GetProducts()
	if err != nil {
		return []Product{}, err
	}

	var result []Product
	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Title), strings.ToLower(title)) {
			result = append(result, product)
		}
	}
	return result, nil
}

func (r *Repository) GetZayavkaDefaults() map[int]ZayavkaDefaults {
	return map[int]ZayavkaDefaults{
		1: {Area: "100", TempOutside: "-15", TempInside: "22", CarrierVol: "50", Result: "25000 руб."},
		2: {Area: "80", TempOutside: "-10", TempInside: "20", CarrierVol: "40", Result: "18000 руб."},
	}
}
