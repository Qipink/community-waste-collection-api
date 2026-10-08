package service

// NormalizePagination menyamakan aturan paging di semua list endpoint.
// Dipakai service (untuk query) dan handler (untuk meta) agar angka yang
// dilaporkan ke klien selalu sama dengan yang benar-benar dieksekusi.
func NormalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}
