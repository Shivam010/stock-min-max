# stock-min-max
stock-min-max is an application based on the **MCX India**'s endpoint for getting Commodities' stock prices throughout the day.  
```shell script
curl 'https://www.mcxindia.com/BackPage.aspx/GetGraphForScrip' \
  -H 'Content-Type: application/json; charset=UTF-8' \
  -H 'Referer: https://www.mcxindia.com/home' \
  --data-binary '{"Commodity":"Crudeoil"}' \
  --compressed
```

## Endpoint

```shell script
GET /api/v1?commodity={value}
```
