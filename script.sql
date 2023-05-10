select * from [dbo].[product]
--delete from [dbo].[product]
--drop table [dbo].[product_info]
/*
CREATE TABLE product (
    id          UNIQUEIDENTIFIER PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       FLOAT NOT NULL,
    sku         VARCHAR(50) NOT NULL,
    created_on  DATETIME,
    updated_on  DATETIME,
    deleted_on  DATETIME
);
*/

SELECT TOP 10
			ID,
			Name,
			Description,
			Price,
			SKU,
			Created_On,
			Updated_On,
			Deleted_On
		FROM
			product
		ORDER BY
			Created_On DESC

