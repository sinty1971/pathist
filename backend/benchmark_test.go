package main

import (
	"slices"
	"testing"
)

// 簡易的なCompany構造体
type Company struct {
	ID         string
	FolderName string
}

// モックデータの準備
func generateCompanies(n int) []Company {
	companies := make([]Company, n)
	for i := 0; i < n; i++ {
		companies[i] = Company{
			ID:         "TC" + string(rune('A'+i%26)) + string(rune('0'+i/26)),
			FolderName: "test-folder",
		}
	}
	return companies
}

// 元のforループ実装
func findCompanyForLoop(companies []Company, id string) *Company {
	for _, company := range companies {
		if company.ID == id {
			return &company
		}
	}
	return nil
}

// slices.IndexFunc実装
func findCompanySlicesIndexFunc(companies []Company, id string) *Company {
	idx := slices.IndexFunc(companies, func(c Company) bool {
		return c.ID == id
	})
	if idx != -1 {
		return &companies[idx]
	}
	return nil
}

// ベンチマーク: 中規模スライス（100要素）
func BenchmarkFindCompany_Medium(b *testing.B) {
	companies := generateCompanies(100)
	targetID := companies[50].ID

	b.Run("ForLoop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findCompanyForLoop(companies, targetID)
		}
	})

	b.Run("SlicesIndexFunc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = findCompanySlicesIndexFunc(companies, targetID)
		}
	})
}

// ベンチマーク: 大規模スライス（1000要素）
func BenchmarkFindCompany_Large(b *testing.B) {
	companies := generateCompanies(1000)
	targetID := companies[500].ID

	b.Run("ForLoop", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = findCompanyForLoop(companies, targetID)
		}
	})

	b.Run("SlicesIndexFunc", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = findCompanySlicesIndexFunc(companies, targetID)
		}
	})
}
