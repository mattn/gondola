package binary

import (
	"errors"
	"io"
	"reflect"
)

// Write writes the binary representation of data into w.
// Data must be a fixed-size value or a slice of fixed-size
// values, or a pointer to such data.
// Bytes written to w are encoded using the specified byte order
// and read from successive fields of the data.
// When writing structs, zero values are written for fields
// with blank (_) field names.
func Write(w io.Writer, order ByteOrder, data interface{}) error {
	// Fast path for basic types and slices of basic types
	var b [8]byte
	var bs []byte
	var err error
	switch v := data.(type) {
	case *int8:
		bs = b[:1]
		b[0] = byte(*v)
	case int8:
		bs = b[:1]
		b[0] = byte(v)
	case *uint8:
		bs = b[:1]
		b[0] = *v
	case uint8:
		bs = b[:1]
		b[0] = byte(v)
	case *int16:
		bs = b[:2]
		order.PutUint16(bs, uint16(*v))
	case int16:
		bs = b[:2]
		order.PutUint16(bs, uint16(v))
	case *uint16:
		bs = b[:2]
		order.PutUint16(bs, *v)
	case uint16:
		bs = b[:2]
		order.PutUint16(bs, v)
	case *int32:
		bs = b[:4]
		order.PutUint32(bs, uint32(*v))
	case int32:
		bs = b[:4]
		order.PutUint32(bs, uint32(v))
	case *uint32:
		bs = b[:4]
		order.PutUint32(bs, *v)
	case uint32:
		bs = b[:4]
		order.PutUint32(bs, v)
	case *int64:
		bs = b[:8]
		order.PutUint64(bs, uint64(*v))
	case int64:
		bs = b[:8]
		order.PutUint64(bs, uint64(v))
	case *uint64:
		bs = b[:8]
		order.PutUint64(bs, *v)
	case uint64:
		bs = b[:8]
		order.PutUint64(bs, v)
	case []int8:
		bs := b[:8]
		count := len(v)
		steps := count / 8
		i := 0
		for j := 0; j < steps; j++ {
			bs[0] = byte(v[i])
			i++
			bs[1] = byte(v[i])
			i++
			bs[2] = byte(v[i])
			i++
			bs[3] = byte(v[i])
			i++
			bs[4] = byte(v[i])
			i++
			bs[5] = byte(v[i])
			i++
			bs[6] = byte(v[i])
			i++
			bs[7] = byte(v[i])
			i++
			if _, err = w.Write(bs); err != nil {
				return err
			}
		}
		if i < count && err == nil {
			rem := count - i
			br := b[:rem]
			for j := 0; j < rem; j++ {
				br[j] = byte(v[i])
				i++
			}
			_, err = w.Write(br)
		}
		return err
	case []uint8:
		_, err := w.Write(v)
		return err
	case []int16:
		bs := b[:8]
		count := len(v)
		steps := count / 4
		i := 0
		for j := 0; j < steps; j++ {
			order.PutUint16(bs[:2], uint16(v[i]))
			i++
			order.PutUint16(bs[2:4], uint16(v[i]))
			i++
			order.PutUint16(bs[4:6], uint16(v[i]))
			i++
			order.PutUint16(bs[6:8], uint16(v[i]))
			i++
			if _, err = w.Write(bs); err != nil {
				return err
			}
		}
		if i < count {
			rem := (count - i) * 2
			br := b[:rem]
			for j := 0; j < rem; j += 2 {
				order.PutUint16(br[j:j+2], uint16(v[i]))
				i++
			}
			_, err = w.Write(br)
		}
		return err
	case []uint16:
		bs := b[:8]
		count := len(v)
		steps := count / 4
		i := 0
		for j := 0; j < steps; j++ {
			order.PutUint16(bs[:2], v[i])
			i++
			order.PutUint16(bs[2:4], v[i])
			i++
			order.PutUint16(bs[4:6], v[i])
			i++
			order.PutUint16(bs[6:8], v[i])
			i++
			if _, err = w.Write(bs); err != nil {
				return err
			}
		}
		if i < count {
			rem := (count - i) * 2
			br := b[:rem]
			for j := 0; j < rem; j += 2 {
				order.PutUint16(br[j:j+2], v[i])
				i++
			}
			_, err = w.Write(br)
		}
		return err
	case []int32:
		bs := b[:8]
		count := len(v)
		steps := count / 2
		i := 0
		for j := 0; j < steps; j++ {
			order.PutUint32(bs[:4], uint32(v[i]))
			i++
			order.PutUint32(bs[4:], uint32(v[i]))
			i++
			if _, err = w.Write(bs); err != nil {
				return err
			}
		}
		if i != count {
			b4 := b[:4]
			order.PutUint32(b4, uint32(v[i]))
			_, err = w.Write(b4)
		}
		return err
	case []uint32:
		bs := b[:8]
		count := len(v)
		steps := count / 2
		i := 0
		for j := 0; j < steps; j++ {
			order.PutUint32(bs[:4], v[i])
			i++
			order.PutUint32(bs[4:], v[i])
			i++
			if _, err = w.Write(bs); err != nil {
				return err
			}
		}
		if i != count {
			b4 := b[:4]
			order.PutUint32(b4, v[i])
			_, err = w.Write(b4)
		}
		return err
	case []int64:
		bs = b[:8]
		for _, val := range v {
			order.PutUint64(bs, uint64(val))
			if _, err = w.Write(bs); err != nil {
				break
			}
		}
		return err
	case []uint64:
		bs = b[:8]
		for _, val := range v {
			order.PutUint64(bs, val)
			if _, err = w.Write(bs); err != nil {
				break
			}
		}
		return err
	}
	if bs != nil {
		_, err := w.Write(bs)
		return err
	}

	// Fallback to reflect-based encoding.
	v := reflect.Indirect(reflect.ValueOf(data))
	if _, err := dataSize(v); err != nil {
		return errors.New("binary.Write: " + err.Error())
	}
	e := encoder{coder: coder{order: order}, writer: w}
	defer e.recover()
	e.value(v)
	return e.err
}
