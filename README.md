# pkgwatch

Checks upstream GitHub/GitLab releases against pkgver in local PKGBUILDs and prints what's outdated.

## Build

```
go build -o pkgwatch .
```

## Usage

```
pkgwatch ./packages/*
pkgwatch --json
```
