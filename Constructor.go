// Package constructor provides a minimal framework for creating and organizing command line
// Go applications. constructor is designed to be easy to understand and write, the most simple
// constructor application can be written as follows:
//
//	func main() {
//		(&constructor.App{}).Run(os.Args)
//	}
//
// Of course this application does not do much, so let's make this an actual application:
//
//	func main() {
//		app := &constructor.App{
//	  		Name: "greet",
//	  		Usage: "say a greeting",
//	  		Action: func(c *constructor.Context) error {
//	  			fmt.Println("Greetings")
//	  			return nil
//	  		},
//		}
//
//		app.Run(os.Args)
//	}
package constructor
