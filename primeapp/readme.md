## Testing comands

_Para ejecutar las pruebas unitarias en un paquete_

```
 go test
```

_Para ejecutar las pruebas unitarias de una función o un grupo de funciones_

```
 go test -run TestName

```

go test

```

_Para ejecutar las pruebas unitarias y ver los detalles_

```

go test -v

```

_Para ejecutar las pruebas unitarias de todo el proyecto_

```

go test -v ./...

```

_Para ejecutar las pruebas unitarias y ver el coverage_

```

go test -v -cover

```

_Para ejecutar las pruebas unitarias y crear archivo de salida para el coverage_

```

go test -coverprofile=coverage.out

```

_Para transformar el file coverage aun archivo html_

```

go tool cover -html=coverage.out -o coverage.html

```

```
