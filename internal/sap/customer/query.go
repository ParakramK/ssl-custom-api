package customer

const CustomerAgingQuery = `
SELECT
    IFNULL(T0."DocEntry", 0) AS "DocEntry",
    T0."DocNum" AS "Invoice Number",
    T0."DocDate" AS "Invoice Date",
    C."CardCode" AS "Customer Code",
    C."CardName" AS "Customer Name",
    C."Balance" AS "Balance",
    T2."TaxId2" AS "Tax Number",
    T0."DocDueDate" AS "Due Date",
    T1."PymntGroup" AS "Payment Terms",
    T3."SlpName" AS "Sales Employee",

    CASE
        WHEN T0."GroupNum" IN (8,12,10,9,21) THEN 'BG'
        WHEN T0."GroupNum" IN (15,14,11,13,26,27) THEN 'LC'
        ELSE 'Other'
    END AS "Payment Term Group",

    IFNULL(
        DAYS_BETWEEN(T0."DocDate", CURRENT_DATE),
        0
    ) AS "Outstanding Days",

    IFNULL(T0."DocTotal", 0) AS "Invoice Amount",

    IFNULL(T0."PaidToDate", 0) AS "Paid Amount",

    IFNULL(T0."DocTotal", 0)
        - IFNULL(T0."PaidToDate", 0)
        AS "Outstanding Amount",

    CASE
        WHEN DAYS_BETWEEN(T0."DocDate", CURRENT_DATE)
            BETWEEN 0 AND 30
        THEN IFNULL(T0."DocTotal", 0)
            - IFNULL(T0."PaidToDate", 0)
        ELSE 0
    END AS "0-30 Days",

    CASE
        WHEN DAYS_BETWEEN(T0."DocDate", CURRENT_DATE)
            BETWEEN 31 AND 60
        THEN IFNULL(T0."DocTotal", 0)
            - IFNULL(T0."PaidToDate", 0)
        ELSE 0
    END AS "31-60 Days",

    CASE
        WHEN DAYS_BETWEEN(T0."DocDate", CURRENT_DATE)
            BETWEEN 61 AND 90
        THEN IFNULL(T0."DocTotal", 0)
            - IFNULL(T0."PaidToDate", 0)
        ELSE 0
    END AS "61-90 Days",

    CASE
        WHEN DAYS_BETWEEN(T0."DocDate", CURRENT_DATE) >= 91
        THEN IFNULL(T0."DocTotal", 0)
            - IFNULL(T0."PaidToDate", 0)
        ELSE 0
    END AS "91+ Days"

FROM OCRD C

LEFT JOIN OINV T0
    ON C."CardCode" = T0."CardCode"
    AND T0."DocStatus" = 'O'

LEFT JOIN OCTG T1
    ON T0."GroupNum" = T1."GroupNum"

LEFT JOIN (
    SELECT
        "CardCode",
        MAX("TaxId2") AS "TaxId2"
    FROM CRD7
    GROUP BY "CardCode"
) T2
    ON C."CardCode" = T2."CardCode"

LEFT JOIN OSLP T3
    ON C."SlpCode" = T3."SlpCode"

WHERE C."CardCode" = ?

ORDER BY T0."DocDate"
`
