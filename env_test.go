package go_env_test

import (
	"fmt"
	"testing"

	go_env "github.com/hehaowen00/go-env"
)

type Config struct {
	A int     `env:"A,default:10"`
	B string  `env:"B,default:hello, world!"`
	C bool    `env:"C,default:false"`
	D float32 `env:"D,default:3.14"`
}

func TestEnv1(t *testing.T) {
	config, err := go_env.Load[Config]()
	if err != nil {
		panic(err)
	}

	if config.A != 10 {
		t.FailNow()
	}

	if config.B != "hello, world!" {
		t.FailNow()
	}

	if config.C != false {
		t.FailNow()
	}

	if config.D != 3.14 {
		t.FailNow()
	}

	fmt.Printf("%+v\n", config)
}
