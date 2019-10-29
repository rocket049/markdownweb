# 关于 MarkdownWeb
本程序用于建立简单的静态展示网站，直接将`Markdown`文件翻译成`HTML`，并且允许使用模板文件`md.tpl`美化页面。
网站的根目录是`files`目录，首页是`index.md`。
如果根目录或子目录中有`md.tpl`文件，会使用目录中的模板文件，否则将使用默认的`md.tpl`。
允许显示一组轮播图片，用`ad.json`指定图片列表和相应的`URL`，使用方法参考例子中的`md.tpl`和`ad.js`程序的配合方法。

# About MarkdownWeb
This program is used to build a simple static display website, translate the 'markdown' file directly into 'HTML', and allow the template file 'md.tpl' to beautify the page.
The root directory of the website is the 'files' directory. The homepage is `index.md`.
If there is an 'md.tpl' file in the directory (or sub directory), the template file in the directory will be used. Otherwise, the default 'md.tpl' will be used.
It is allowed to display a group of rotation pictures. Use 'ad.json' to specify the picture list and the corresponding 'URL'. Refer to the matching method of 'md.tpl' and 'ad.js' in the example for the usage method.