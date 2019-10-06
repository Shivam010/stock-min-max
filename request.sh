curl 'https://www.mcxindia.com/BackPage.aspx/GetIntradayGraph_Quote' \
  -H 'content-type: application/json; charset=UTF-8' \
  --data-binary '{"Commodity":"CRUDEOIL","Expiry":"21OCT2019"}' \
  --compressed

curl 'https://www.mcxindia.com/BackPage.aspx/GetGraphForScrip' \
  -H 'Content-Type: application/json; charset=UTF-8' \
  -H 'Referer: https://www.mcxindia.com/home' \
  --data-binary '{"Commodity":"Crudeoil"}' \
  --compressed
