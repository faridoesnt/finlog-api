# Finlog API

Finlog API is a privacy-focused personal finance backend built with **Golang (Fiber)**, **MySQL**, and deployed using **Docker**.

This service powers encrypted financial operations such as transactions, categories, import batches, and user account management. All sensitive payloads are encrypted client-side before reaching the API.

---

## 📦 Tech Stack

- **Go (Fiber Framework)**
- **MySQL 8**
- **Docker & Docker Compose**
- **Caddy (HTTPS & Reverse Proxy)**
- **GitHub Actions (CI/CD Deployment)**

---

## 🔐 Security Notes

- All sensitive content (transactions, categories, imported data) is encrypted on the **client side** before being sent.
- Server stores **only ciphertext**.
- API keys, database credentials, and runtime configs are managed via environment variables.

---

## 📜 License

Proprietary — internal use only.

---

## 👤 Author

**Farid Haikal**  
Finlog — Privacy-First Personal Finance System