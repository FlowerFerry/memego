package memego

import "testing"

func TestAssginAndToString(t *testing.T) {
	t.Run("TestAssginAndToString 1", func(t *testing.T) {
		s, _ := CreateStringByGoStr("Hello, world!")
		defer s.destroy()
		if s.toString() != "Hello, world!" {
			t.Fatal("TestAssginAndToString failed")
		}
		if s.notEqualGoStr("Hello, world!") {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 2", func(t *testing.T) {
		s, _ := CreateStringByBytes([]byte("Hello, world!"))
		defer s.destroy()
		if s.toString() != "Hello, world!" {
			t.Fatal("TestAssginAndToString failed")
		}
		if s.notEqualGoStr("Hello, world!") {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 3", func(t *testing.T) {
		s1, _ := CreateStringByGoStr("Hello, world!")
		defer s1.destroy()
		s2, _ := CreateStringByOther(s1)
		defer s2.destroy()
		if s2.toString() != "Hello, world!" {
			t.Fatal("TestAssginAndToString failed")
		}
		if s2.notEqualGoStr("Hello, world!") {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 4", func(t *testing.T) {
		s := CreateString()
		defer s.destroy()
		gs := "Go（又称Golang）是Google开发的一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言。\n罗伯特·格瑞史莫（Robert Griesemer），罗布·派克（Rob Pike）及肯·汤普逊（Ken Thompson）于2007年9月开始设计Go，稍后Ian Lance Taylor、Russ Cox加入项目。Go是基于Inferno操作系统所开发的。Go于2009年11月正式宣布推出，成为开放源代码项目，并在Linux及Mac OS X平台上进行了实现，后来追加了Windows系统下的实现。在2016年，Go被软件评价公司TIOBE 选为“TIOBE 2016 年最佳语言”。 目前，Go每半年发布一个二级版本（即从a.x升级到a.y）。"
		s.assignByGoStr(gs)
		if s.toString() != gs {
			t.Fatal("TestAssginAndToString failed")
		}
		if s.notEqualGoStr(gs) {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 5", func(t *testing.T) {
		s := CreateString()
		defer s.destroy()
		gs := "Go is a statically typed, compiled high-level programming language designed at Google[11] by Robert Griesemer, Rob Pike, and Ken Thompson.[12] It is syntactically similar to C, but with memory safety, garbage collection, structural typing,[6] and CSP-style concurrency.[13] It is often referred to as Golang because of its former domain name, golang.org, but its proper name is Go."
		s.assignByBytes([]byte(gs))
		if s.toString() != gs {
			t.Fatal("TestAssginAndToString failed")
		}
		if s.notEqualGoStr(gs) {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 6", func(t *testing.T) {
		s1, _ := CreateStringByGoStr("Hello, world!")
		defer s1.destroy()
		s2 := CreateString()
		defer s2.destroy()
		s2.assignByOther(s1)
		if s2.toString() != "Hello, world!" {
			t.Fatal("TestAssginAndToString failed")
		}
		if s2.notEqualGoStr("Hello, world!") {
			t.Fatal("TestAssginAndToString failed")
		}
	})

	t.Run("TestAssginAndToString 7", func(t *testing.T) {
		s1, _ := CreateStringByGoStr("Hello, world!")
		defer s1.destroy()
		s2, _ := CreateStringByGoStr("Hello, world!")
		defer s2.destroy()
		if s1.notEqual(s2) {
			t.Fatal("TestAssginAndToString failed")
		}
	})

}
