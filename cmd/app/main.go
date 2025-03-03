package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type InvertedLine struct {
	LineNumber int    `json:"line_number"`
	Inverted   string `json:"inverted"`
}

type Repository struct {
	FilePath string
}

func NewRepository(filePath string) *Repository {
	return &Repository{FilePath: filePath}
}

func (r *Repository) SaveData(data string) error {
	file, err := os.OpenFile(r.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	return err
}

func (r *Repository) ReadData() ([]string, error) {
	file, err := os.ReadFile(r.FilePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines, nil
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveData(data string) error {
	return s.repo.SaveData(data)
}

func (s *Service) GetInvertedLines() ([]InvertedLine, error) {
	lines, err := s.repo.ReadData()
	if err != nil {
		return nil, err
	}

	var invertedLines []InvertedLine
	for i, line := range lines {
		invertedLines = append(invertedLines, InvertedLine{
			LineNumber: i + 1,
			Inverted:   reverseString(line),
		})
	}

	return invertedLines, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SaveDataHandler(c echo.Context) error {
	var input struct {
		Data string `json:"data"`
	}

	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Неверный формат JSON")
	}

	if input.Data == "" {
		return c.String(http.StatusBadRequest, "Поле 'data' не может быть пустым")
	}

	if err := h.service.SaveData(input.Data); err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при записи в файл")
	}

	return c.String(http.StatusOK, "Данные успешно записаны в файл")
}

func (h *Handler) GetInvertedLinesHandler(c echo.Context) error {
	invertedLines, err := h.service.GetInvertedLines()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Ошибка при чтении файла")
	}

	return c.JSON(http.StatusOK, invertedLines)
}

func main() {
	repo := NewRepository("data.txt")
	service := NewService(repo)
	handler := NewHandler(service)

	e := echo.New()

	api := e.Group("/api")
	{
		api.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})

		api.POST("/input", handler.SaveDataHandler)
		api.GET("/input", handler.GetInvertedLinesHandler)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
