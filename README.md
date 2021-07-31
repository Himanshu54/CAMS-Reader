### CAMS ( Consolidate Account Statement - Detailed ) Reader

#### Install
```go get -u https://github.com/Himanshu54/CAMS-Reader```

#### Usage
```go run app.go```

```json
{
  "Folios": [
    {
      "FolioNo": "0123456789 / 0 ",
      "AMC": "XXX MUTUAL FUND",
      "PAN": "NOT OK",
      "KYC": "OK",
      "Schemes": [
        {
          "Scheme": "XXXX   Fund  -  Direct Plan Growth ",
          "Registrar": "KFINTECH",
          "Advisor": "ABC 12312312 ",
          "Scheme_type": "",
          "Rta_code": "123 A 123 A ",
          "Open": 0,
          "Close": 0,
          "Value": 01.38680076599121,
          "Nav": 01.38680076599121,
          "Valuation": {
            "Date": "05-May-1999",
            "Value": 0,
            "Nav": 01.38680076599121
          },
          "Transactions": [
            {
              "Date": "07-Sep-1999",
              "Unit": 5.55400085449219,
              "Balance": 5.55400085449219,
              "Amount": 9.9500122070312,
              "Type": "PURCHASE",
              "Price": 7.99959945678711
            },
            {
              "Date": "15-Sep-1999",
              "Unit": 8.2429962158203,
              "Balance": 20.7969970703125,
              "Amount": 06.8701171875,
              "Type": "PURCHASE",
              "Price": 7.585100173950195
            },
          ],
          "Charges": 1.38680076599121,
          "PL": 96.490051269531
        },
        {
          "Scheme": "XXX  Index Fund  -  Direct Growth ",
          "Registrar": "KFINTECH",
          "Advisor": "ABC 12121212 ",
          "Scheme_type": "",
          "Rta_code": "123 ABCD ",
          "Open": 0,
          "Close": 8.1649780273438,
          "Value": 8.384309122106,
          "Nav": 3.92959976196289,
          ...
        }
      ]
    },
    {
      "FolioNo": "1038994199 ",
      "AMC": "ABC Mutual Fund",
      "PAN": "OK",
      "KYC": "OK",
      ...
```
