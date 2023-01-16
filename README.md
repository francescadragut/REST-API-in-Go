# REST-API-in-Go

This REST API in Go is used for a data analysis project on Romanian companies. The database is in PostgreSQL. Since the data is protected, check the next section on how to create the database in order to match the structure of the original database.

## Create database

```
CREATE TABLE romanian_companies_financial_data_serial(
	id SERIAL PRIMARY KEY,
	denumire VARCHAR(500),
	cui VARCHAR(500),
	an VARCHAR(500),
	caen VARCHAR(500),
	den_caen VARCHAR(500),
	numar_mediu_de_salariati VARCHAR(500),
	pierdere_neta VARCHAR(500),
	profit_net VARCHAR(500),
	cifra_de_afaceri_neta VARCHAR(500),
	stocuri VARCHAR(500)
)
```

** Important: ** The `id` field should be created as `SERIAL`, because the `POST` request uses the incremental id automatically.
