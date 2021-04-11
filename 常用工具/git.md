# git

## 一、常用命令

### 1. git commit --amend 改写最近一次提交

例如，

```shell
git commit --amend #改写最近一次提交
```

### 2. git cherry-pick <提交记录的hash值> 仅把在其他分支执行的某个修改记录导入到指定分支

例如，

```shell
git checkout master
git cherry-pick 99daed2
```

### 3. Git仓库删除全部历史提交

```
# 1.Checkout
git checkout --orphan latest_branch
# 2. Add all the files
git add -A
# 3. Commit the changes
git commit -am "init"
# 4. Delete the branch
git branch -D master
# 5.Rename the current branch to master
git branch -m master
# 6.Finally, force update your repository
git push -f origin master
```

### 4. 仓库进行fork后与原仓库同步

1. 进入本地仓库目录下，执行命令`git remote -v`，一般情况下只有两行指向自己库的链接，说明此时你没有设置过`upstream`（上游代码库）；

    ```shell
    origin  git@github.com:ultraji/halo.git (fetch)
    origin  git@github.com:ultraji/halo.git (push)
    ```

2. 首先需要设置一下上游代码库`upstream`，通过以下命令（命令中的链接为原仓库的链接，这个操作设置一次即可，无需重复操作）：`git remote add upstream git@github.com:halo-dev/halo.git`，再通过`git remote -v`就能看到`upstream`了

    ```shell
    origin  git@github.com:ultraji/halo.git (fetch)
    origin  git@github.com:ultraji/halo.git (push)
    upstream        git@github.com:halo-dev/halo.git (fetch)
    upstream        git@github.com:halo-dev/halo.git (push)
    ```

**接下来是merge的关键步骤**

3. 执行命令`git fetch upstream`抓取原仓库的更新：

4. 执行命令`git checkout master`切到master分支：

5. 执行命令`git merge upstream/master`合并远程的master分支：

6. 执行命令`git push`同步到远程仓库。

**提交Pull Request的正确姿势**：先同步上游修订，如果和本地有冲突, 先在本地解决冲突，然后再提Pull Request。

---

## 二、配置远程仓库（例如Github）

## 配置密钥对

本地仓库与远程GitHub仓库是通过SSH加密传输的，因此需要创建一对SSH key。

1. 创建SSH Key，必须填写注册github时的邮箱
    ```shell
    ssh-keygen -t rsa -C "ultraji@live.com"
    ```

2. 将公钥填入Github

    查看当前用户目录的.ssh文件夹中，找到新创建的key的公钥，通过cat命令查看公钥内容，复制所有内容到GitHub的SSH and GPG keys中。
    ```shell
    cat id_rsa.pub
    ```

3. 测试连接
    ```shell
    # 通过这个命令测试是否能连接上github服务器
    ssh -T git@github.com
    # 成功连接会显示以下内容
    Hi ultraji! You've successfully authenticated, but GitHub does not provide shell access.
    ```

4. 设置git配置文件
    ```shell
    git config --global user.name "ultraji"
    git config --global user.email "ultraji@live.com"
    # 打印配置信息
    git config --list
    ```

---

## 三、gitignore配置

### 忽略无后缀的文件

在gitignore前加上以下内容：

```txt
*
!*.*
!*/
```

## 资料

* [猴子都能懂的GIT入门](https://backlog.com/git-tutorial/cn/)