package heapImpl

/**
* @Author: super
* @Date: 2020-07-28 11:12
* @Description:
**/

type Int int

func (x Int) Less(than Item) bool {
	return x < than.(Int)
}