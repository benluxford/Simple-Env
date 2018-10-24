// Copyright Ben Luxford.
// All Rights Reserved

/*
Package senv (Senv) is a no bells, damn straight, dead simple, environment variable management package.



Save your environment variables in a file as follows:

	// VALUES - Inline style comments are ignored.

	key = value
	key2 = value
	...

	// VALUE GROUP - Tabs are ignored.
		etc = another value

	----------------
	Random block text is also ignored.
	----------------
	Sed auctor fermentum sollicitudin. Orci varius natoque penatibus et magnis dis parturient montes,
	nascetur ridiculus mus. Aliquam euismod lobortis purus. Nunc in arcu nec lorem fermentum fringilla
	vel ut turpis. Sed id magna sodales, dignissim nunc vitae, mollis diam. Etiam volutpat elementum libero,
	luctus euismod nisl convallis vehicula. Cras justo lacus, feugiat ac dapibus id, faucibus ac ipsum. Quisque
	ultricies feugiat ornare. Curabitur id mattis arcu. Suspendisse vel risus vel turpis placerat vehicula eu id sem.
	----------------

All keys and values have whitespace trimmed. Keys are converted to:
	ALL_UPPERCASE_AND_UNDERSCORED

A custom prefix is applied to all variable keys to ensure no clashes within the system:
	PREFIX_ALL_UPPERCASE_AND_UNDERSCORED

The package was designed to be a basic parser and environment variable schema, between Golang and Node
applications. To keep environment variables within scope of Node applications, these applications MUST be
direct children of the go application using Senv.

For example "within go application" - start a basic Node application:

	cmd := exec.Command("node", "path/to/index.js")
	err := cmd.Run()
	handle error...

Accessing variables within Node:
	// Example get function
	let GetVar = key = process.env[`${process.env.SENV_PREFIX}_${key.toUpperCase().replace(" ", "_")}`]

	// Example get
	let variable = GetVar("some key")

*/
//
package senv
