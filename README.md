# docSyncKit
文档同步工具，适用于本地个人知识库搭建

## 工作原理
1. 在github创建一个仓库用于存放文档资料
2. 本地Mac电脑安装git客户端，并配置Github ssh密钥，用于git提交
3. 安装本工具，并生成rsa密钥，需要自行备份保存
4. 补充配置文件，即可将本地文件加密存储到github


## 使用环境 
- 平台：Mac
- 账号：自行注册Github账号，并创建一个用于数据同步的仓库
- 依赖：git (开发环境 version 2.36.1)


## Github配置
参考：https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent

```shell
# 进入ssh配置目录
➜  docSyncKit git:(main) ✗ cd ~/.ssh
# 根据github账号创建私钥
➜  docSyncKit git:(main) ✗ ssh-keygen -t ed25519 -C "coderant@163.com"
Generating public/private ed25519 key pair.
Enter file in which to save the key (/Users/xxx/.ssh/id_ed25519): coderant@163.com  # 这里输入保存私钥的文件名
Enter passphrase (empty for no passphrase):   # 密码，这里直接按回车键
Enter same passphrase again:   # 重复密码，这里直接按回车键
Your identification has been saved in coderant@163.com
Your public key has been saved in coderant@163.com.pub
The key fingerprint is:
SHA256:m2OG3mtsqB7q/C10FYNiz/XNRRdM0fCxqUqtYGOgU3c coderant@163.com
The key's randomart image is:
+--[ED25519 256]--+
|       .     .=*=|
|    o . +     .+=|
|   . + + = E . o.|
|      = + o + .  |
|     o .S= . o   |
|    . o.oo+ o    |
|   ....o*  o     |
| . ..+.++.       |
| .+o+oooo.       |
+----[SHA256]-----+

# 后台启动ssh-agent
➜  docSyncKit git:(main) ✗ eval "$(ssh-agent -s)"
Agent pid 91373

# 编辑ssh配置文件添加github配置
➜  docSyncKit git:(main) ✗ vi ~/.ssh/config
Host github.com
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/coderant@163.com

# 加载文件
➜  docSyncKit git:(main) ✗ ssh-add --apple-use-keychain ~/.ssh/coderant@163.com
Identity added: /Users/xxx/.ssh/coderant@163.com (coderant@163.com)

```

## 开发

### 开发环境
mac book arm64
go version go1.20.6

### 编译
```shell
make clean && make
```
执行编译命令后，会生成build目录

### 运行
```shell
# 初始化生成自己的rsa私钥与公钥，需要自行备份./build/rsa_keys
./build/bin/docSyncKit rsa create
# 修改./build/conf/conf.toml配置
# 比如，WorkDir、PrivateKeyFile、PublicKeyFile、Repository等参数
vi ./build/conf/conf.toml

# 首先将项目从远端同步到本地。注意，本命令会从远端覆盖本地同名文件。如果本地已经作出更改，但是没同步到远端，会丢失。本操作适合在初次执行
./build/bin/docSyncKit sync local

# 后续若要将本地修改同步到远端。注意执行第一次同步操作后，会在conf/sync.txt记录上次同步的时间，下次执行同步，只会同步修改时间在该时间之后的文件
./build/bin/docSyncKit sync remote
```

## 后续
后续待开发内容：
1. 对比原始目录与本地git目录，找出差异的文件
2. 支持同步单个文件到远端
3. 监控本地目录修改，识别目录移动、文件重命名、文件移动、删除等事件
4. 实现多客户端的同步（目前只支持，一端写，多端只读）
