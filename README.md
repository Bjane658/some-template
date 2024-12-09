# some-template

some is simple CLI tool to manage and add templates of source code files. With this tool, you can easily store reusable code templates and quickly create new files in the current folder based on those templates.

## How to build and install `some`

### build
```
make 
```

### install `some` to /usr/local/bin
```
make install
```

# How to use `some`

### Add template

```
some -a java spring HttpClient
```

### View template

```
some -v java spring HttpClient
```

### Edit template

```
some -e java spring HttpClient
```

### Apply template

```
some java spring HttpClient
```
