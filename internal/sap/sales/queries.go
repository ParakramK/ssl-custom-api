package sales

const SALES_ORDER_VS_BILLING_QUERY = `
SELECT
    C."CardCode" AS "Customer Code",
    C."CardName" AS "Customer Name",
    T3."SlpName" AS "Sales Employee",


    COALESCE(O."Ordered Qty", 0) AS "Ordered Quantity",
    COALESCE(I."Billed Qty", 0) AS "Billed Quantity",

    COALESCE(I."Billed Qty", 0)
        - COALESCE(O."Ordered Qty", 0) AS "Quantity Variance",

    COALESCE(O."Ordered Amount", 0) AS "Ordered Amount",
    COALESCE(I."Billed Amount", 0) AS "Billed Amount",

    COALESCE(I."Billed Amount", 0)
        - COALESCE(O."Ordered Amount", 0) AS "Amount Variance"

FROM OCRD C

LEFT JOIN (
    SELECT
        O."CardCode",

        SUM(R."Quantity") AS "Ordered Qty",

        SUM(R."LineTotal") AS "Ordered Amount"

    FROM ORDR O

    INNER JOIN RDR1 R
        ON O."DocEntry" = R."DocEntry"

    WHERE O."CANCELED" = 'N'

    GROUP BY
        O."CardCode"
) O
    ON C."CardCode" = O."CardCode"

LEFT JOIN (
    SELECT
        I."CardCode",

        SUM(I1."Quantity") AS "Billed Qty",

        SUM(I1."LineTotal") AS "Billed Amount"

    FROM OINV I

    INNER JOIN INV1 I1
        ON I."DocEntry" = I1."DocEntry"

    WHERE I."CANCELED" = 'N'

    GROUP BY
        I."CardCode"
) I
    ON C."CardCode" = I."CardCode"

LEFT JOIN OSLP T3
    ON C."SlpCode" = T3."SlpCode"

WHERE C."CardType" = 'C'

ORDER BY
    "Billed Amount" DESC
LIMIT %d;
`
