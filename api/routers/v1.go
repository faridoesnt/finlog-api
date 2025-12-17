package routers

import (
	"finlog-api/api/constants"
	"finlog-api/api/contracts"
	"finlog-api/api/handlers"
	"finlog-api/api/middlewares"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Init(app *contracts.App) {
	app.Fiber.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Fiber.Get("/activated", handlers.ActivatedHandler)
	app.Fiber.Get("/privacy", func(c *fiber.Ctx) error {
		currentYear := time.Now().Year()
		lastUpdated := time.Now().Format("02 January 2006")

		html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="utf-8" />
		<title>Privacy Policy – FinLog</title>
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<meta name="robots" content="index, follow" />
		<style>
			body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI",
						Roboto, Helvetica, Arial, sans-serif;
			background-color: #ffffff;
			color: #111827;
			margin: 0;
			padding: 0;
			}
			.container {
			max-width: 720px;
			margin: 40px auto;
			padding: 0 20px 60px;
			line-height: 1.7;
			}
			h1 {
			font-size: 28px;
			margin-bottom: 8px;
			}
			h2 {
			font-size: 18px;
			margin-top: 32px;
			margin-bottom: 8px;
			}
			p, li {
			font-size: 15px;
			}
			ul {
			padding-left: 20px;
			}
			.updated {
			color: #6b7280;
			font-size: 13px;
			margin-bottom: 24px;
			}
			footer {
			margin-top: 48px;
			font-size: 13px;
			color: #6b7280;
			}
		</style>
		</head>
		<body>
		<div class="container">
			<h1>Privacy Policy – FinLog</h1>
			<div class="updated">Last updated: %s</div>

			<p>
			FinLog (“kami”, “aplikasi”) adalah aplikasi pencatat keuangan pribadi
			yang dirancang untuk membantu pengguna mengelola pemasukan dan
			pengeluaran secara mandiri. Privasi pengguna adalah hal yang penting
			bagi kami. Dokumen ini menjelaskan bagaimana data dikumpulkan,
			digunakan, dan dilindungi.
			</p>

			<h2>1. Informasi yang Kami Kumpulkan</h2>
			<ul>
			<li>
				<strong>Informasi Akun:</strong> alamat email yang digunakan untuk
				pendaftaran, verifikasi akun, dan autentikasi pengguna.
			</li>
			<li>
				<strong>Data Keuangan:</strong> data transaksi seperti pemasukan,
				pengeluaran, kategori, dan catatan yang dimasukkan langsung oleh
				pengguna untuk keperluan pencatatan keuangan pribadi.
			</li>
			</ul>

			<h2>2. Cara Kami Menggunakan Informasi</h2>
			<ul>
			<li>Menyediakan fungsi utama aplikasi</li>
			<li>Menyimpan dan menampilkan data keuangan pengguna</li>
			<li>Mengelola autentikasi dan keamanan akun</li>
			<li>Meningkatkan stabilitas dan pengalaman penggunaan aplikasi</li>
			</ul>

			<h2>3. Penyimpanan dan Keamanan Data</h2>
			<p>
			Data diproses hanya untuk keperluan aplikasi dan dilindungi oleh
			mekanisme keamanan aplikasi, termasuk perlindungan tambahan seperti
			PIN atau biometrik jika diaktifkan oleh pengguna.
			</p>

			<h2>4. Berbagi Data dengan Pihak Ketiga</h2>
			<p>
			FinLog tidak membagikan data pribadi atau data keuangan pengguna kepada
			pihak ketiga, kecuali jika diwajibkan oleh hukum yang berlaku.
			</p>

			<h2>5. Hak Pengguna</h2>
			<ul>
			<li>Mengakses data pribadi dan data keuangan</li>
			<li>Menghapus data yang tersimpan di aplikasi</li>
			<li>Berhenti menggunakan aplikasi kapan saja</li>
			</ul>

			<h2>6. Data Anak-anak</h2>
			<p>
			FinLog tidak ditujukan untuk anak-anak dan tidak secara sengaja
			mengumpulkan data dari pengguna di bawah usia yang diizinkan oleh
			hukum.
			</p>

			<h2>7. Perubahan Kebijakan Privasi</h2>
			<p>
			Kebijakan Privasi ini dapat diperbarui dari waktu ke waktu dan akan
			ditampilkan melalui halaman ini.
			</p>

			<h2>8. Kontak</h2>
			<p>
			Jika Anda memiliki pertanyaan mengenai Kebijakan Privasi ini, silakan
			hubungi kami melalui email:
			<br />
			<strong>faridhaikaal@gmail.com</strong>
			</p>

			<footer>
			© %d FinLog. All rights reserved.
			</footer>
		</div>
		</body>
		</html>
		`, lastUpdated, currentYear)

		return c.Type("html").SendString(html)
	})

	app.Fiber.Get("/account-deletion", func(c *fiber.Ctx) error {
		currentYear := time.Now().Year()
		lastUpdated := time.Now().Format("02 January 2006")

		html := fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="utf-8" />
		<title>Penghapusan Akun dan Data – FinLog</title>
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<meta name="robots" content="index, follow" />
		<style>
			body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI",
						Roboto, Helvetica, Arial, sans-serif;
			background-color: #ffffff;
			color: #111827;
			margin: 0;
			padding: 0;
			}
			.container {
			max-width: 720px;
			margin: 40px auto;
			padding: 0 20px 60px;
			line-height: 1.7;
			}
			h1 {
			font-size: 28px;
			margin-bottom: 8px;
			}
			h2 {
			font-size: 18px;
			margin-top: 32px;
			margin-bottom: 8px;
			}
			p, li {
			font-size: 15px;
			}
			ul {
			padding-left: 20px;
			}
			.updated {
			color: #6b7280;
			font-size: 13px;
			margin-bottom: 24px;
			}
			footer {
			margin-top: 48px;
			font-size: 13px;
			color: #6b7280;
			}
		</style>
		</head>
		<body>
		<div class="container">
			<h1Penghapusan Akun dan Data - FinLog</h1>
			<div class="updated">Last updated: %s</div>

			<p>
			Pengguna dapat meminta penghapusan akun dan data dengan cara:
			</p>

			<p>1. Mengirim email ke faridhaikaal@gmail.com</p>
			<p>2. Gunakan subjek: Permintaan Penghapusan Akun FinLog</p>
			<p>3. Sertakan email yang terdaftar di aplikasi</p>

			<p>
			Data yang dihapus:
			</p>

			<ul>
			<li>Akun pengguna</li>
			<li>Data transaksi</li>
			<li>Data autentikasi</li>
			</ul>

			<p>
			Retensi data:
			</p>

			<ul>
			<li>Data dihapus permanen maksimal 30 hari setelah permintaan diterima</li>
			</ul>

			<footer>
			© %d FinLog. All rights reserved.
			</footer>
		</div>
		</body>
		</html>
		`, lastUpdated, currentYear)

		return c.Type("html").SendString(html)
	})

	// Support both /api and /api/v1 prefixes for compatibility.
	registerAPIRoutes(app, app.Fiber.Group("/api"))
	registerAPIRoutes(app, app.Fiber.Group("/api/v1"))
	registerAPIRoutes(app, app.Fiber.Group("/v1"))
	RegisterWebhook(app)
}

