# go-env

Environment variable loader for structs in go

```go
type Config struct {
	A int     `env:"A,default:10"`
	B string  `env:"B,default:hello, world!"`
	C bool    `env:"C,default:false"`
	D float32 `env:"D,default:3.14"`
}

func main() {
	config, err := go_env.Load[Config]()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)
}
```
