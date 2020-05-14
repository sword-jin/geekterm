# Geekhub.com 终端版

![](./doc/m1.png)

### 安装
#### 下载二进制运行

进入 [release](https://github.com/rrylee/geekterm/releases) 下载对应系统的版本

#### 手动安装
```bash
git clone https://github.com/rrylee/geekterm.git

cd geekterm/cli

go run --mod=vendor . --cookie="你的浏览器cookie"
```

获取登录 cookie 请移动到 [wiki](https://github.com/rrylee/geekterm/wiki/%E6%89%8B%E5%8A%A8%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95-cookie)

### 常用操作

1. enter 或者 -> 进入帖子列表
2. j,k 上下，然后 enter 进入内容，继续 j 滑动到评论列表

#### 查看个人动态->回帖

1. i 进入个人动态
2. k 向下移动，j 向上移动，选择要查看的帖子
3. enter 进入帖子
4. k 向下移动，j 向上移动，选择评论列表
5. r 回帖，并且会回复当前楼层
6. R 直接回帖
7. 如果不想回帖了，ESC 返回

### 详细键位

如果你之前习惯使用 vim，那么这套操作你会非常舒服，基本上可以单手操作

|  键位      | 功能   |  说明  |
| --------   | -----:  | :----:  |
| h      | <-   |   可以在各个面板中跳转     |
| j        |   down   |    |
| k        |    up    |    |
| l        |    <-    |    |
| <-        |    <-    |  增加了跳转  |
| ->        |    ->    |  增加了跳转  |
| i        |    查看个人动态    |    |
| o        |    打开浏览器    |  选中帖子，打开帖子，在动态，打开动态  |
| M        |    进入留言列表    |    |
| r        |    回帖    |  如果有选中的留言，直接回复留言  |
| R        |    直接回帖    |    |
| m        |    上一页    |  帖子列表和留言  |
| n        |    下一页    |  帖子列表和留言  |
| q        |    退出    |    |
| enter    |    回复评论    | 在评论列表和个人中心  |
| enter    |    加载    | 加载选择帖子内容  |

### TODO
- [x] 基础功能
- [x] vim 键位
- [x] 从浏览器打开文章页面
- [x] 回贴
- [ ] 评论分页加载
- [ ] 自动签到
- [ ] 自动登录
- [ ] 帖子图片
- [ ] 分子详情展示
- [ ] 评论列表
- [ ] 配色自定义
- [ ] 自定义键位
- [ ] 未识别动态， 按enter 打开 github issue 页面
- [ ] 自定义签名

