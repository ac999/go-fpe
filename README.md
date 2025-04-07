# 🔐 Go-FPE

![GitHub last commit](https://img.shields.io/github/last-commit/ac999/go-fpe)

A lightweight and modular cryptographic library written in Go, implementing secure and standards-compliant algorithms such as AES (Advanced Encryption Standard) and FF1 (Format-Preserving Encryption). Built for educational purposes.

## 📦 Features

- 🔒 AES Encryption (FIPS 197 standard)
- 🔁 FF1 (NIST Format-Preserving Encryption)
- 🧩 Component-based design
- ⚙️ Utilities for encoding, transformation, and debugging
- ✅ Test coverage for core components

---

## 📂 Project Structure

```
.
├── algorithms/              # Core cryptographic implementations
│   ├── aes.go               # AES block cipher
│   ├── ff1.go               # Format-preserving encryption (FF1)
│   ├── component.go         # Cipher interfaces and wrappers
│   └── helpers.go           # Internal utility functions
├── tests/
│   └── algorithms_test.go   # Unit tests
├── main.go                  # Example or entry point (WIP)
├── .gitignore
└── LICENSE
```

---

## 🚀 Getting Started

### Requirements

- Go 1.18+

---

## 🧪 Testing

Run the test suite:

```bash
go test ./tests/...
```

---

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to fork the repository and submit a pull request.

---

## 🛡️ Disclaimer

This library is provided as-is and should not be used in production systems without a full security audit. For high-stakes cryptographic applications, always use vetted libraries and consult experts.
