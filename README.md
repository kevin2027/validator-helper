# validator-helper

## Installation

---

Use go get.

```go
go get github.com/yaofei517/validator-helper
```

Then import the validator package into your own code.

```go
import "github.com/yaofei517/validator-helper"
```

## 使用

通过 msg 标签标识错误返回消息

```go
type A struct {
	String      string   `validate:"required,max=3" msg:"A 至少需要一个字符串"`
	Slice       []string `validate:"min=1,dive,min=5" msg:"A 切片至少需要一个大于5字符的字符串"`
	SliceStruct [1]*B    `validate:"min=1,dive" msg:"A 至少需要一个参数"`
	Struct      *B       `validate:"required" msg:"A 没有结构体参数"`
}
type B struct {
	String string `validate:"required,max=3" msg:"B 至少需要一个小于3字符字符串"`
}

func TestValidator(t *testing.T) {
	validate := validator.New()
	obj := A{
		String: "15",
		Slice:  []string{"1266676"},
		SliceStruct: [1]*B{
			{String: "1278787"},
		},
		Struct: &B{
			"136",
		},
	}
	err := validate.Struct(obj)
	if err != nil {
		err = util.ValidatorHelper(obj, err)
		fmt.Printf("fail: %s\n", err)
	}
	fmt.Println("ok")
}

```
