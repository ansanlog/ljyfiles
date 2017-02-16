package com

import "mifpay-reconclt/bo"

/*集合操作*/

type SetHandle interface{
	Duplicate() Set
	Count() int
	Clear()
	Empty() bool
	Minus(sets ...Set)
	Intersect(sets ...Set)
}

type Set map[bo.Sbtest1]bool

// 新建集合对象
// 可以传入初始元素
/*func New(arr T) Set {
	s := make(Set, len(arr))
	s.Add(arr)
	return s
}*/

// 创建副本
func (s Set) Duplicate() Set {
	r := make(map[bo.Sbtest1]bool, len(s))
	for e := range s {
		r[e] = true
	}
	return r
}

// 添加元素
func (s Set) Add(item []bo.Sbtest1) {
	for _, v := range item {
		s[v] = true
	}
}

// 删除元素
func (s Set) Remove(items ...bo.Sbtest1) {
	for _, v := range items {
		delete(s, v)
	}
}

// 判断元素是否存在
func (s Set) Has(items ...bo.Sbtest1) bool {
	for _, v := range items {
		if _, ok := s[v]; !ok {
			return false
		}
	}
	return true
}

// 统计元素个数
func (s Set) Count() int {
	return len(s)
}

// 清空集合
func (s Set) Clear() {
	s = map[bo.Sbtest1]bool{}
}

// 空集合判断
func (s Set) Empty() bool {
	return len(s) == 0
}

// 差集
// 获取 s 与所有参数的差集，结果存入 s
func (s Set) Minus(sets ...Set) {
	for _, set := range sets {
		for e := range set {
			delete(s, e)
		}
	}
}

// 差集（函数）
// 获取第 1 个参数与其它参数的差集，并返回
func Minus(sets ...Set) Set {
	// 获取差集
	r := sets[0].Duplicate()
	for _, set := range sets[1:] {
		for e := range set {
			delete(r, e)
		}
	}
	return r
}

// 交集
// 获取 s 与其它参数的交集，结果存入 s
func (s Set) Intersect(sets ...Set) {
	for _, set := range sets {
		for e := range s {
			if _, ok := set[e]; !ok {
				delete(s, e)
			}
		}
	}
}

// 交集（函数）
// 获取所有参数的交集，并返回
func Intersect(sets ...Set) Set {
	// 获取交集
	r := sets[0].Duplicate()
	for _, set := range sets[1:] {
		for e := range r {
			if _, ok := set[e]; !ok {
				delete(r, e)
			}
		}
	}
	return r
}
