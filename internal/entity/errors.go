package entity

import "errors"

var (
	ErrHouseholdNotFound = errors.New("household tidak ditemukan")
	ErrPickupNotFound    = errors.New("pickup tidak ditemukan")
	ErrPaymentNotFound   = errors.New("payment tidak ditemukan")
	ErrInvalidUUID       = errors.New("format ID UUID tidak valid")
	ErrDatabase          = errors.New("terjadi kesalahan pada sistem database")
	ErrValidation        = errors.New("validasi input gagal")

	ErrHasPendingPayment     = errors.New("household memiliki payment berstatus pending, tidak dapat membuat pickup baru")
	ErrPickupNotPending      = errors.New("pickup hanya dapat dijadwalkan jika status saat ini pending")
	ErrSafetyCheckRequired   = errors.New("sampah elektronik hanya bisa dijadwalkan jika pemeriksaan keamanan sudah lolos")
	ErrPickupNotScheduled    = errors.New("pickup hanya dapat diselesaikan jika status saat ini scheduled")
	ErrPickupAlreadyFinal    = errors.New("pickup sudah completed/canceled, tidak dapat diubah")
	ErrPaymentNotPending     = errors.New("payment hanya dapat dikonfirmasi jika status saat ini pending")
	ErrProofRequired         = errors.New("konfirmasi payment mewajibkan upload file bukti bayar")
	ErrPaymentAmountMismatch = errors.New("amount tidak sesuai tarif")
)
