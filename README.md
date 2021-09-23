# Jencryption

一个图片加密算法

### 特点：

加密前后均为图片格式且加密后图片的尺寸与加密前保持一致。

加密Demo:

加密前: ![befo](docs/img.jpeg)  加密后: ![after](docs/encryption.png)

### 有啥用？

1、反扒

2、盗链 (将自己的私有图片加密后放到公开的cdn上，不怕对方知道图片内容)

3、防盗链（图片都加密了，别人盗过去也不知道杂用）

4、用于在公开论坛传输一些不符合社会主义核心价值观的图片

### 使用方法：

1、使用二进制文件加解密

加密当前文件夹内全部图片 `Jencryption encrypt [密码]`

加密指定文件夹/图片 `Jencryption encrypt [路径] [密码]`

解密当前文件夹 `Jencryption decrypt [密码]`

解密指定文件夹/图片 `Jencryption decrypt [路径] [密码] `

2、使用js sdk在你的网站接入

```
<script src="/js/md5.min.js"></script>
<script src="/js/core.js"></script>
<script>
DecryptAllImage("你的密码")  // 解密页面上全部图片

// DecryptImage(document.getElementById("img"),"你的密码")  // 解密单张图片


</script>
```

3、在线工具 [https://xjiangwei.cn/Jencryption/](https://xjiangwei.cn/Jencryption/)

### 有谁用？

[Jinnrry Blog 全站文章图片均加密 ](https://xjiangwei.cn)


