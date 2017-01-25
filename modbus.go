package modbuscli

import (
  "net"
)

const (
  READHOLDINGREGISTER = 0x03
  WRITEREGISTER = 0x06
  WRITEREGISTERS = 0x10
)

func intTo2Byte(val int) ([]byte) { /* binary.BigEndian.PutUint16 */
  b:= make([]byte,2)
  b[0] = byte(val >> 8)
  b[1] = byte(val)
  return b
}

func byteTo16int(val []byte) (int) {
  _ = val[1] // bounds check hint to compiler
  return int(val[1]) | int(val[0])<<8
}

type ModBus struct {
  EndPoint string
  Addr byte
  Code byte
  Data []byte
}

func (m *ModBus) envelope() ([]byte) {
  if m.Addr == byte(0) { m.Addr = 0x01 }
  head := []byte{0x00, 0x00, 0x00, 0x00, 0x00, byte(len(m.Data) + 2), m.Addr, m.Code}
  body := []byte{}
  body = append(body, head...)
  body = append(body, m.Data...)
  return body
}

func (m *ModBus) WriteRegister(regAddr int, regVal int) ([]int, error){
  m.Code = WRITEREGISTER
  m.Data = []byte{}
  m.Data = append(m.Data, intTo2Byte(regAddr)...)
  m.Data = append(m.Data, intTo2Byte(regVal)...)
  tmp, err := m.send()
  return []int{ byteTo16int(tmp[8:10]) , byteTo16int(tmp[10:12]) }, err
}

func (m *ModBus) ReadHoldingRegister(regAddr int, regSize int) ([]int, error) {
  m.Code = READHOLDINGREGISTER
  m.Data = []byte{}
  m.Data = append(m.Data, intTo2Byte(regAddr)...)
  m.Data = append(m.Data, intTo2Byte(regSize)...)
  tmp, err := m.send()
  outp := []int{}
  for i := 0; i < (regSize * 2); i+=2 {
    outp = append(outp, []int{byteTo16int(tmp[9+i:9+i+2])}...)
  }
  return outp, err
}

func (m *ModBus) WriteRegisters(regAddr int, regVal []int) ([]int, error) {
  m.Code = WRITEREGISTERS
  m.Data = []byte{}
  m.Data = append(m.Data, intTo2Byte(regAddr)...)
  m.Data = append(m.Data, intTo2Byte(len(regVal))...)
  m.Data = append(m.Data, []byte{byte(len(regVal)*2)}...)
  for _,v := range regVal{
    m.Data = append(m.Data, intTo2Byte(int(v))...)
  }
  tmp,err := m.send()
  return []int{ byteTo16int(tmp[8:10]) , byteTo16int(tmp[10:12]) }, err
}

func (m *ModBus) send() ([]byte, error) {
  outp := make([]byte, 0x40)
  addr, err := net.ResolveTCPAddr("tcp4", m.EndPoint); if err == nil {
    c, err := net.DialTCP("tcp", nil, addr); if err == nil {
      _, err = c.Write(m.envelope()); if err == nil {
        _, err := c.Read(outp); if err == nil {
          return outp, nil
        }
      }
    }
  }
  return []byte{}, err
}


