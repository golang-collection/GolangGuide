# 文件名

golang中文件名命名规则如下：

## 平台区分

文件名\_平台。

例： file\_windows.go, file\_unix.go

可选为`windows, unix, posix, plan9, darwin, bsd, linux, freebsd, nacl, netbsd, openbsd, solaris, dragonfly, bsd, notbsd， android，stubs`

## 单元测试

文件名_test.go或者 文件名_平台\_test.go。

例： path\_test.go, path\_windows\_test.go

## 版本区分

文件名\_版本号等。

例：trap\_windows\_1.4.go

## CPU类型区分

文件名\_\(平台:可选\)\_CPU类型.

例：vdso\_linux\_amd64.go

可选为`amd64, none, 386, arm, arm64, mips64, s390,mips64x,ppc64x, nonppc64x, s390x, x86,amd64p32`

