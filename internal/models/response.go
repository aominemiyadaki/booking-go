package models

import "log/slog"

// ErrorResponse — стандартный ответ об ошибке.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse — стандартный успешный ответ.
type SuccessResponse struct {
	Message string `json:"message"`
}

// PaginatedResponse — ответ с пагинацией.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// PaginationRequest — параметры пагинации из запроса.
type PaginationRequest struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize устанавливает значения по умолчанию для пагинации.
func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}
}

// Offset вычисляет смещение для SQL-запроса.
// Пример: страница 2, размер 20 → offset = 20 (пропускаем первые 20 записей).
func (p *PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// InternalError логирует реальную ошибку и возвращает безопасное сообщение для клиента.
// Это важно для безопасности: клиент не должен видеть внутренние ошибки БД.
func InternalError(err error, context string) string {
	slog.Error(context, "error", err)
	return "внутренняя ошибка сервера"
}

// TotalPages вычисляет общее количество страниц.
func TotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	// Округляем вверх: 21 запись / 20 на страницу = 2 страницы
	pages := (total + int64(pageSize) - 1) / int64(pageSize)
	return int(pages)
}
