package writescript

import (
	"fmt"
	"github.com/paulvollmer/otto"
	"github.com/paulvollmer/textbackend"
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
		} else {
			// set the header to the first line
			w.Content.Writeln(header)
		}
	}

	// initialize otto
	vm := otto.New()
	// infos about the software
	vm.Set("VERSION", Version)
	// create api we can use at the plugin
	vm.Set("writeln", func(call otto.FunctionCall) otto.Value {
		// check if args are empty...
		if len(call.ArgumentList) == 0 {
			w.Content.Writeln("")
		} else {
			tmpLine := ""
			for l, v := range call.ArgumentList {
				// val, _ := call.Argument(0).ToString()
				val, errVal := v.ToString()
				if errVal != nil {
					fmt.Println("cannot convert variable", errVal)
				}
				fmt.Println(">>", l, val)
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

	// run the vm and get the result
	_, err := vm.Run(CreateVMScript(plugin, data))
	if err != nil {
		return err
	}
	return nil
}

// CreateVMScript creates the javascript script core wrapper.
func CreateVMScript(plugin, data string) string {
	//fmt.Println("CreateVMScript")

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
	// fmt.Println("SCRIPT:", script)

	return script
}
