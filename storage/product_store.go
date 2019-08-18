package storage

import (
	"fmt"

	"github.com/triasmoro/inventory-api/model"
)

// SaveProduct method
func (s *Store) SaveProduct(product *model.Product) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	// save product
	stmt, err := tx.Prepare("INSERT INTO products(name) VALUES (?)")
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := stmt.Exec(product.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	productID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// assign id generated
	product.ID = int(productID)

	for i, variant := range product.Variants {
		// save product variant
		stmt, err := tx.Prepare("INSERT INTO product_variants(product_id, sku) VALUES (?, ?)")
		if err != nil {
			tx.Rollback()
			return err
		}

		res, err := stmt.Exec(product.ID, variant.SKU)
		if err != nil {
			tx.Rollback()
			return err
		}

		variantID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		// assign id generated
		product.Variants[i].ProductID = int(productID)
		product.Variants[i].ID = int(variantID)

		for _, option := range variant.Options {
			// save product variant options
			stmt, err := tx.Prepare("INSERT INTO product_variant_options(product_variant_id, product_option_value_id) VALUES (?, ?)")
			if err != nil {
				tx.Rollback()
				return err
			}

			if _, err := stmt.Exec(variantID, option.ID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// GetOptionValueID method
func (s *Store) GetOptionValueID(name, value string) (int, error) {
	query := fmt.Sprintf(`SELECT pov.id 
		FROM product_options po
		INNER JOIN product_option_values pov
			ON pov.product_option_id = po.id
		WHERE po.name = "%s"
		AND pov.value = "%s"`, name, value)

	var id int64
	if err := s.DB.QueryRow(query).Scan(&id); err != nil {
		return 0, err
	}

	return int(id), nil
}
