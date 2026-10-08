-- Seed production-like: 32 household, 42 pickup, 14 payment.
-- Idempotent (ON CONFLICT DO NOTHING) sehingga aman dijalankan ulang.
-- created_at disebar 60 hari ke belakang agar pagination & filter tanggal bermakna.
-- ID legacy (111..., 222..., a111..., a222...) dipertahankan agar Postman collection lama tetap jalan.

-- ============ HOUSEHOLDS (32) ============
INSERT INTO households (id, owner_name, address, created_at, updated_at) VALUES
  ('11111111-1111-1111-1111-111111111111', 'Budi Santoso', 'Jl. Merdeka No. 10, Klojen, Malang', NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days'),
  ('22222222-2222-2222-2222-222222222222', 'Siti Aminah', 'Jl. Soekarno Hatta No. 25, Lowokwaru, Malang', NOW() - INTERVAL '59 days', NOW() - INTERVAL '59 days'),
  ('10000000-0000-0000-0000-000000000003', 'Agus Wijaya', 'Jl. Ijen No. 88, Klojen, Malang', NOW() - INTERVAL '58 days', NOW() - INTERVAL '58 days'),
  ('10000000-0000-0000-0000-000000000004', 'Dewi Lestari', 'Jl. Sigura-gura No. 14, Lowokwaru, Malang', NOW() - INTERVAL '56 days', NOW() - INTERVAL '56 days'),
  ('10000000-0000-0000-0000-000000000005', 'Rina Marlina', 'Jl. Tlogomas No. 32, Lowokwaru, Malang', NOW() - INTERVAL '55 days', NOW() - INTERVAL '55 days'),
  ('10000000-0000-0000-0000-000000000006', 'Joko Prasetyo', 'Jl. Raya Tlogomas No. 57, Lowokwaru, Malang', NOW() - INTERVAL '54 days', NOW() - INTERVAL '54 days'),
  ('10000000-0000-0000-0000-000000000007', 'Endang Sulastri', 'Perum Sawojajar Blok C-12, Kedungkandang, Malang', NOW() - INTERVAL '52 days', NOW() - INTERVAL '52 days'),
  ('10000000-0000-0000-0000-000000000008', 'Herman Saputra', 'Jl. Danau Toba No. 9, Sawojajar, Malang', NOW() - INTERVAL '51 days', NOW() - INTERVAL '51 days'),
  ('10000000-0000-0000-0000-000000000009', 'Yuni Astuti', 'Jl. Danau Sentani No. 21, Madyopuro, Malang', NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days'),
  ('10000000-0000-0000-0000-000000000010', 'Dedi Kurniawan', 'Jl. Raya Madyopuro No. 44, Kedungkandang, Malang', NOW() - INTERVAL '48 days', NOW() - INTERVAL '48 days'),
  ('10000000-0000-0000-0000-000000000011', 'Maya Putri', 'Jl. A. Yani No. 17, Blimbing, Malang', NOW() - INTERVAL '47 days', NOW() - INTERVAL '47 days'),
  ('10000000-0000-0000-0000-000000000012', 'Andi Nugroho', 'Jl. R. Tumenggung Suryo No. 63, Blimbing, Malang', NOW() - INTERVAL '46 days', NOW() - INTERVAL '46 days'),
  ('10000000-0000-0000-0000-000000000013', 'Fitri Handayani', 'Jl. Laksda Adisucipto No. 28, Blimbing, Malang', NOW() - INTERVAL '44 days', NOW() - INTERVAL '44 days'),
  ('10000000-0000-0000-0000-000000000014', 'Bambang Sutrisno', 'Jl. Cengger Ayam No. 5, Lowokwaru, Malang', NOW() - INTERVAL '43 days', NOW() - INTERVAL '43 days'),
  ('10000000-0000-0000-0000-000000000015', 'Ratna Dewi', 'Jl. Veteran No. 19, Ketawanggede, Malang', NOW() - INTERVAL '42 days', NOW() - INTERVAL '42 days'),
  ('10000000-0000-0000-0000-000000000016', 'Eko Firmansyah', 'Jl. Gajayana No. 50, Dinoyo, Malang', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days'),
  ('10000000-0000-0000-0000-000000000017', 'Sri Wahyuni', 'Jl. MT Haryono No. 167, Dinoyo, Malang', NOW() - INTERVAL '39 days', NOW() - INTERVAL '39 days'),
  ('10000000-0000-0000-0000-000000000018', 'Hendra Gunawan', 'Jl. Candi Panggung No. 8, Mojolangu, Malang', NOW() - INTERVAL '38 days', NOW() - INTERVAL '38 days'),
  ('10000000-0000-0000-0000-000000000019', 'Novi Rahmawati', 'Jl. Simpang Candi Panggung No. 33, Lowokwaru, Malang', NOW() - INTERVAL '36 days', NOW() - INTERVAL '36 days'),
  ('10000000-0000-0000-0000-000000000020', 'Tono Sugiarto', 'Jl. Joyo Agung No. 71, Merjosari, Malang', NOW() - INTERVAL '35 days', NOW() - INTERVAL '35 days'),
  ('10000000-0000-0000-0000-000000000021', 'Lina Marlina', 'Jl. Joyo Utomo No. 12, Merjosari, Malang', NOW() - INTERVAL '34 days', NOW() - INTERVAL '34 days'),
  ('10000000-0000-0000-0000-000000000022', 'Wahyu Hidayat', 'Jl. Raya Langsep No. 6, Sukun, Malang', NOW() - INTERVAL '32 days', NOW() - INTERVAL '32 days'),
  ('10000000-0000-0000-0000-000000000023', 'Dian Puspita', 'Jl. S. Supriadi No. 41, Sukun, Malang', NOW() - INTERVAL '31 days', NOW() - INTERVAL '31 days'),
  ('10000000-0000-0000-0000-000000000024', 'Rudi Hartono', 'Jl. Kebonsari No. 23, Sukun, Malang', NOW() - INTERVAL '30 days', NOW() - INTERVAL '30 days'),
  ('10000000-0000-0000-0000-000000000025', 'Sari Wulandari', 'Jl. Bandulan No. 77, Sukun, Malang', NOW() - INTERVAL '28 days', NOW() - INTERVAL '28 days'),
  ('10000000-0000-0000-0000-000000000026', 'Fajar Ramadhan', 'Jl. Raya Mulyorejo No. 15, Sukun, Malang', NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
  ('10000000-0000-0000-0000-000000000027', 'Intan Permata', 'Jl. Keben No. 3, Klojen, Malang', NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),
  ('10000000-0000-0000-0000-000000000028', 'Galih Prakoso', 'Jl. Arjosari No. 52, Blimbing, Malang', NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
  ('10000000-0000-0000-0000-000000000029', 'Wulan Sari', 'Jl. Plaosan No. 18, Blimbing, Malang', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
  ('10000000-0000-0000-0000-000000000030', 'Yoga Saputra', 'Jl. Polowijen No. 29, Blimbing, Malang', NOW() - INTERVAL '18 days', NOW() - INTERVAL '18 days'),
  ('10000000-0000-0000-0000-000000000031', 'Nita Kurnia', 'Jl. Kedawung No. 11, Lowokwaru, Malang', NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
  ('10000000-0000-0000-0000-000000000032', 'Dadan Hidayat', 'Jl. Tata Surya No. 36, Dinoyo, Malang', NOW() - INTERVAL '12 days', NOW() - INTERVAL '12 days')
ON CONFLICT (id) DO NOTHING;

-- ============ WASTE PICKUPS (42) ============
-- pending 12, scheduled 10, completed 14, canceled 6. Semua tipe terwakili di tiap status.
INSERT INTO waste_pickups (id, household_id, type, status, pickup_date, safety_check, created_at, updated_at) VALUES
  -- legacy (dipertahankan untuk Postman lama)
  ('a1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'organic', 'pending', NULL, NULL, NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
  ('a2222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'electronic', 'pending', NULL, false, NOW() - INTERVAL '9 days', NOW() - INTERVAL '9 days'),
  -- pending tambahan
  ('a0000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', 'plastic', 'pending', NULL, NULL, NOW() - INTERVAL '8 days', NOW() - INTERVAL '8 days'),
  ('a0000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', 'paper', 'pending', NULL, NULL, NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
  ('a0000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'organic', 'pending', NULL, NULL, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
  ('a0000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000006', 'electronic', 'pending', NULL, true, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
  ('a0000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000007', 'plastic', 'pending', NULL, NULL, NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),
  ('a0000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000008', 'paper', 'pending', NULL, NULL, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
  ('a0000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000009', 'organic', 'pending', NULL, NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
  ('a0000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000010', 'electronic', 'pending', NULL, false, NOW() - INTERVAL '1 days', NOW() - INTERVAL '1 days'),
  ('a0000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000011', 'plastic', 'pending', NULL, NULL, NOW() - INTERVAL '20 hours', NOW() - INTERVAL '20 hours'),
  ('a0000000-0000-0000-0000-000000000012', '10000000-0000-0000-0000-000000000012', 'organic', 'pending', NULL, NULL, NOW() - INTERVAL '10 hours', NOW() - INTERVAL '10 hours'),
  -- scheduled (pickup_date masa depan, electronic selalu safety_check=true)
  ('a0000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000013', 'organic', 'scheduled', NOW() + INTERVAL '1 days', NULL, NOW() - INTERVAL '9 days', NOW() - INTERVAL '2 days'),
  ('a0000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000014', 'plastic', 'scheduled', NOW() + INTERVAL '2 days', NULL, NOW() - INTERVAL '8 days', NOW() - INTERVAL '2 days'),
  ('a0000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000015', 'paper', 'scheduled', NOW() + INTERVAL '2 days', NULL, NOW() - INTERVAL '8 days', NOW() - INTERVAL '1 days'),
  ('a0000000-0000-0000-0000-000000000016', '10000000-0000-0000-0000-000000000016', 'electronic', 'scheduled', NOW() + INTERVAL '3 days', true, NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 days'),
  ('a0000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000017', 'organic', 'scheduled', NOW() + INTERVAL '3 days', NULL, NOW() - INTERVAL '6 days', NOW() - INTERVAL '1 days'),
  ('a0000000-0000-0000-0000-000000000018', '10000000-0000-0000-0000-000000000018', 'electronic', 'scheduled', NOW() + INTERVAL '4 days', true, NOW() - INTERVAL '5 days', NOW() - INTERVAL '12 hours'),
  ('a0000000-0000-0000-0000-000000000019', '10000000-0000-0000-0000-000000000019', 'plastic', 'scheduled', NOW() + INTERVAL '5 days', NULL, NOW() - INTERVAL '4 days', NOW() - INTERVAL '10 hours'),
  ('a0000000-0000-0000-0000-000000000020', '10000000-0000-0000-0000-000000000020', 'paper', 'scheduled', NOW() + INTERVAL '6 days', NULL, NOW() - INTERVAL '3 days', NOW() - INTERVAL '8 hours'),
  ('a0000000-0000-0000-0000-000000000021', '10000000-0000-0000-0000-000000000021', 'organic', 'scheduled', NOW() + INTERVAL '7 days', NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '6 hours'),
  ('a0000000-0000-0000-0000-000000000022', '10000000-0000-0000-0000-000000000022', 'electronic', 'scheduled', NOW() + INTERVAL '7 days', true, NOW() - INTERVAL '1 days', NOW() - INTERVAL '5 hours'),
  -- completed (pickup_date masa lalu; tiap completed punya 1 payment di bawah)
  ('a0000000-0000-0000-0000-000000000023', '10000000-0000-0000-0000-000000000023', 'organic', 'completed', NOW() - INTERVAL '30 days', NULL, NOW() - INTERVAL '35 days', NOW() - INTERVAL '30 days'),
  ('a0000000-0000-0000-0000-000000000024', '10000000-0000-0000-0000-000000000024', 'plastic', 'completed', NOW() - INTERVAL '28 days', NULL, NOW() - INTERVAL '33 days', NOW() - INTERVAL '28 days'),
  ('a0000000-0000-0000-0000-000000000025', '10000000-0000-0000-0000-000000000025', 'paper', 'completed', NOW() - INTERVAL '26 days', NULL, NOW() - INTERVAL '31 days', NOW() - INTERVAL '26 days'),
  ('a0000000-0000-0000-0000-000000000026', '10000000-0000-0000-0000-000000000026', 'electronic', 'completed', NOW() - INTERVAL '24 days', true, NOW() - INTERVAL '29 days', NOW() - INTERVAL '24 days'),
  ('a0000000-0000-0000-0000-000000000027', '10000000-0000-0000-0000-000000000027', 'organic', 'completed', NOW() - INTERVAL '22 days', NULL, NOW() - INTERVAL '27 days', NOW() - INTERVAL '22 days'),
  ('a0000000-0000-0000-0000-000000000028', '10000000-0000-0000-0000-000000000028', 'plastic', 'completed', NOW() - INTERVAL '20 days', NULL, NOW() - INTERVAL '25 days', NOW() - INTERVAL '20 days'),
  ('a0000000-0000-0000-0000-000000000029', '10000000-0000-0000-0000-000000000029', 'paper', 'completed', NOW() - INTERVAL '18 days', NULL, NOW() - INTERVAL '23 days', NOW() - INTERVAL '18 days'),
  ('a0000000-0000-0000-0000-000000000030', '10000000-0000-0000-0000-000000000030', 'electronic', 'completed', NOW() - INTERVAL '16 days', true, NOW() - INTERVAL '21 days', NOW() - INTERVAL '16 days'),
  ('a0000000-0000-0000-0000-000000000031', '10000000-0000-0000-0000-000000000031', 'organic', 'completed', NOW() - INTERVAL '14 days', NULL, NOW() - INTERVAL '19 days', NOW() - INTERVAL '14 days'),
  ('a0000000-0000-0000-0000-000000000032', '10000000-0000-0000-0000-000000000032', 'plastic', 'completed', NOW() - INTERVAL '12 days', NULL, NOW() - INTERVAL '17 days', NOW() - INTERVAL '12 days'),
  ('a0000000-0000-0000-0000-000000000041', '11111111-1111-1111-1111-111111111111', 'organic', 'completed', NOW() - INTERVAL '11 days', NULL, NOW() - INTERVAL '16 days', NOW() - INTERVAL '11 days'),
  ('a0000000-0000-0000-0000-000000000042', '22222222-2222-2222-2222-222222222222', 'electronic', 'completed', NOW() - INTERVAL '10 days', true, NOW() - INTERVAL '15 days', NOW() - INTERVAL '10 days'),
  ('a0000000-0000-0000-0000-000000000033', '10000000-0000-0000-0000-000000000003', 'paper', 'completed', NOW() - INTERVAL '9 days', NULL, NOW() - INTERVAL '14 days', NOW() - INTERVAL '9 days'),
  ('a0000000-0000-0000-0000-000000000034', '10000000-0000-0000-0000-000000000004', 'electronic', 'completed', NOW() - INTERVAL '8 days', true, NOW() - INTERVAL '13 days', NOW() - INTERVAL '8 days'),
  -- canceled
  ('a0000000-0000-0000-0000-000000000035', '10000000-0000-0000-0000-000000000005', 'organic', 'canceled', NULL, NULL, NOW() - INTERVAL '25 days', NOW() - INTERVAL '24 days'),
  ('a0000000-0000-0000-0000-000000000036', '10000000-0000-0000-0000-000000000006', 'plastic', 'canceled', NULL, NULL, NOW() - INTERVAL '23 days', NOW() - INTERVAL '22 days'),
  ('a0000000-0000-0000-0000-000000000037', '10000000-0000-0000-0000-000000000007', 'electronic', 'canceled', NULL, false, NOW() - INTERVAL '21 days', NOW() - INTERVAL '20 days'),
  ('a0000000-0000-0000-0000-000000000038', '10000000-0000-0000-0000-000000000008', 'paper', 'canceled', NULL, NULL, NOW() - INTERVAL '19 days', NOW() - INTERVAL '18 days'),
  ('a0000000-0000-0000-0000-000000000039', '10000000-0000-0000-0000-000000000009', 'organic', 'canceled', NULL, NULL, NOW() - INTERVAL '17 days', NOW() - INTERVAL '16 days'),
  ('a0000000-0000-0000-0000-000000000040', '10000000-0000-0000-0000-000000000010', 'plastic', 'canceled', NULL, NULL, NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days')
ON CONFLICT (id) DO NOTHING;

-- ============ PAYMENTS (16) ============
-- 1 payment per completed pickup di atas. Nominal ikut aturan: non-electronic 50000, electronic 100000.
-- Status disebar: pending 4, paid 8 (dengan proof + payment_date), failed 2.
-- created_at disebar agar filter start_date/end_date bisa diuji.
INSERT INTO payments (id, household_id, waste_id, amount, payment_date, status, proof_file_url, created_at, updated_at) VALUES
  ('b0000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000023', 'a0000000-0000-0000-0000-000000000023', 50000, NOW() - INTERVAL '29 days', 'paid', '/uploads/seed_bukti_001.jpg', NOW() - INTERVAL '30 days', NOW() - INTERVAL '29 days'),
  ('b0000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000024', 'a0000000-0000-0000-0000-000000000024', 50000, NOW() - INTERVAL '27 days', 'paid', '/uploads/seed_bukti_002.jpg', NOW() - INTERVAL '28 days', NOW() - INTERVAL '27 days'),
  ('b0000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000025', 'a0000000-0000-0000-0000-000000000025', 50000, NULL, 'pending', NULL, NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
  ('b0000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000026', 'a0000000-0000-0000-0000-000000000026', 100000, NOW() - INTERVAL '23 days', 'paid', '/uploads/seed_bukti_004.jpg', NOW() - INTERVAL '24 days', NOW() - INTERVAL '23 days'),
  ('b0000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000027', 'a0000000-0000-0000-0000-000000000027', 50000, NULL, 'failed', NULL, NOW() - INTERVAL '22 days', NOW() - INTERVAL '21 days'),
  ('b0000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000028', 'a0000000-0000-0000-0000-000000000028', 50000, NOW() - INTERVAL '19 days', 'paid', '/uploads/seed_bukti_006.jpg', NOW() - INTERVAL '20 days', NOW() - INTERVAL '19 days'),
  ('b0000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000029', 'a0000000-0000-0000-0000-000000000029', 50000, NULL, 'pending', NULL, NOW() - INTERVAL '18 days', NOW() - INTERVAL '18 days'),
  ('b0000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000030', 'a0000000-0000-0000-0000-000000000030', 100000, NOW() - INTERVAL '15 days', 'paid', '/uploads/seed_bukti_008.jpg', NOW() - INTERVAL '16 days', NOW() - INTERVAL '15 days'),
  ('b0000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000031', 'a0000000-0000-0000-0000-000000000031', 50000, NOW() - INTERVAL '13 days', 'paid', '/uploads/seed_bukti_009.jpg', NOW() - INTERVAL '14 days', NOW() - INTERVAL '13 days'),
  ('b0000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000032', 'a0000000-0000-0000-0000-000000000032', 50000, NULL, 'failed', NULL, NOW() - INTERVAL '12 days', NOW() - INTERVAL '11 days'),
  ('b0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111111', 'a0000000-0000-0000-0000-000000000041', 50000, NULL, 'pending', NULL, NOW() - INTERVAL '11 days', NOW() - INTERVAL '11 days'),
  ('b0000000-0000-0000-0000-000000000012', '22222222-2222-2222-2222-222222222222', 'a0000000-0000-0000-0000-000000000042', 100000, NOW() - INTERVAL '9 days', 'paid', '/uploads/seed_bukti_012.jpg', NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days'),
  ('b0000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000033', 50000, NOW() - INTERVAL '8 days', 'paid', '/uploads/seed_bukti_013.jpg', NOW() - INTERVAL '9 days', NOW() - INTERVAL '8 days'),
  ('b0000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000034', 100000, NULL, 'pending', NULL, NOW() - INTERVAL '8 days', NOW() - INTERVAL '8 days')
ON CONFLICT (id) DO NOTHING;
