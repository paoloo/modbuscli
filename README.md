# modbuscli

A working Golang modbus client driver, implementing all 16 bits `Holding Register` functions:
**0x03**(*Read Holding Register*), **0x06**(*Write Single Register*) and
**0x10**(*Write Multiple Registers*) and handle their output. To use it, add to your project
imports `"github.com/paoloo/modbuscli"`, `go get` it then instance it and set the
endpoint(modbus TCP server):

```go
 mt := new(ModBus)
 mt.EndPoint = "127.0.0.1:502"
```
if you want to provide a **Slave Addres** other than **0x01**, just
```go
 mt.Addr = 0x05
```
then use:
- `res,_ := mt.WriteRegister(addr, value)` to write a single register **value** into **addr**.
   Returns ([]int, error) where first position of []int is the address writen and the second is
   the register itself;

- `res,_ := mt.ReadHoldingRegister(addr, size)` to read **size** 16 bits register from **addr**.
   Returns ([]int, error) where []int is an array of **size** values readen from **addr**;

- `res,_ := mt.WriteRegisters(addr,[]int{values})` to write an array of **values** into **addr**.
   Returns ([]int, error) where first position of []int is the addres written and the second is
   the amount of registers written.


### Enjoy
