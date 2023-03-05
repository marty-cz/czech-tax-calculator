# Czech Tax Calculator

Go console application which processes buys, sells, dividends, additional income and fees of Stocks and Crypto currencies and creates "Tax Year Report" using **Czech** law and output currency.

When executed it generates reports for each year (with some revenue) in Excel format. The file contains separate sheets for Stocks and Cryptocurrency overview.

## Taxes from Czech law point of view

It is mandatory to file taxes of year profit:

* if year  *stocks* (real *stocks*, *ETF*s, *fonds* or *bonds*) revenue is more than 100,000 CZK
* or if year *other* (*cryptos*, *dividends*, *CFD*s, ...) revenue is more than 6,000 CZK

The revenue might be lowered by revenue of *stocks* (real *stocks*/*ETF*s/...) hold more than 3 years. In case of this time tested revenue is more than 5,000,000 CZK it needs to be noted (but no tax is paid).

### Exchange Rate

There are two options but they have to be used consistently through whole tax report:

1. **Uniform year exchange rate** is published by ČNB (Czech National Bank) for *previous calendar year*.
2. **Daily exchange rate** is published by ČNB (Czech National Bank) for *past business day*. This option is available for people doing *bookkeeping* only (see [Pokyn GFŘ-D-54](https://www.sagit.cz/info/fz22001)).

### Purchase Price of Stock/Cryptocurrency

The purchase price contains buy price of sold particular stock/crypto and its buy fee and a broker's provision.

:warning: A sell transaction might be selling stocks/cryptos from multiple buy transactions!

So, there are also two options how to calculate purchase price but they have to be used consistently through whole tax report:

1. **FIFO** method is based on selling oldest bought item (buy transaction) first.
2. **Weighted arithmetic average** method is based on averaging total purchases. It is complicated and requires recalculation of the average item price. :warning: Thus it is not implemented.

### Cryptocurrencies

Cryptocurrencies are treated as an *Intangible moving asset* ("Nehmotný movitý majetek") => *Other income* ("Ostatní příjmy")  by Czech law (at least in 2022).

This means there is no time test available and it is not possible to combine profit from stocks and cryptos!

## Application Parameters

```raw
Usage of ./out/bin/czech-tax-calculator-linux:
  --crypto-input string
        File path to input file with Crypto-currencies transaction records
  --stock-input string
        File path to input file with Stocks transaction records
  --year string
        Target year for taxes (default "Previous Tax Year")
```

Run command below in case of this documentation is out of date:

```shell
./out/bin/czech-tax-calculator-linux -h
```

### Input data format

Please see [examples](./examples) directory which covers form of Stock and Cryptocurrency source data.

## Build and Run

See [Makefile](./Makefile) for more details

```shell
make build
./out/bin/czech-tax-calculator-linux --stock-input ./examples/Ucetni-kniha-Akcie.xlsx --crypto-input ./examples/Ucetni-kniha-Crypto.xlsx
```

or

```shell
make buildAndRun
```

## References

* [Zdanění kryptoměn](https://finex.cz/zdaneni-kryptomen-kompletni-navod/)
* [Zdanění příjmů z akcií](https://luciekocmanova.cz/zdaneni-prijmu-z-akcii/)
* [Prodali jste akcie nebo podílové listy. Co teď s daněmi](https://www.penize.cz/dan-z-prijmu-fyzickych-osob/425326-jak-zdanit-prijmy-z-prodeje-akcii-a-podilovych-listu-investice-a-danove-priznani)
* [Akcie a daňové přiznání. Jak určit pořizovací cenu a kurz](https://www.penize.cz/investice/425626-dan-z-prodeje-cennych-papiru-jak-urcit-porizovaci-cenu-a-kurz)
* [Zdanění dividend: Jakou sazbou se daní dividendy a jak postupovat při vyplnění daňového přiznání?](https://finex.cz/zdaneni-dividend-jak-postupovat-pri-vyplneni-danoveho-priznani/)
