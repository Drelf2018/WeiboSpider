<p align="center">
  <a href="https://github.com/Drelf2018/WeiboSpider/">
    <img src="https://user-images.githubusercontent.com/41439182/181650651-8a837630-3a75-47bd-9140-b9a0a14b6fac.png" alt="WeiboSpider">
  </a>
</p>

<div align="center">

# WeiboSpider
绝赞 golang 重构 [rubbish/weibo_detector](https://github.com/Drelf2018/rubbish/tree/main/weibo_detector) 中
</div>
计划利用 [go-python/cpy3](https://github.com/go-python/cpy3) 使用 Python Pillow 实现图片生成 ~~(懒)~~

错误的 是利用 Python 强大的图像处理技术提高生产效率

但是这个 `cpy3` 我研究了好一会没搞明白，贴两串错误代码看看

```go
> go run .\gothon\import.go
panic: object is nil

goroutine 1 [running]:
main.main()
        C:/Users/drelf/go/src/github.com/Drelf2018/WeiboSpider/gothon/import.go:33 +0x434
exit status 2
```
```go
> go run .\gothon\import.go
200
[MODULE] repr(hello) = <module 'hello' from '.\\gothon\\hello.py'>
[VARS] a = "10"
Exception 0xc0000005 0x0 0x8 0x7ffda9322f15
PC=0x7ffda9322f15
signal arrived during external code execution

runtime.cgocall(0xea49a0, 0xc0000cde28)
        C:/Program Files/Go/src/runtime/cgocall.go:157 +0x4a fp=0xc0000cde00 sp=0xc0000cddc8 pc=0xe137aa
github.com/go-python/cpy3._Cfunc_PyUnicode_AsUTF8(0x0)
        _cgo_gotypes.go:4939 +0x57 fp=0xc0000cde28 sp=0xc0000cde00 pc=0xe9eb77
github.com/go-python/cpy3.PyUnicode_AsUTF8.func1(0xc00004e070?)
        C:/Users/drelf/go/pkg/mod/github.com/go-python/cpy3@v0.2.0/unicode.go:89 +0x3a fp=0xc0000cde60 sp=0xc0000cde28 pc=0xe9fe1a
github.com/go-python/cpy3.PyUnicode_AsUTF8(0xedb3d2?)
        C:/Users/drelf/go/pkg/mod/github.com/go-python/cpy3@v0.2.0/unicode.go:89 +0x19 fp=0xc0000cde78 sp=0xc0000cde60 pc=0xe9fdb9
main.main()
        C:/Users/drelf/go/src/github.com/Drelf2018/WeiboSpider/gothon/import.go:49 +0x236 fp=0xc0000cdf80 sp=0xc0000cde78 pc=0xea0f36
runtime.main()
        C:/Program Files/Go/src/runtime/proc.go:250 +0x1fe fp=0xc0000cdfe0 sp=0xc0000cdf80 pc=0xe464be
runtime.goexit()
        C:/Program Files/Go/src/runtime/asm_amd64.s:1571 +0x1 fp=0xc0000cdfe8 sp=0xc0000cdfe0 pc=0xe6ce01
rax     0xc0000ce000
rbx     0xc0000cde28
rcx     0x0
rdi     0xc0000ce000
rsi     0xf64060
rbp     0xc0000cddf0
rsp     0xbf845ffaa0
r8      0xf643c0
r9      0x0
r10     0x227b102da30
r11     0xbf845ffa70
r13     0x10
r14     0xc000044000
r15     0x227b059ffff
rip     0x7ffda9322f15
rflags  0x10206
cs      0x33
gs      0x2b
exit status 2
```

### 2022.7.27 3: 12 已完成微博生成图片开发
