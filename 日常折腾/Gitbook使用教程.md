# gitbook使用教程

> Document Everything! For you, your users and your team.

GitBook是一个基于Node.js的命令行工具，可使用Markdown来制作精美的电子书。

## 一、安装及使用

1. 首先前往官网[https://nodejs.org/](https://nodejs.org/)下载nodejs环境并安装，建议下载长期支持版；

2. 安装完成后，在命令行输入`npm install gitbook-cli -g`安装gitbook；

3. 新建文件夹，切到对应文件夹下执行`gitbook init`初始化电子书；

4. 使用`gitbook serve`运行，可在浏览器查看效果。

PS：运行gitbook命令时可能出现`因为在此系统上禁止运行脚本`的提示，使用命令`set-executionpolicy remotesigned`修改Powershell的执行策略。

## 二、常用插件

在电子书根目录增加`book.json`用于插件管理，基本模板如下：

```json
{
    "title": "菜鸟云计算工程师的学习笔记",
    "author": "ultraji",
    "links": {
        "sidebar": {
            "Home": "http://ultraji.xyz"
        }
    },
    "plugins": [
        "-lunr", "-search", "search-pro",
        "custom-favicon",
        "edit-link"
    ],
    "pluginsConfig": {
        "favicon": "/images/profile.jpg",
        "edit-link": {
            "base": "https://github.com/ultraji/the-notes-of-cloud-computing",
            "label": "编辑此页面"
        }
    }
}
```

然后，利用`gitbook install`安装json中列出的所有插件，或者`npm install gitbook-plugin-插件名`或`gitbook install 插件名`单独安装指定的插件。

### 2.1 search-pro

原生插件不支持中文搜索。搜索插件search-pro支持中文搜索，在使用此插件之前，需要将默认的search和lunr插件去掉。在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["-lunr", "-search", "search-pro"]
}
```

### 2.2 chapter-fold

gitbook支持多层目录，安装此插件后，点击导航栏的标题名就可以实现折叠扩展。在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["chapter-fold"]
}
```

### 2.3 anchor-navigation-ex

该插件会在页面右上角生成一个灰色的按钮，鼠标移入后会显示灰色的目录。在book.json的plugins参数中添加插件名：

```json
{
    "plugins" : ["anchor-navigation-ex"],
    "pluginsConfig": {
        "anchor-navigation-ex": {
            "showLevel": false, //标题是否显示层级序号.页面标题和导航中的标题都会加上层级显示。
            "showGoTop": false // 是否显示返回顶部按钮
        },
    }
}
```

### 2.4 back-to-top-button

回到顶部按钮，在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["back-to-top-button"],
}
```

### 2.5 edit-link

添加编辑此页面的超链接，在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["edit-link"],
    "pluginsConfig": {
      "edit-link": {
            "base": "https://github.com/ultraji/the-notes-of-cloud-computing",
            "label": "编辑此页面"
       }
    }
}
```

### 2.6 favicon

修改网站显示的图标，在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["favicon"],
    "pluginsConfig": {
      "favicon": {
            "shortcut": "asset/img/favicon.ico",
            "bookmark": "asset/img/favicon.ico",
            "appleTouch": "asset/img/favicon.ico",
            "appleTouchMore": {
                "120x120": "asset/img/favicon.ico",
                "180x180": "asset/img/favicon.ico"
            }
        }
    }
}
```

### 2.7 lightbox

点击图片弹窗显示，在book.json的plugins参数中添加插件名：

```json
{
    "plugins": ["lightbox"],
}
```

### 2.8 tbfed-pagefooter

添加页面显示最后更新时间，在book.json的plugins参数中添加插件名:

```json
{
    "plugins": ["tbfed-pagefooter"],
    "pluginsConfig": {
        "tbfed-pagefooter": {
            "copyright":"Copyright &copy ultraji.xyz 2019",
            "modify_label": "该文章修订时间：",
            "modify_format": "YYYY-MM-DD HH:mm:ss"
        },
    }
}
```