## Usage

All commands could be used are defined by the Makefile.

1.1 To build tool with sanitizing data race
```bash
make build-race
```

1.2 To build tool (no sanitizing data race)
```bash
make build
```

2 To run tests (with sanitizing data race case)
```bash
make test
```

3 To use built tool please run something like this:
```bash
./bin/myhttp -parallel=8 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
```
