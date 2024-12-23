# Clever Cloud GoAccess Logs Converter

This application uses [Clever Tools](https://github.com/CleverCloud/clever-tools) to fetch access logs from Clever Cloud and convert them to use with GoAccess (COMMON logs format).

## How to use

1. Install Clever Tools: `npm install -g clever-tools` or with [your favorite package manager](https://github.com/CleverCloud/clever-tools?tab=readme-ov-file#installation)
2. Run `clever login` and login with your [Clever Cloud account](https://console.clever-cloud.com/)
3. Check you're logged in with `clever profile`

To fetch access logs, you'll need :
* The application ID
* Your user or organization ID

You can find them in [Clever Cloud's Console](https://console.clever-cloud.com/).

Then, run the following command:

```bash
./clever-access-logs-converter -org <ORG_ID> -app <APP_ID> -since <ISO_8601_DATE_HOUR>
./clever-access-logs-converter -org <ORG_ID> -app <APP_ID> -since <ISO_8601_DATE_HOUR> --until <ISO_8601_DATE_HOUR>
```

## Advanced usage

You can go further with options:

```
./clever-access-logs-converter --help
Configuration error: missing required parameters

clever-access-logs-converter (vfccb720-dirty)
Convert Clever Cloud access logs to use with GoAccess (COMMON format)

Required:
  -org string        Organization ID
  -app string        Application ID
  -since string      Start date (ISO 8601 format)

Optional:

  -deployment string Filter by deployment ID
  -instance string   Filter by instance ID
  -limit int         Limit number of results (min: 1)
  -until string      End date (ISO 8601 format)

  -out string        Output file name (default "goaccess_logs.txt")
```

## Fetch access logs from GoAccess

You can now use the generated file with GoAccess:

```bash
# Fetch data in GoAccess CLI interface
goaccess goaccess_logs.txt --log-format=COMMON

# Open the generated report.html file in your browser with real-time update or not
goaccess goaccess_logs.txt --log-format=COMMON -o report.html
goaccess goaccess_logs.txt --log-format=COMMON --real-time-html -o report.html
```

## Build & Install

You can build the application with the provided Makefile:

```bash
git clone https://github.com/davlgd/goaccess-logs-converter.git
cd goaccess-logs-converter

make clean
make build
```

You can also install it in your `$GOPATH/bin`:

```bash
make install
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

## Contributing

Feel free to contribute to this project by creating [issues](https://github.com/davlgd/goaccess-logs-converter/issues) or a [pull request](https://github.com/davlgd/goaccess-logs-converter/pulls) with your changes.
