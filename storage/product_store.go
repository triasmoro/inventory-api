package storage

import (
	"database/sql"
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

// ExportProducts method
func (s *Store) ExportProducts() ([]model.Product, error) {
	query := `SELECT
		p.id AS product_id,
		p.name,
		pv.id AS variant_id,
		pv.sku,
		pov.id AS option_value_id,
		po.name AS option_name,
		pov.value AS option_value
	FROM products p
	INNER JOIN product_variants pv ON pv.product_id = p.id
	INNER JOIN product_variant_options pvo ON pvo.product_variant_id = pv.id
	INNER JOIN product_option_values pov ON pov.id = pvo.product_option_value_id
	INNER JOIN product_options po ON po.id = pov.product_option_id
	WHERE pv.fg_delete = 0
	ORDER BY p.id, pv.id ASC`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	var variants []model.Variant
	var options []model.Option

	var prevProduct model.Product
	var prevVariant model.Variant
	for rows.Next() {
		var productItem model.Product
		var variantItem model.Variant
		var optionItem model.Option

		err = rows.Scan(
			&productItem.ID,
			&productItem.Name,
			&variantItem.ID,
			&variantItem.SKU,
			&optionItem.ID,
			&optionItem.Name,
			&optionItem.Value,
		)
		if err != nil {
			return products, err
		}

		// assign product id at variant object
		variantItem.ProductID = productItem.ID

		// populate variant
		if prevVariant.ID != 0 && variantItem.ID != prevVariant.ID {
			prevVariant.Options = options
			variants = append(variants, prevVariant)
			options = []model.Option{} // reset options
		}

		// populate option
		options = append(options, optionItem)

		// populate product
		if prevProduct.ID != 0 && productItem.ID != prevProduct.ID {
			prevProduct.Variants = variants
			products = append(products, prevProduct)
			variants = []model.Variant{} // reset variants
		}

		prevProduct = productItem
		prevVariant = variantItem
	}

	// insert last iteration data
	prevVariant.Options = options
	variants = append(variants, prevVariant)
	prevProduct.Variants = variants
	products = append(products, prevProduct)

	return products, nil
}

// GetProductByID method
func (s *Store) GetProductByID(id int) (model.Product, error) {
	query := fmt.Sprintf(`SELECT id, name FROM products WHERE id = %d`, id)

	var p model.Product
	if err := s.DB.QueryRow(query).Scan(&p.ID, &p.Name); err != nil {
		return p, err
	}

	return p, nil
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

// IsProductVariantExist method
func (s *Store) IsProductVariantExist(id int) (bool, error) {
	var exist int
	query := "SELECT 1 FROM product_variants WHERE id = ? AND fg_delete = 0"
	if err := s.DB.QueryRow(query, id).Scan(&exist); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// UpdateProduct method
func (s *Store) UpdateProduct(p model.Product) error {
	query := fmt.Sprintf(`UPDATE products SET name = "%s" WHERE id = %d`, p.Name, p.ID)
	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	return nil
}
