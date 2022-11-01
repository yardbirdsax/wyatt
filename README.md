# wyatt

![](doc/img/gophermarshall_small.png)

Wyatt is Go library that marshals / unmarshals GitHub Action context data into / from structs. It's named for
the (in)famous U.S. Marshal [Wyatt
Earp](https://history.howstuffworks.com/historical-figures/wyatt-earp.htm).

Wyatt makes it easy to populate Go variables with data provided by the GitHub
Actions runtime environment, specifically in the context of taking inputs. The Actions platform
provides as [specially named environment
variables](https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#example-specifying-inputs),
and populating their values into structs.

The sensible way to do this is with struct tags, so that one can simply decorate their struct and
let the library interpret how to build the bridge between the Actions runtime and the struct fields.

One way to do this would be creating our own struct tags and using methods of the `reflect` standard
library to interpret them. However, that involves a lot of rather complex code that isn't strictly
necessary, because code similar to it has already been written by very smart folks in the standard
library! Specifically, the `encoding/json` package already does the heavy lifting for us, such as
parsing struct tags, dealing with values, etc. Therefore, we take the following approach:

1. Consumers of the library create structs that are decorated with the `json` struct tag, as
documented by [the `encoding/json` package docs](https://pkg.go.dev/encoding/json#Marshal). The
tag value should be the name of the input, with spaces replaced by `_`. Casing does not matter.
2. Pass a pointer to the struct to the `Unmarshal` method, which in turn:
    1. Creates a JSON object of all the GitHub Action inputs and their values by parsing the environment
    variables injected by the Actions runtime. Simple types are converted to their JSON equivalent. For
    example, if an input as a value of "32", it will be converted to a floating point number; if an
    input has a value of "true", then it will be converted to a boolean.
    2. Uses the native `json.Unmarshal` method to populate the struct.

Here's a simple example:
```go
var myStruct struct{
  // Input name: doSomething
  DoSomething bool `json:"dosomething"`
  // Input name: a_string_value
  AStringValue string `json:"a_string_value"`
  // Input name: AnIntValue
  AnIntValue int `json:"anintvalue"`
}{}

err := wyatt.UnMarshal(&myStruct)
```

You can also have more complex types, such as lists, maps, or structs, by passing valid JSON
as the input. For example, say you have a struct like this:

```go
type Bong struct {
  Foo string `json:"foo"`
  Bar int    `json:"bar"`
}
type ActionConfig struct{
  Foo    string   `json:"foo"`
  Fooers []string `json:"fooers"`
  Bing   Bong     `json:"bing"`
}
```

The input to your Action might look something like this:
```yaml
- uses: someone/my-go-action
  with:
    foo: bar
    fooers: |
      [
        "joe",
        "jane"
      ]
    bing: |
      {
        "foo": "foo-er",
        "bar": "bar-er"
      }
```

## Roadmap

In a future iteration, we'd like Wyatt to also be able to translate from structs to [GitHub Actions
outputs](https://docs.github.com/en/actions/creating-actions/metadata-syntax-for-github-actions#outputs-for-docker-container-and-javascript-actions),
using a similar approach to the one we use for inputs.
