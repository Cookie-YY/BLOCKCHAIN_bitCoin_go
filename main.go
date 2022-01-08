package main

import (
	"blcokChain/client"
)

/* ========= Welcome to the demo block-chain =========
@author: cookie
@course: Block Chain of PKU
@date: 2021-11 ~ 2022-01
*/

func main() {
	cli := client.GetCMDCli()
	cli.Run()
}
