-- Remove only the additional variant SKUs introduced by 000007.
DELETE FROM product_variants
WHERE sku ~ '^HELM-(001|002|003|004|005|006|007|008|009|010)-[ABCD]$';
