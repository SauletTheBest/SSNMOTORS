package model

type Product struct {
    ID          string  // auto-generated hex
    Name        string  // required
    Description string  // required
    Category    string  // required
    Stock       int32   // required, >=0
    Price       float64 // required, >=0
}
