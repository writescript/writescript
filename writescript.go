package writescript

import (
	"fmt"
	"github.com/writescript/otto"
	"github.com/writescript/textbackend"
)

// Version of the script engine.
const Version = "0.2.1"

// WriteScript Core
type WriteScript struct {
	Content textbackend.TextContent // create output storage the plugin can write content to
}

// Process the plugin generator.
func (w *WriteScript) Process(plugin, data, header string, headerOn bool) error {
	// do you want to write a header?
	if headerOn {
		if header == "" {
			// if no header was set, create a default header
			w.Content.Writeln("// written by writescript v" + Version)
			w.Content.Writeln("// DO NOT EDIT!")
			w.Content.Writeln("")
		} else {
			// set the header to the first line
			w.Content.Writeln(header)
		}
	}

	// initialize otto
	vm := otto.New()
	// infos about the software
	vm.Set("version", Version)
	// create api we can use at the plugin
	vm.Set("writeln", func(call otto.FunctionCall) otto.Value {
		// check if args are empty...
		if len(call.ArgumentList) == 0 {
			w.Content.Writeln("")
		} else {
			tmpLine := ""
			for _, v := range call.ArgumentList {
				val, errVal := v.ToString()
				if errVal != nil {
					fmt.Println("cannot convert variable", errVal)
				}
				tmpLine += val

			}

			w.Content.Writeln(tmpLine)
		}

		return otto.Value{}
	})
	// vm.Set("TODO: write", func(call otto.FunctionCall) otto.Value {
	// 	g.Content.AddLine(ContentLine{g.level, val})
	// 	return otto.Value{}
	// })
	vm.Set("pushLevel", func(call otto.FunctionCall) otto.Value {
		w.Content.PushLevel()
		return otto.Value{}
	})

	vm.Set("popLevel", func(call otto.FunctionCall) otto.Value {
		w.Content.PopLevel()
		return otto.Value{}
	})

	vm.Set("getLevel", func(call otto.FunctionCall) otto.Value {
		val := w.Content.GetLevel()
		result, _ := otto.ToValue(val)
		return result
	})

	vm.Set("setLevel", func(call otto.FunctionCall) otto.Value {
		val, err := call.Argument(0).ToInteger()
		if err != nil {
			w.Content.SetLevel(int(val))
		}
		return otto.Value{}
	})

	// run the vm and get the result
	_, err := vm.Run(CreateVMScript(plugin, data))
	if err != nil {
		return err
	}
	return nil
}

// CreateVMScript creates the javascript script core wrapper.
func CreateVMScript(plugin, data string) string {
	script := `
	function RUN(data) {
		` + plugin + `
	};`
	script += `RUN(`
	if data == "" {
		script += `{}` // if data is empty string, set it to an empty object
	} else {
		script += `JSON.parse('` + data + `')`
	}
	script += `);`
	return script
}
