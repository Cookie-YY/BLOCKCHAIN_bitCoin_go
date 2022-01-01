package main

import "fmt"

/* Stage2: Alice wants to transfer some money to Bob
1. init state:
	- cookie is the god-address, so it has the fist reward in the first block
	- other address don't have money
2. make Transaction:
	- makeTransaction: only one transaction in one time.
	- makeTransactions: can take several transactions in one time
3. after transaction:
	- remember to delete the json file to get the same result.
*/
func main() {
	// init state
	fmt.Print("\n\nInit State: \n")
	getBalanceOf("cookie") // cookie: 12.5
	getBalanceOf("Alice")  // Alice: 0
	getBalanceOf("Bob")    // Bob: 0

	// make transaction
	fmt.Print("\n\nStart Transaction: \n")
	makeTransaction("cookie", "Bob", 10, "Alice") // cookie: 2.5, Bob: 10, Alice: 12.5
	makeTransaction("Alice", "Bob", 12, "cookie") // cookie: 15, Alice: 0.5, Bob: 22
	makeTransactions(
		[]map[string]string{
			{"from": "Alice", "to": "Bob", "amount": "0.2"},
			{"from": "Alice", "to": "Bob", "amount": "0.1"},
		},
		"cookie")

	// after transaction
	fmt.Print("\n\nAfter Transaction: \n")
	getBalanceOf("cookie") // cookie: 15
	getBalanceOf("Alice")  // Alice: 0.2
	getBalanceOf("Bob")    // Bob: 22.3
}
