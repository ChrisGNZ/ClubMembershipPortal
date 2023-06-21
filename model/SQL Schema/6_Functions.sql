---------------------------------------------------------------------------------------------------------------
--https://stackoverflow.com/questions/12504732/strip-non-numeric-characters-from-a-string
---------------------------------------------------------------------------------------------------------------
go
CREATE FUNCTION dbo.StripNonNumeric (@string VARCHAR(5000))
    RETURNS VARCHAR(1000)
AS
BEGIN
    SET @string = REPLACE(@string, ',', '.')
    SET @string = (SELECT   SUBSTRING(@string, v.number, 1)
                   FROM     master..spt_values v
                   WHERE    v.type = 'P'
                     AND v.number BETWEEN 1 AND LEN(@string)
                     AND (SUBSTRING(@string, v.number, 1) LIKE '[0-9]'
                       OR SUBSTRING(@string, v.number, 1) LIKE '[.]')
                   ORDER BY v.number
                   FOR
                       XML PATH('')
    )
    RETURN @string
END
GO
-------------------------------------