func registerAPIRoutes(app *contracts.App, api fiber.Router) {
	authGroup := api.Group("/auth")
	authGroup.Post("/login", handlers.AuthLogin)
	authGroup.Post("/register", handlers.Register)
	authGroup.Post("/resend-verification", handlers.ResendVerification)
	authGroup.Get("/verify", handlers.VerifyEmail)
	authGroup.Post("/refresh", handlers.Refresh)

	jwtTTL := parseDuration(app.Config[constants.JWT_TTL], time.Hour)
	protected := api.Group("", middlewares.JWT([]byte(app.Config[constants.JWT_SECRET]), jwtTTL))

	protected.Post("/auth/logout", handlers.Logout)

	protected.Get("/categories", handlers.GetCategories)
	protected.Post("/categories", handlers.CreateCategory)
	protected.Put("/categories/:id", handlers.UpdateCategory)
	protected.Delete("/categories/:id", handlers.DeleteCategory)

	protected.Get("/recent-transactions", handlers.GetRecentTransactions)
	protected.Get("/transactions", handlers.GetTransactions)
	protected.Post("/transactions", handlers.CreateTransaction)
	protected.Post("/transactions/import", handlers.ImportTransactions)
	protected.Get("/transactions/import/history", handlers.ImportHistory)
	protected.Delete("/transactions/import/:batch_id", handlers.UndoImportBatch)
	protected.Put("/transactions/:id", handlers.UpdateTransaction)
	protected.Put("/transactions/bulk/notes", handlers.UpdateTransactionNotes)
	protected.Put("/transactions/bulk/amounts", handlers.UpdateTransactionAmount)
	protected.Put("/transactions/bulk/dates", handlers.UpdateTransactionDate)
	protected.Delete("/transactions/:id", handlers.DeleteTransaction)
	protected.Delete("/transactions/bulk/delete", handlers.DeleteTransactions)

	protected.Get("/budget", handlers.GetBudget)

	keyGroup := protected.Group("/keys")
	keyGroup.Post("/backup", handlers.StoreKeyBackup)
	keyGroup.Put("/backup/rotate", handlers.RotateKeyBackup)
	keyGroup.Get("/backup", handlers.GetActiveKeyBackup)
	keyGroup.Get("/backup/status", handlers.GetKeyBackupStatus)
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	return fallback
}
